package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
)

type Config struct {
	ModuleName  string
	ServiceName string
	OutputPath  string
}

func main() {
	cfg := Config{}
	flag.StringVar(&cfg.ModuleName, "module-name", "", "the name of the new module")
	flag.StringVar(&cfg.ServiceName, "service-name", "", "the name of the new service")
	flag.StringVar(&cfg.OutputPath, "output", "", "the path to output the files")
	flag.Parse()

	if cfg.ModuleName == "" {
		fmt.Println("Module name cannot be empty!")
		os.Exit(1)
	}

	if cfg.ServiceName == "" {
		fmt.Println("Service name cannot be empty!")
		os.Exit(1)
	}

	if err := parseTemplates(cfg); err != nil {
		log.Fatalln(err)
	}
}

func parseTemplates(cfg Config) error {
	return filepath.WalkDir("./templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		newFilePath := strings.Join(strings.Split(path, "templates/")[1:], "")
		newFilePath = strings.ReplaceAll(newFilePath, ".tmpl", "")
		outputLocation := fmt.Sprintf("%s/%s", cfg.OutputPath, newFilePath)
		if err := os.MkdirAll(filepath.Dir(outputLocation), 0755); err != nil {
			if !os.IsExist(err) {
				return fmt.Errorf("failed creating dir: %w", err)
			}
		}
		f, err := os.Create(outputLocation)
		if err != nil {
			return fmt.Errorf("failed creating output file: %w", err)
		}
		defer f.Close()

		if !strings.Contains(filepath.Base(path), ".tmpl") {
			scaffoldFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer scaffoldFile.Close()
			_, err = io.Copy(f, scaffoldFile)
			if err != nil {
				return err
			}
			return nil
		}
		return template.Must(template.New(filepath.Base(path)).Funcs(template.FuncMap{
			"ToCamel": strcase.ToCamel,
		}).ParseFiles(path)).Execute(f, cfg)
	})
}
