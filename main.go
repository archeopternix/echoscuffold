// AccountManager project main.go
package main

import (
	"log"
	"os"

	"strings"
	"text/template"

	"github.com/gertd/go-pluralize"
)

type Page struct {
	Title      string
	Navigation interface{}
	Content    interface{}
}

func main() {
	var err error

	pl := pluralize.NewClient()

	// First we create a FuncMap with which to register the function.
	funcMap := template.FuncMap{
		"lowercase": strings.ToLower, "singular": pl.Singular, "plural": pl.Plural,
	}

	// Rendering the model components
	//	tmpl, err := template.New("model").Funcs(funcMap).ParseFiles("template/model.tmpl", "template/types.tmpl")

	_, entities := getAllEntities()
	/*
		for _, e := range entities {
			var output *os.File
			defer output.Close()

			output, _ = os.Create("model/" + e.Name + ".go")

			err = tmpl.ExecuteTemplate(output, "model", e)
			if err != nil {
				log.Fatalf("template execution: %s", err)
			}
		}
	*/
	// Rendering the view components
	tmpl, _ := template.New("base").Funcs(funcMap).ParseFiles("template/base.html", "template/list.html")

	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
	//	_, entities = getAllEntities()

	for _, e := range entities {
		var output *os.File
		defer output.Close()

		output, _ = os.Create("view/" + e.Name + ".html")

		err = tmpl.ExecuteTemplate(output, "base", e)
		if err != nil {
			log.Fatalf("template execution: %s", err)
		}
	}

}
