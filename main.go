// AccountManager project main.go
package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gertd/go-pluralize"
)

func main() {
	pl := pluralize.NewClient()

	// First we create a FuncMap with which to register the function.
	funcMap := template.FuncMap{
		"lowercase": strings.ToLower, "singular": pl.Singular, "plural": pl.Plural,
	}

	// pattern is the glob pattern used to find all the template files.
	pattern := filepath.Join("template", "*.tmpl")

	tmpl, err := template.New("model").Funcs(funcMap).ParseGlob(pattern)

	_, entities := getAllEntities()

	for _, e := range entities {
		var output *os.File
		defer output.Close()

		output, _ = os.Create("model/" + e.Name + ".go")

		err = tmpl.ExecuteTemplate(output, "model", e)
		if err != nil {
			log.Fatalf("template execution: %s", err)
		}
	}
}
