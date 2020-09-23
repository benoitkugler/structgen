package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/benoitkugler/structgen/data"
	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/loader"
	"github.com/benoitkugler/structgen/orm/composites"
	"github.com/benoitkugler/structgen/orm/creation"
	"github.com/benoitkugler/structgen/orm/crud"
	"github.com/benoitkugler/structgen/tstypes"
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

func main() {
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
		var typeHandler loader.Handler
		switch m.mode {
		case "ts":
			typeHandler = tstypes.NewHandler(en)
		case "rand":
			typeHandler = data.Handler{PackageName: packageName, EnumsTable: en}
		case "sql":
			typeHandler = crud.NewHandler(packageName, false)
		case "sql_test":
			typeHandler = crud.NewHandler(packageName, true)
		case "sql_gen":
			typeHandler = creation.NewGenHandler(en)
		case "sql_composite":
			typeHandler = &composites.Composites{OriginPackageName: packageName}
		case "enums":
			typeHandler = enums.Handler{PackageName: packageName, Enums: en}
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

		if err := typeHandler.WriteHeader(f); err != nil {
			log.Fatal(err)
		}

		if err := decls.Render(f); err != nil {
			log.Fatal(err)
		}

		if err := typeHandler.WriteFooter(f); err != nil {
			log.Fatal(err)
		}

		if err := f.Close(); err != nil {
			log.Fatal(err)
		}

		log.Printf("Code for mode %s written in %s \n", m.mode, m.output)
	}
	log.Println("Done.")
}
