// AccountManager project main.go
package main

import (
	. "echoscuffold/model"
	"fmt"
	"io"
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

// Holds all entitites or one dedicated Object for template generation
type AppModel struct {
	Entities []Entity
	Object   Entity
	Config
}

// App configurtations
type Config struct {
	Path  string
	Name  string
	Title string
}

// Creates standard app directories
func (c Config) CreateTargetApp() {
	_, err := os.Lstat(c.Path)
	if err != nil {
		err = os.Mkdir(c.Path, os.ModeDir)
		if err != nil {
			log.Fatalf("App directory: %s", err)
		}
	}

	_, err = os.Lstat(c.Path + "/static")
	if err != nil {
		err = os.Mkdir(c.Path+"/static", os.ModeDir)
		if err != nil {
			log.Fatalf("Static files directory: %s", err)
		}
	}
	_, err = os.Lstat(c.Path + "/view")
	if err != nil {
		err = os.Mkdir(c.Path+"/view", os.ModeDir)
		if err != nil {
			log.Fatalf("Subdir view: %s", err)
		}
	}
	_, err = os.Lstat(c.Path + "/data")
	if err != nil {
		err = os.Mkdir(c.Path+"/data", os.ModeDir)
		if err != nil {
			log.Fatalf("Subdir data: %s", err)
		}
	}
	_, err = os.Lstat(c.Path + "/model")
	if err != nil {
		err = os.Mkdir(c.Path+"/model", os.ModeDir)
		if err != nil {
			log.Fatalf("Subdir model: %s", err)
		}
	}
	_, err = os.Lstat(c.Path + "/controller")
	if err != nil {
		err = os.Mkdir(c.Path+"/controller", os.ModeDir)
		if err != nil {
			log.Fatalf("Subdir controller: %s", err)
		}
	}
	log.Println("Directory structure generated or verified")
}

// copies file from "src" to "dst"
func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

// loops through all fields for "lookup" fields and adds a lookup entity
func identifyLookups(list []Entity) (entities []Entity) {
	lookups := make(map[string]Entity)
	_, es := GetAllEntities()
	for _, entity := range es {
		for _, field := range entity.Fields {
			if field.Type == "lookup" {
				lk := Entity{
					Name:       field.Object,
					EntityType: 1,
				}
				lk.AddField(Field{Name: "Text", Required: true, Type: "string"})
				lk.AddField(Field{Name: "Order", Type: "int"})
				lookups[lk.Name] = lk
			}
		}
	}

	for _, val := range lookups {
		entities = append(entities, val)
	}
	return entities
}

// loops through all relations and adds parent/child fields or many-to-many mappingttable
func parseRelations(list []Entity) {
	rels := GetAllRelations()
	fmt.Printf("%d relations loaded.\n", len(rels))
	for _, relation := range rels {
		if relation.Kind == "one2many" {
			for i, entity := range list {
				if relation.Child == entity.Name {
					list[i].AddField(Field{Name: relation.Parent, Type: "child", Object: relation.Parent})
				}
				if relation.Parent == entity.Name {
					list[i].AddField(Field{Name: relation.Child, Type: "parent", Object: relation.Child})
				}
			}
		}

		// TODO: many2many
		/*
			if relation.Kind == "many2many" {
				lk := NewEntity()
				lk.Name = relation.Parent + relation.Child
				lk.EntityType = 2
				lk.AddField(Field{Name: relation.Child, Type: "manychild", Object: lk.Name})
				lk.AddField(Field{Name: relation.Parent, Type: "manychild", Object: lk.Name})
				list = append(list, *lk)
				for i, entity := range list {
					if relation.Child == entity.Name {
						list[i].AddField(Field{Name: relation.Parent, Type: "manyparent", Object: lk.Name})
					}
					if relation.Parent == entity.Name {
						list[i].AddField(Field{Name: relation.Child, Type: "manyparent", Object: lk.Name})
					}
				}
			} */
	}
}

// generate app models and copies basic database functions
func generateModel(mt *template.Template) {
	var err error

	for _, e := range app.Entities {
		var output *os.File
		defer output.Close()
		app.Object = e
		output, err = os.Create(app.Path + "/model/" + strings.ToLower(e.Name) + ".go")
		if err != nil {
			log.Fatalf("File creation: %s", err)
		}

		err = mt.ExecuteTemplate(output, "model", app)
		if err != nil {
			log.Fatalf("template execution: %s", err)
		}

		fmt.Println("model generated: " + output.Name())
	}

	// Copy of database.go
	err = copyFile("model/database.go", app.Path+"/model/database.go")
	if err != nil {
		log.Fatalf("Copy of .go file: %s", err)
	}
	fmt.Println("models finished")
}

// generate app controller and copies basic functions
func generateController(ct *template.Template) {
	var err error

	for _, e := range app.Entities {
		var output *os.File
		defer output.Close()
		app.Object = e
		output, err = os.Create(app.Path + "/controller/" + strings.ToLower(e.Name) + ".go")
		if err != nil {
			log.Fatalf("File creation: %s", err)
		}

		err = ct.ExecuteTemplate(output, "controller", app)
		if err != nil {
			log.Fatalf("template execution: %s", err)
		}

		fmt.Println("controller generated: " + output.Name())
	}

	// Copy of dashboard.go
	err = copyFile("template/dashboard.tmpl", app.Path+"/controller/dashboard.go")
	if err != nil {
		log.Fatalf("Copy of .go file: %s", err)
	}
	fmt.Println("controllers finished")
}

var app AppModel

func main() {
	var err error

	app.Name = "ProjectMgnt"
	app.Title = "Usermanagement for eTracker Accounts"
	app.Path = "/Users/Andreas Eisner/go/src/" + app.Name
	app.CreateTargetApp()

	_, app.Entities = GetAllEntities()
	app.Entities = append(app.Entities, identifyLookups(app.Entities)...)
	parseRelations(app.Entities)

	// TODO: many to many relations with mapping objects missing

	repo := NewYAMLRepository("repo.yaml", app.Entities)
	repo.Save()

	//	app.Entities = append(app.Entities, parseRelations(app.Entities)...)
	fmt.Printf("%d entities loaded.\n", len(app.Entities))

	// pluralize and sigularize functions for templates
	pl := pluralize.NewClient()
	// First we create a FuncMap with which to register the function.
	funcMap := template.FuncMap{
		"lowercase": strings.ToLower, "singular": pl.Singular, "plural": pl.Plural,
	}

	// Rendering the model components
	tmpl, err := template.New("model").Funcs(funcMap).ParseFiles("template/model.tmpl", "template/types.tmpl")
	if err != nil {
		log.Fatalf("Error in model templates: %s", err)
	}
	generateModel(tmpl)

	// Rendering the controller components
	tmpl, err = template.New("controller").Funcs(funcMap).ParseFiles("template/controller.tmpl")
	if err != nil {
		log.Fatalf("Error in controller templates: %s", err)
	}
	generateController(tmpl)

	fmt.Println("views")
	// Copy the view components
	err = copyFile("template/base.html", app.Path+"/view/_base.html")
	if err != nil {
		log.Fatalf("Copy of base: %s", err)
	}
	err = copyFile("template/dashboard.html", app.Path+"/view/dashboard.html")
	if err != nil {
		log.Fatalf("Copy of dashboard: %s", err)
	}
	err = copyFile("template/copy/_footer.html", app.Path+"/view/_footer.html")
	if err != nil {
		log.Fatalf("Copy of _footer: %s", err)
	}
	err = copyFile("template/copy/_header.html", app.Path+"/view/_header.html")
	if err != nil {
		log.Fatalf("Copy of _header: %s", err)
	}
	err = copyFile("template/copy/_hero.html", app.Path+"/view/_hero.html")
	if err != nil {
		log.Fatalf("Copy of _hero: %s", err)
	}
	err = copyFile("template/copy/_mainnav.html", app.Path+"/view/_mainnav.html")
	if err != nil {
		log.Fatalf("Copy of _mainnav: %s", err)
	}

	// Rendering the sidenav
	sidenavtmpl, err := template.New("sidenav").Funcs(funcMap).ParseFiles("template/sidenav.html")
	if err != nil {
		log.Fatalf("Parse sidenav template: %s", err)
	}
	var output *os.File
	defer output.Close()

	output, err = os.Create(app.Path + "/view/_sidenav.html")
	if err != nil {
		log.Fatalf("File creation: %s", err)
	}
	err = sidenavtmpl.ExecuteTemplate(output, "sidenav", app)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
	output.Close()

	// Rendering the dashboard
	dashtmpl, err := template.New("dashboard").Funcs(funcMap).ParseFiles("template/dashboard.html")
	if err != nil {
		log.Fatalf("Parse dashboard template: %s", err)
	}

	defer output.Close()

	output, err = os.Create(app.Path + "/view/dashboard.html")
	if err != nil {
		log.Fatalf("File creation: %s", err)
	}
	err = dashtmpl.ExecuteTemplate(output, "dashboard", app)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
	output.Close()

	// Rendering the list view components
	listtmpl, _ := template.New("list").Funcs(funcMap).ParseFiles("template/list.html")
	if err != nil {
		log.Fatalf("Parse list template: %s", err)
	}
	for _, e := range app.Entities {
		var output *os.File
		defer output.Close()

		output, _ = os.Create(app.Path + "/view/" + strings.ToLower(e.Name) + "list.html")

		err = listtmpl.ExecuteTemplate(output, "list", e)
		if err != nil {
			log.Fatalf("template execution: %s", err)
		}
	}

	// Rendering the listtable view components
	listtabletmpl, _ := template.New("list").Funcs(funcMap).ParseFiles("template/listtable.html")
	if err != nil {
		log.Fatalf("Parse list template: %s", err)
	}
	for _, e := range app.Entities {
		var output *os.File
		defer output.Close()

		output, _ = os.Create(app.Path + "/view/" + strings.ToLower(e.Name) + "listtable.html")

		err = listtabletmpl.ExecuteTemplate(output, "listtable", e)
		if err != nil {
			log.Fatalf("template execution: %s", err)
		}
	}

	// Rendering the detailview components
	detailtmpl, _ := template.New("detail").Funcs(funcMap).ParseFiles("template/detail.html")
	if err != nil {
		log.Fatalf("Parse detail template: %s", err)
	}
	for _, e := range app.Entities {
		var output *os.File
		defer output.Close()

		output, _ = os.Create(app.Path + "/view/" + strings.ToLower(e.Name) + "detail.html")

		err = detailtmpl.ExecuteTemplate(output, "detail", e)
		if err != nil {
			log.Fatalf("template execution: %s", err)
		}
	}

	// Application classes:

	// Rendering the main.go
	maintmpl, err := template.New("main").Funcs(funcMap).ParseFiles("template/main.tmpl")
	if err != nil {
		log.Fatalf("Parse main.go template: %s", err)
	}

	output, err = os.Create(app.Path + "/main.go")
	if err != nil {
		log.Fatalf("File creation: %s", err)
	}
	err = maintmpl.ExecuteTemplate(output, "main", app)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
	output.Close()

}
