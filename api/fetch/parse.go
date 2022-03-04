package fetch

import (
	"errors"
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"path/filepath"

	"github.com/benoitkugler/structgen/api/gents"
	"golang.org/x/tools/go/packages"
)

func isHttpMethod(name string) bool {
	switch name {
	case "GET", "PUT", "POST", "DELETE":
		return true
	default:
		return false
	}
}

// LoadSource loads the given Go source file, loading
// its package and the AST.
func LoadSource(sourceFile string) (*packages.Package, *ast.File, error) {
	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedImports | packages.NeedDeps}
	pkgs, err := packages.Load(cfg, "file="+sourceFile)
	if err != nil {
		return nil, nil, err
	}
	if len(pkgs) != 1 {
		return nil, nil, fmt.Errorf("only one package expected, got %d", len(pkgs))
	}

	absSourceFile, err := filepath.Abs(sourceFile)
	if err != nil {
		return nil, nil, err
	}

	pkg := pkgs[0]
	for _, file := range pkg.Syntax { // restrict to input file for simplicity
		if absSourceFile == pkg.Fset.File(file.Package).Name() {
			return pkg, file, nil
		}
	}

	return nil, nil, fmt.Errorf("internal error: file not found %s", absSourceFile)
}

// Parse looks for method calls .GET .POST .PUT .DELETE
// inside all top level functions in `f`.
func Parse(pkg *packages.Package, f *ast.File) gents.Service {
	var out gents.Service
	for _, decl := range f.Decls {
		funcStm, ok := decl.(*ast.FuncDecl)
		if !ok || funcStm.Body == nil {
			continue
		}

		for _, stm := range funcStm.Body.List {
			call, ok := stm.(*ast.ExprStmt)
			if !ok {
				continue
			}
			callExpr, ok := call.X.(*ast.CallExpr)
			if !ok {
				continue
			}
			selector, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			methodName := selector.Sel.Name
			if !isHttpMethod(methodName) || len(callExpr.Args) < 2 {
				// we are looking for .<METHOD>(url, handler)
				continue
			}
			path, err := parseArgPath(callExpr.Args[0], pkg, f.Imports)
			if err != nil {
				fmt.Printf("%s\n", err)
				continue
			}
			contrat, err := parseArgHandler(callExpr.Args[1], pkg)
			if err != nil {
				fmt.Printf("%s\n", err)
				continue
			}
			out = append(out, gents.API{Url: path, Method: methodName, Contrat: contrat})
		}
	}
	return out
}

// look for ident in the imported packages of the source file
func isImportedPacakge(ident *ast.Ident, pkg *packages.Package, fileImports []*ast.ImportSpec) (*packages.Package, bool) {
	for _, imported := range fileImports {
		impPkg := pkg.Imports[stringLitteral(imported.Path)]
		pkgName := impPkg.Name
		if imported.Name != nil { // use local package name
			pkgName = imported.Name.String()
		}
		if pkgName == ident.Name { // use the local name of imported package
			return impPkg, true
		}
	}
	return nil, false
}

func resolveStringConst(arg *ast.Ident, pkg *packages.Package, local bool) (string, error) {
	var obj types.Object
	if local {
		// start by local scope
		localScope := pkg.Types.Scope().Innermost(arg.Pos())
		if localScope != nil {
			obj = localScope.Lookup(arg.Name)
		}
	}
	if obj == nil { // package scope
		obj = pkg.Types.Scope().Lookup(arg.Name)
	}
	if obj == nil {
		return "", fmt.Errorf("can't resolve constant at %s", pkg.Fset.Position(arg.Pos()))
	}
	val := obj.(*types.Const).Val()
	if val.Kind() == constant.String {
		return constant.StringVal(val), nil
	}
	return "", fmt.Errorf("can't resolve constant at %s", pkg.Fset.Position(arg.Pos()))
}

func parseAddStrings(x, y ast.Expr, pkg *packages.Package, fileImports []*ast.ImportSpec) (string, error) {
	valueX, err := parseArgPath(x, pkg, fileImports)
	if err != nil {
		return "", err
	}
	valueY, err := parseArgPath(y, pkg, fileImports)
	if err != nil {
		return "", err
	}
	return valueX + valueY, nil
}

func stringLitteral(arg *ast.BasicLit) string {
	if arg.Kind == token.STRING { // remove quotes
		return arg.Value[1 : len(arg.Value)-1]
	}
	return ""
}

// we support string litteral or string const
func parseArgPath(arg ast.Expr, pkg *packages.Package, fileImports []*ast.ImportSpec) (string, error) {
	switch arg := arg.(type) {
	case *ast.Ident:
		if arg.Obj.Kind == ast.Con { // constant of the package
			return resolveStringConst(arg, pkg, true)
		}
	case *ast.SelectorExpr: // looking for imported constants
		if pkgIdent, ok := arg.X.(*ast.Ident); ok {
			if pkgImported, ok := isImportedPacakge(pkgIdent, pkg, fileImports); ok {
				return resolveStringConst(arg.Sel, pkgImported, false)
			}
		}
	case *ast.BinaryExpr:
		if arg.Op == token.ADD {
			return parseAddStrings(arg.X, arg.Y, pkg, fileImports)
		}
	case *ast.BasicLit:
		if out := stringLitteral(arg); out != "" {
			return out, nil
		}
	}
	return "", fmt.Errorf("ignoring invalid url at %s", pkg.Fset.Position(arg.Pos()))
}

func resolveMethodReceiver(x *ast.Ident, pkg *packages.Package) *types.Named {
	localScope := pkg.Types.Scope().Innermost(x.Pos())
	obj := localScope.Lookup(x.Name)
	if obj == nil {
		obj = pkg.Types.Scope().Lookup(x.Name)
	}
	if obj == nil {
		panic(fmt.Sprintf("can't resolve name %s", x.Name))
	}

	type_ := obj.Type()
	if ptr, isPointer := type_.(*types.Pointer); isPointer { // remove indirection
		type_ = ptr.Elem()
	}

	if named, ok := type_.(*types.Named); ok {
		return named
	}
	panic(fmt.Sprintf("unexpected type for %s: %T", x.Name, type_))
}

func extractMethodBody(f *ast.File, pos token.Pos) (body []ast.Stmt, err error) {
	for _, decl := range f.Decls {
		funcDecl, isFunc := decl.(*ast.FuncDecl)
		if !isFunc {
			continue
		}
		if funcDecl.Name.NamePos == pos {
			return funcDecl.Body.List, nil
		}
	}
	return nil, errors.New("method not found")
}

// return the body of `fn`
func findMethod(fn *types.Func, rootPkg *packages.Package) (body []ast.Stmt, err error) {
	declFile := rootPkg.Fset.Position(fn.Pos()).Filename

	search := func(pkg *packages.Package) *ast.File {
		for i, file := range pkg.GoFiles {
			if file == declFile {
				return pkg.Syntax[i]
			}
		}
		return nil
	}
	// search in current package
	if f := search(rootPkg); f != nil {
		return extractMethodBody(f, fn.Pos())
	}
	// search into imports
	for _, importedPkg := range rootPkg.Imports {
		if f := search(importedPkg); f != nil {
			return extractMethodBody(f, fn.Pos())
		}
	}
	return nil, errors.New("method not found")
}

func parseArgHandler(arg ast.Expr, pkg *packages.Package) (gents.Contrat, error) {
	if method, ok := arg.(*ast.SelectorExpr); ok {
		if ident, ok := method.X.(*ast.Ident); ok {
			named := resolveMethodReceiver(ident, pkg)
			for i := 0; i < named.NumMethods(); i++ {
				fn := named.Method(i)
				if method.Sel.Name == fn.Name() {
					funcBody, err := findMethod(fn, pkg)
					if err != nil {
						return gents.Contrat{}, err
					}
					contrat := analyzeHandler(funcBody, named.Obj().Pkg())
					contrat.HandlerName = fn.Name()
					return contrat, nil
				}
			}
		}
	}

	return gents.Contrat{}, fmt.Errorf("ignoring invalid handler at %s : only methods are supported", pkg.Fset.Position(arg.Pos()))
}
