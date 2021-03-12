// AccountManager project main.go
package main

import (
	. "echoscuffold/model"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"text/template"

	. "github.com/archeopternix/filegenerator"
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

func (app *AppModel) GenerateModel() (*Engine, error) {
	e := NewEngine("Models")
	// generate directory structure
	dirs := NewDirectoryGenerator()
	if err := dirs.Add(filepath.Join(app.Path, "model")); err != nil {
		return nil, err
	}
	if err := e.AddGenerator(dirs); err != nil {
		return nil, err
	}

	// generate model files
	// pluralize and singularize functions for templates
	pl := pluralize.NewClient()
	// First we create a FuncMap with which to register the function.
	funcMap := template.FuncMap{
		"lowercase": strings.ToLower, "singular": pl.Singular, "plural": pl.Plural,
	}
	tg := NewTemplateGenerator(funcMap)
	if err := tg.Add(tg.Template("model").ParseFiles(filepath.Join("template", "model.tmpl"), filepath.Join("template", "types.tmpl"))); err != nil {
		return nil, (err)
	}

	for _, entity := range app.Entities {
		app.Object = entity

		// add output
		file := filepath.Join(app.Path, "model", strings.ToLower(entity.Name)) + ".go"
		if err := tg.ParseFilename("model", file, entity); err != nil {
			return nil, err
		}
	}
	if err := e.AddGenerator(tg); err != nil {
		return nil, err
	}

	return e, nil
}

func (a *AppModel) GenerateJSONDataStore() (*Engine, error) {
	e := NewEngine("Datastore")

	// create data directory
	dirs := NewDirectoryGenerator()
	if err := dirs.Add(filepath.Join(app.Path, "data")); err != nil {
		return nil, err
	}
	if err := e.AddGenerator(dirs); err != nil {
		return nil, err
	}

	// Copy of database.go
	cp := NewCopyGenerator()
	if err := cp.Add(filepath.Join("model", "database.go"), filepath.Join(app.Path, "model")); err != nil {
		return nil, err
	}
	if err := e.AddGenerator(cp); err != nil {
		return nil, err
	}

	return e, nil
}

func (a *AppModel) GenerateMainApplication() (*Engine, error) {
	e := NewEngine("Mainapp")

	// create main directory
	dirs := NewDirectoryGenerator()
	if err := dirs.Add(filepath.Join(app.Path, "cmd")); err != nil {
		return nil, err
	}
	if err := e.AddGenerator(dirs); err != nil {
		return nil, err
	}

	// generate model files
	// pluralize and singularize functions for templates
	pl := pluralize.NewClient()
	// First we create a FuncMap with which to register the function.
	funcMap := template.FuncMap{
		"lowercase": strings.ToLower, "singular": pl.Singular, "plural": pl.Plural,
	}
	tg := NewTemplateGenerator(funcMap)
	if err := tg.Add(tg.Template("main").ParseFiles(filepath.Join("template", "main.tmpl"))); err != nil {
		return nil, (err)
	}

	file := filepath.Join(app.Path, "cmd", "main.go")
	if err := tg.ParseFilename("main", file, app); err != nil {
		return nil, err
	}

	if err := e.AddGenerator(tg); err != nil {
		return nil, err
	}

	return e, nil
}

func (a *AppModel) GenerateController() (*Engine, error) {
	e := NewEngine("Controllers")

	// create main directory
	dirs := NewDirectoryGenerator()
	if err := dirs.Add(filepath.Join(app.Path, "controller")); err != nil {
		return nil, err
	}
	if err := e.AddGenerator(dirs); err != nil {
		return nil, err
	}

	cp := NewCopyGenerator()
	if err := cp.Add(filepath.Join("template", "page.go"), filepath.Join(app.Path, "controller")); err != nil {
		return nil, err
	}
	if err := e.AddGenerator(cp); err != nil {
		return nil, err
	}

	// generate model files
	// pluralize and singularize functions for templates
	pl := pluralize.NewClient()
	// First we create a FuncMap with which to register the function.
	funcMap := template.FuncMap{
		"lowercase": strings.ToLower, "singular": pl.Singular, "plural": pl.Plural,
	}
	tg := NewTemplateGenerator(funcMap)
	if err := tg.Add(tg.Template("controller").ParseFiles(filepath.Join("template", "controller.tmpl"))); err != nil {
		return nil, (err)
	}

	for _, entity := range app.Entities {
		app.Object = entity
		data := struct {
			App   AppModel
			Model Entity
		}{app, entity}

		// add output
		file := filepath.Join(app.Path, "controller", strings.ToLower(entity.Name)) + ".go"
		if err := tg.ParseFilename("controller", file, data); err != nil {
			return nil, err
		}
	}
	if err := e.AddGenerator(tg); err != nil {
		return nil, err
	}

	return e, nil
}

//-------
func (app *AppModel) GenerateView() (*Engine, error) {
	e := NewEngine("Views")

	// create view directory
	dirs := NewDirectoryGenerator()
	if err := dirs.Add(filepath.Join(app.Path, "view")); err != nil {
		return nil, err
	}
	if err := e.AddGenerator(dirs); err != nil {
		return nil, err
	}

	// Copy of files
	cp := NewCopyGenerator()
	if err := cp.Add(filepath.Join("template", "_base.html"), filepath.Join(app.Path, "view")); err != nil {
		return nil, err
	}
	if err := cp.Add(filepath.Join("template", "_dashboard.html"), filepath.Join(app.Path, "view")); err != nil {
		return nil, err
	}
	if err := cp.Add(filepath.Join("template", "copy", "_footer.html"), filepath.Join(app.Path, "view")); err != nil {
		return nil, err
	}
	if err := cp.Add(filepath.Join("template", "copy", "_header.html"), filepath.Join(app.Path, "view")); err != nil {
		return nil, err
	}
	if err := cp.Add(filepath.Join("template", "copy", "_hero.html"), filepath.Join(app.Path, "view")); err != nil {
		return nil, err
	}
	if err := cp.Add(filepath.Join("template", "copy", "_mainnav.html"), filepath.Join(app.Path, "view")); err != nil {
		return nil, err
	}
	if err := e.AddGenerator(cp); err != nil {
		return nil, err
	}

	// generate model files
	// pluralize and singularize functions for templates
	pl := pluralize.NewClient()
	// First we create a FuncMap with which to register the function.
	funcMap := template.FuncMap{
		"lowercase": strings.ToLower, "singular": pl.Singular, "plural": pl.Plural,
	}
	tg := NewTemplateGenerator(funcMap)

	// SideNav
	if err := tg.Add(tg.Template("sidenav").ParseFiles(filepath.Join("template", "_sidenav.html"))); err != nil {
		return nil, (err)
	}
	file := filepath.Join(app.Path, "view", "_sidenav.html")
	if err := tg.ParseFilename("sidenav", file, app); err != nil {
		return nil, err
	}

	// Dashboard
	if err := tg.Add(tg.Template("dashboard").ParseFiles(filepath.Join("template", "_dashboard.html"))); err != nil {
		return nil, (err)
	}
	file = filepath.Join(app.Path, "view", "_dashboard.html")
	if err := tg.ParseFilename("dashboard", file, app); err != nil {
		return nil, err
	}

	// List Views
	if err := tg.Add(tg.Template("list").ParseFiles(filepath.Join("template", "list.html"))); err != nil {
		return nil, (err)
	}
	for _, entity := range app.Entities {
		// add output
		file := filepath.Join(app.Path, "view", strings.ToLower(entity.Name)) + "list.html"
		if err := tg.ParseFilename("list", file, entity); err != nil {
			return nil, err
		}
	}

	// Listtable Views
	if err := tg.Add(tg.Template("listtable").ParseFiles(filepath.Join("template", "listtable.html"))); err != nil {
		return nil, (err)
	}
	for _, entity := range app.Entities {
		app.Object = entity

		// add output
		file := filepath.Join(app.Path, "view", strings.ToLower(entity.Name)) + "listtable.html"
		if err := tg.ParseFilename("listtable", file, entity); err != nil {
			return nil, err
		}
	}

	// Detail Views
	if err := tg.Add(tg.Template("detail").ParseFiles(filepath.Join("template", "detail.html"))); err != nil {
		return nil, (err)
	}
	for _, entity := range app.Entities {
		app.Object = entity

		// add output
		file := filepath.Join(app.Path, "view", strings.ToLower(entity.Name)) + "detail.html"
		if err := tg.ParseFilename("detail", file, entity); err != nil {
			return nil, err
		}
	}

	if err := e.AddGenerator(tg); err != nil {
		return nil, err
	}

	return e, nil
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

var app AppModel

func main() {

	app.Name = "XXX"
	app.Title = "Usermanagement for eTracker Accounts"
	app.Path = filepath.Join("\\Users", "A.Eisner", "go", "src", app.Name)
	//	app.CreateTargetApp()

	_, app.Entities = GetAllEntities()
	app.Entities = append(app.Entities, identifyLookups(app.Entities)...)
	parseRelations(app.Entities)

	// TODO: many to many relations with mapping objects missing

	repo := NewYAMLRepository("repo.yaml", app.Entities)
	repo.Save()

	// main application
	e, err := app.GenerateMainApplication()
	if err != nil {
		log.Fatal(err)
	}
	if err := e.Run(); err != nil {
		log.Fatal(err)
	}

	e, err = app.GenerateModel()
	if err != nil {
		log.Fatal(err)
	}
	if err := e.Run(); err != nil {
		log.Fatal(err)
	}

	e, err = app.GenerateJSONDataStore()
	if err != nil {
		log.Fatal(err)
	}
	if err := e.Run(); err != nil {
		log.Fatal(err)
	}

	e, err = app.GenerateController()
	if err != nil {
		log.Fatal(err)
	}
	if err := e.Run(); err != nil {
		log.Fatal(err)
	}

	e, err = app.GenerateView()
	if err != nil {
		log.Fatal(err)
	}
	if err := e.Run(); err != nil {
		log.Fatal(err)
	}
}
