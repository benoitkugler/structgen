package loader

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io"
	"strings"

	"github.com/benoitkugler/structgen/utils"
	"golang.org/x/tools/go/packages"
)

type Handler interface {
	WriteHeader(w io.Writer) error
	HandleType(topLevelDecl *Declarations, typ types.Type)
	HandleComment(topLevelDecl *Declarations, comment Comment) error
	WriteFooter(w io.Writer) error
}

type Declaration interface {
	Render() string
	Id() string
}

// Declarations stores all top-level declarations
// to write.
type Declarations struct {
	list []Declaration
	keys map[string]struct{}
}

func NewDeclarations() *Declarations {
	return &Declarations{keys: map[string]struct{}{}}
}

func (ds Declarations) List() []Declaration {
	return ds.list
}

func (ds *Declarations) Add(d Declaration) {
	if _, alreadyHere := ds.keys[d.Id()]; !alreadyHere {
		ds.list = append(ds.list, d)
		ds.keys[d.Id()] = struct{}{}
	}
}

func (ds Declarations) ToString() string {
	var out strings.Builder
	for _, decl := range ds.list {
		_, _ = out.WriteString(decl.Render()) // err is nil
		_, _ = out.WriteString("\n")          // idem
	}
	return out.String()
}

func (ds Declarations) Render(out io.Writer) error {
	_, err := io.WriteString(out, ds.ToString())
	return err
}

type Comment struct {
	TypeName string
	Tag      string
	Content  string
}

func LoadSource(sourceFile string) (*packages.Package, error) {
	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedImports | packages.NeedDeps}
	pkgs, err := packages.Load(cfg, "file="+sourceFile)
	if err != nil {
		return nil, err
	}
	if len(pkgs) != 1 {
		return nil, fmt.Errorf("only one package expected, got %d", len(pkgs))
	}
	return pkgs[0], nil
}

func WalkFile(absPathOrigin string, pkg *packages.Package, handler Handler) (*Declarations, error) {
	scope := pkg.Types.Scope()
	fset := pkg.Fset
	accu := NewDeclarations()
	for _, name := range scope.Names() {
		object := scope.Lookup(name)
		if fset.Position(object.Pos()).Filename != absPathOrigin {
			// retrict to file declaration
			continue
		}
		if _, isTypeName := object.(*types.TypeName); !isTypeName {
			// ignore non-type declarations
			continue
		}
		handler.HandleType(accu, object.Type())
	}

	for _, file := range pkg.Syntax {
		if fset.Position(file.Pos()).Filename != absPathOrigin {
			// retrict to file declaration
			continue
		}
		for _, decl := range file.Decls {
			if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.TYPE && decl.Doc != nil {
				typeName := decl.Specs[0].(*ast.TypeSpec).Name.String()
				for _, line := range decl.Doc.List {
					if tag, content := utils.IsSpecialComment(line.Text); tag != "" {
						err := handler.HandleComment(accu, Comment{
							TypeName: typeName,
							Tag:      tag,
							Content:  content,
						})
						if err != nil {
							return nil, err
						}
					}
				}
			}
		}
	}
	return accu, nil
}
