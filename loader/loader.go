package loader

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io"
	"path/filepath"
	"strings"

	"github.com/benoitkugler/structgen/utils"
	"golang.org/x/tools/go/packages"
)

// Type is the common interface shared by all items
// to generate
type Type interface {
	// Returns the declarations to write
	Render() []Declaration
}

// Declaration is a top level declaration to write to the generated file.
type Declaration struct {
	Id      string // uniquely identifies the item, used to avoid duplicated declarations
	Content string // actual code to write
}

// Handler handles the specifity of the generated target.
// Is is used by calling `HandleType` and `HandleComment` for each type and comments
// in the source file, then calling Header(), then Declaration.Render(), then Footer().
type Handler interface {
	// HandleType process the given type, returning the corresponding Declaration.
	// It may return nil to ignore the type.
	// If needed, additional information way be stored and used in Header() or Footer().
	HandleType(typ types.Type) Type

	// HandleComment process special comments used to specify context-dependent
	// information, allowing to modify the returned values of Header() or Footer().
	HandleComment(comment Comment) error

	// Header returns the start of the generated content.
	// After being called, all the stored declarations are rendered
	// and added to the file, then Footer is called.
	Header() string

	// Footer returns the remaining generated instructions
	Footer() string
}

// Declarations stores all top-level declarations
// to write.
type Declarations []Type

func (ds Declarations) render() (out []Declaration) {
	for _, ty := range ds {
		out = append(out, ty.Render()...)
	}
	return out
}

// Render remove duplicates and aggregate the declarations
func Render(decls []Declaration) string {
	keys := map[string]bool{}
	var out strings.Builder
	for _, decl := range decls {
		if alreadyHere := keys[decl.Id]; !alreadyHere {
			keys[decl.Id] = true
			out.WriteString(decl.Content)
			out.WriteByte('\n')
		}
	}
	return out.String()
}

func (ds Declarations) Generate(out io.Writer, handler Handler) error {
	_, err := io.WriteString(out, handler.Header()+"\n")
	if err != nil {
		return fmt.Errorf("generating header: %s", err)
	}

	if _, err = io.WriteString(out, Render(ds.render())); err != nil {
		return fmt.Errorf("generating declarations: %s", err)
	}

	_, err = io.WriteString(out, handler.Footer())
	if err != nil {
		return fmt.Errorf("generating footer: %s", err)
	}

	return nil
}

// Comment is a special comment with the following syntax // <tag>:<content>
type Comment struct {
	TypeName string // the type where the comment was found
	Tag      string // the first part of the comment
	Content  string // the actual content
}

func LoadSource(sourceFile string) (*packages.Package, error) {
	dir := filepath.Dir(sourceFile)
	cfg := &packages.Config{
		Dir: dir,
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax |
			packages.NeedTypes | packages.NeedImports | packages.NeedDeps,
	}
	pkgs, err := packages.Load(cfg, "file="+sourceFile)
	if err != nil {
		return nil, err
	}
	if len(pkgs) != 1 {
		return nil, fmt.Errorf("only one package expected, got %d", len(pkgs))
	}
	if len(pkgs[0].Errors) != 0 {
		return nil, fmt.Errorf("errors during package loading:\n%v", pkgs[0].Errors)
	}
	return pkgs[0], nil
}

// WalkFile uses the package information to analyse the defined types.
func WalkFile(absPathOrigin string, pkg *packages.Package, handler Handler) (Declarations, error) {
	scope := pkg.Types.Scope()
	fset := pkg.Fset

	var accu Declarations
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

		decls := handler.HandleType(object.Type())
		if decls != nil {
			accu = append(accu, decls)
		}
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
						err := handler.HandleComment(Comment{
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
