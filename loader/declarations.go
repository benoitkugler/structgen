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

type Handler interface {
	// HandleType process the given type, creating a Declaration
	// and storing it (if needed) into `topLevelDecl`.
	HandleType(topLevelDecl *Declarations, typ types.Type)

	// HandleComment process special comments used to specify cnntext-dependent
	// information.
	HandleComment(topLevelDecl *Declarations, comment Comment) error

	// WriteHeader writes the start of the generated content.
	// After being called, all the stored declarations are rendered
	// and added to the file, then WriteFooter is called.
	WriteHeader(w io.Writer) error

	// WriteFooter writes the remaining generated instructions
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
	keys map[string]bool
}

func NewDeclarations() *Declarations {
	return &Declarations{keys: map[string]bool{}}
}

func (ds Declarations) List() []Declaration { return ds.list }

func (ds *Declarations) Add(d Declaration) {
	if alreadyHere := ds.keys[d.Id()]; !alreadyHere {
		ds.list = append(ds.list, d)
		ds.keys[d.Id()] = true
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
