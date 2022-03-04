package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	darttypes "github.com/benoitkugler/structgen/dart-types"
	"github.com/benoitkugler/structgen/data"
	"github.com/benoitkugler/structgen/enums"
	"github.com/benoitkugler/structgen/formatter"
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

var fmts formatter.Formatters

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

	packageName := pkg.Name
	for _, m := range modes {
		var (
			typeHandler loader.Handler
			format      formatter.Format // format if true
		)
		switch m.mode {
		case "ts":
			typeHandler = tstypes.NewHandler(en)
			format = formatter.Ts
		case "dart":
			typeHandler = darttypes.NewHandler(en)
			format = formatter.Dart
		case "itfs-json":
			typeHandler = interfaces.NewHandler(packageName)
			format = formatter.Go
		case "rand":
			typeHandler = data.NewHandler(packageName, en)
			format = formatter.Go
		case "sql":
			typeHandler = crud.NewHandler(packageName, false)
			format = formatter.Go
		case "sql_test":
			typeHandler = crud.NewHandler(packageName, true)
			format = formatter.Go
		case "sql_gen":
			typeHandler = creation.NewGenHandler(en)
		case "sql_composite":
			typeHandler = &composites.Composites{OriginPackageName: packageName}
			format = formatter.Go
		case "enums":
			typeHandler = enums.Handler{PackageName: packageName, Enums: en}
			format = formatter.Go
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

		err = fmts.FormatFile(format, m.output)
		if err != nil {
			log.Fatalf("formatting failed: generated code is probably incorrect: %s", err)
		}

		log.Printf("Code for mode %s written in %s \n", m.mode, m.output)
	}
	log.Println("Done.")
}
