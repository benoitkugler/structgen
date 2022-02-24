package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	darttypes "github.com/benoitkugler/structgen/dart-types"
	"github.com/benoitkugler/structgen/data"
	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/interfaces"
	"github.com/benoitkugler/structgen/loader"
	"github.com/benoitkugler/structgen/orm/composites"
	"github.com/benoitkugler/structgen/orm/creation"
	"github.com/benoitkugler/structgen/orm/crud"
	tstypes "github.com/benoitkugler/structgen/ts-types"
)

type mode struct {
	mode   string
	output string
}

type Modes []mode

func (i *Modes) String() string {
	if i == nil {
		return ""
	}
	return fmt.Sprint(*i)
}

func (i *Modes) Set(value string) error {
	chuncks := strings.Split(value, ":")
	if len(chuncks) != 2 {
		return fmt.Errorf("expected colon separated <mode>:<output>, got %s", value)
	}
	m := mode{mode: chuncks[0], output: chuncks[1]}
	if m.output == "" {
		return fmt.Errorf("output not specified for mode %s", m.mode)
	}
	*i = append(*i, m)
	return nil
}

// check if the goimports command is working
func hasGoFormatter() bool {
	err := exec.Command("which", "goimports").Run()
	if err != nil {
		log.Printf("no formatter for Go (%s)", err)
	} else {
		log.Println("formatter for Go detected")
	}
	return err == nil
}

// check if the dart command is working
func hasDartFormatter() bool {
	err := exec.Command("dart", "format", "--help").Run()
	if err != nil {
		log.Printf("no formatter for Dart (%s)", err)
	} else {
		log.Println("formatter for Dart detected")
	}
	return err == nil
}

// check if the prettier command is working
func hasTypescriptFormatter() bool {
	err := exec.Command("npx", "prettier", "-v").Run()
	if err != nil {
		log.Printf("no formatter for Typescript (%s)", err)
	} else {
		log.Println("formatter for Typescript detected")
	}
	return err == nil
}

func main() {
	const (
		noFormat = iota
		formatGo
		formatDart
		formatTs
	)

	hasGoFmt := hasGoFormatter()
	hasDartFmt := hasDartFormatter()
	hasTsFmt := hasTypescriptFormatter()

	source := flag.String("source", "", "go source file to convert")
	var modes Modes
	flag.Var(&modes, "mode", "list of modes <mode>:<output>")

	flag.Parse()
	if source == nil || *source == "" {
		log.Fatal("Please define input source file")
	}
	if len(modes) == 0 {
		return
	}

	fullPath, err := filepath.Abs(*source)
	if err != nil {
		log.Fatal(err)
	}

	pkg, err := loader.LoadSource(*source)
	if err != nil {
		log.Fatal(err)
	}

	en, err := enums.FetchEnums(pkg)
	if err != nil {
		log.Fatal(err)
	}

	packageName := filepath.Base(pkg.ID)
	for _, m := range modes {
		var (
			typeHandler loader.Handler
			format      int // format if true
		)
		switch m.mode {
		case "ts":
			typeHandler = tstypes.NewHandler(en)
			format = formatTs
		case "dart":
			typeHandler = darttypes.NewHandler(en)
			format = formatDart
		case "itfs-json":
			typeHandler = interfaces.NewHandler(packageName)
			format = formatGo
		case "rand":
			typeHandler = data.NewHandler(packageName, en)
			format = formatGo
		case "sql":
			typeHandler = crud.NewHandler(packageName, false)
			format = formatGo
		case "sql_test":
			typeHandler = crud.NewHandler(packageName, true)
			format = formatGo
		case "sql_gen":
			typeHandler = creation.NewGenHandler(en)
		case "sql_composite":
			typeHandler = &composites.Composites{OriginPackageName: packageName}
			format = formatGo
		case "enums":
			typeHandler = enums.Handler{PackageName: packageName, Enums: en}
			format = formatGo
		default:
			log.Printf("mode %s not supported - skipping \n", m.mode)
		}

		decls, err := loader.WalkFile(fullPath, pkg, typeHandler)
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.Create(m.output)
		if err != nil {
			log.Fatal(err)
		}

		err = decls.Generate(f, typeHandler)
		if err != nil {
			log.Fatal(err)
		}

		if err = f.Close(); err != nil {
			log.Fatal(err)
		}

		switch format {
		case formatGo:
			if hasGoFmt {
				err = exec.Command("goimports", "-w", m.output).Run()
			}
		case formatDart:
			if hasDartFmt {
				err = exec.Command("dart", "format", m.output).Run()
			}
		case formatTs:
			if hasTsFmt {
				err = exec.Command("npx", "prettier", "--write", m.output).Run()
			}
		}

		if err != nil {
			log.Fatalf("formatting failed: generated code is probably incorrect: %s", err)
		}

		log.Printf("Code for mode %s written in %s \n", m.mode, m.output)
	}
	log.Println("Done.")
}
