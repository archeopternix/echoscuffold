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

type ObjectModel struct {
	Entities []Entity
	Object   Entity
}

func main() {
	var err error
	var obj ObjectModel

	_, obj.Entities = getAllEntities()

	pl := pluralize.NewClient()

	// First we create a FuncMap with which to register the function.
	funcMap := template.FuncMap{
		"lowercase": strings.ToLower, "singular": pl.Singular, "plural": pl.Plural,
	}

	// Rendering the model components
	modeltmpl, err := template.New("model").Funcs(funcMap).ParseFiles("template/model.tmpl", "template/types.tmpl")
	for _, e := range obj.Entities {
		var output *os.File
		defer output.Close()
		obj.Object = e
		output, _ = os.Create("model/" + e.Name + ".go")

		err = modeltmpl.ExecuteTemplate(output, "model", obj.Object)
		if err != nil {
			log.Fatalf("template execution: %s", err)
		}
	}

	/*
		// Rendering the view components
		listtmpl, _ := template.New("base").Funcs(funcMap).ParseFiles("template/base.html", "template/list.html")

		if err != nil {
			log.Fatalf("template execution: %s", err)
		}

		for _, e := range entities {
			var output *os.File
			defer output.Close()

			output, _ = os.Create("view/" + strings.ToLower(e.Name) + "list.html")

			err = listtmpl.ExecuteTemplate(output, "base", e)
			if err != nil {
				log.Fatalf("template execution: %s", err)
			}
		}
	*/
}
