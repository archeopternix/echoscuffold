// AccountManager project main.go
package main

import (
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

type ObjectModel struct {
	Entities []Entity
	Object   Entity
	Config
}

type Config struct {
	ApplicationPath string
	ApplicationName string
}

func (c Config) CreateTargetApp() {
	_, err := os.Lstat(c.ApplicationPath + "/static")
	if err != nil {
		err = os.Mkdir(c.ApplicationPath+"/static", os.ModeDir)
		if err != nil {
			log.Fatalf("Static files directory created: %s", err)
		}
	}
	_, err = os.Lstat(c.ApplicationPath + "/view")
	if err != nil {
		err = os.Mkdir(c.ApplicationPath+"/view", os.ModeDir)
		if err != nil {
			log.Fatalf("Subdir created: %s", err)
		}
	}
	_, err = os.Lstat(c.ApplicationPath + "/data")
	if err != nil {
		err = os.Mkdir(c.ApplicationPath+"/data", os.ModeDir)
		if err != nil {
			log.Fatalf("Subdir created: %s", err)
		}
	}
	_, err = os.Lstat(c.ApplicationPath + "/model")
	if err != nil {
		err = os.Mkdir(c.ApplicationPath+"/model", os.ModeDir)
		if err != nil {
			log.Fatalf("Subdir created: %s", err)
		}
	}
	_, err = os.Lstat(c.ApplicationPath + "/controller")
	if err != nil {
		err = os.Mkdir(c.ApplicationPath+"/controller", os.ModeDir)
		if err != nil {
			log.Fatalf("Subdir created: %s", err)
		}
	}
}

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

func identifyLookups(list []Entity) []Entity {
	_, es := getAllEntities()
	for _, entity := range es {
		for _, field := range entity.Fields {
			if field.Type == "lookup" {
				lk := NewEntity()
				lk.Name = field.Object
				lk.EntityType = 1
				lk.addField(Field{Name: "Text", Required: true, Type: "string"})
				lk.addField(Field{Name: "Order", Type: "int"})
				list = append(list, *lk)
			}
		}
	}
	return list
}

func parseRelations(list []Entity) []Entity {
	rels := getAllRelations()
	fmt.Print("Relationen:")
	fmt.Println(rels)
	for _, relation := range rels {
		if relation.Kind == "one2many" {
			for i, entity := range list {
				if relation.Child == entity.Name {
					list[i].addField(Field{Name: relation.Parent, Type: "child", Object: relation.Parent})
				}
				if relation.Parent == entity.Name {
					list[i].addField(Field{Name: relation.Child, Type: "parent", Object: relation.Child})
				}
			}
		}
		if relation.Kind == "many2many" {
			lk := NewEntity()
			lk.Name = relation.Parent + relation.Child
			lk.EntityType = 2
			lk.addField(Field{Name: relation.Child, Type: "manychild", Object: lk.Name})
			lk.addField(Field{Name: relation.Parent, Type: "manychild", Object: lk.Name})
			list = append(list, *lk)
			for i, entity := range list {
				if relation.Child == entity.Name {
					list[i].addField(Field{Name: relation.Parent, Type: "manyparent", Object: lk.Name})
				}
				if relation.Parent == entity.Name {
					list[i].addField(Field{Name: relation.Child, Type: "manyparent", Object: lk.Name})
				}
			}
		}
	}
	return list
}

func main() {
	var err error
	var obj ObjectModel
	obj.ApplicationName = "CRUD"
	obj.ApplicationPath = "/Users/A.Eisner/go/src/" + obj.ApplicationName
	obj.CreateTargetApp()

	_, obj.Entities = getAllEntities()
	fmt.Println(len(obj.Entities))
	obj.Entities = identifyLookups(obj.Entities)
	obj.Entities = parseRelations(obj.Entities)
	fmt.Println(obj)

	pl := pluralize.NewClient()

	// First we create a FuncMap with which to register the function.
	funcMap := template.FuncMap{
		"lowercase": strings.ToLower, "singular": pl.Singular, "plural": pl.Plural,
	}

	fmt.Println("models")
	// Rendering the model components
	modeltmpl, err := template.New("model").Funcs(funcMap).ParseFiles("template/model.tmpl", "template/types.tmpl")
	for _, e := range obj.Entities {
		var output *os.File
		defer output.Close()
		obj.Object = e

		output, err = os.Create(obj.ApplicationPath + "/" + strings.ToLower(e.Name) + ".go")
		if err != nil {
			log.Fatalf("File creation: %s", err)
		}

		err = modeltmpl.ExecuteTemplate(output, "model", obj.Object)
		if err != nil {
			log.Fatalf("template execution: %s", err)
		}

		fmt.Println(output.Name())
	}

	fmt.Println("views")
	// Copy the view components
	err = copyFile("template/base.html", obj.ApplicationPath+"/view/base.html")
	if err != nil {
		log.Fatalf("Copy of base: %s", err)
	}
	err = copyFile("template/dashboard.html", obj.ApplicationPath+"/view/dashboard.html")
	if err != nil {
		log.Fatalf("Copy of dashboard: %s", err)
	}
	err = copyFile("template/copy/_footer.html", obj.ApplicationPath+"/view/_footer.html")
	if err != nil {
		log.Fatalf("Copy of _footer: %s", err)
	}
	err = copyFile("template/copy/_header.html", obj.ApplicationPath+"/view/_header.html")
	if err != nil {
		log.Fatalf("Copy of _header: %s", err)
	}
	err = copyFile("template/copy/_hero.html", obj.ApplicationPath+"/view/_hero.html")
	if err != nil {
		log.Fatalf("Copy of _hero: %s", err)
	}
	err = copyFile("template/copy/_mainnav.html", obj.ApplicationPath+"/view/_mainnav.html")
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

	output, err = os.Create(obj.ApplicationPath + "/view/sidenav.html")
	if err != nil {
		log.Fatalf("File creation: %s", err)
	}
	err = sidenavtmpl.ExecuteTemplate(output, "sidenav", obj)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
	output.Close()

	// 		listtmpl, _ := template.New("list").Funcs(funcMap).ParseFiles("template/base.html", "template/list.html")

	// Rendering the listview components
	listtmpl, _ := template.New("list").Funcs(funcMap).ParseFiles("template/list.html")
	if err != nil {
		log.Fatalf("Parse list template: %s", err)
	}
	for _, e := range obj.Entities {
		var output *os.File
		defer output.Close()

		output, _ = os.Create(obj.ApplicationPath + "/view/" + strings.ToLower(e.Name) + "list.html")

		err = listtmpl.ExecuteTemplate(output, "list", e)
		if err != nil {
			log.Fatalf("template execution: %s", err)
		}
	}

	// Rendering the listview components
	listtabletmpl, _ := template.New("list").Funcs(funcMap).ParseFiles("template/listtable.html")
	if err != nil {
		log.Fatalf("Parse list template: %s", err)
	}
	for _, e := range obj.Entities {
		var output *os.File
		defer output.Close()

		output, _ = os.Create(obj.ApplicationPath + "/view/" + strings.ToLower(e.Name) + "listtable.html")

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
	for _, e := range obj.Entities {
		var output *os.File
		defer output.Close()

		output, _ = os.Create(obj.ApplicationPath + "/view/" + strings.ToLower(e.Name) + "detail.html")

		err = detailtmpl.ExecuteTemplate(output, "detail", e)
		if err != nil {
			log.Fatalf("template execution: %s", err)
		}
	}

	// Application classes:
	// Copy of database.go
	err = copyFile("database.go", obj.ApplicationPath+"/database.go")
	if err != nil {
		log.Fatalf("Copy of .go files: %s", err)
	}

	// Rendering the main.go
	maintmpl, err := template.New("main").Funcs(funcMap).ParseFiles("template/main.tmpl")
	if err != nil {
		log.Fatalf("Parse main.go template: %s", err)
	}

	output, err = os.Create(obj.ApplicationPath + "/main.go")
	if err != nil {
		log.Fatalf("File creation: %s", err)
	}
	err = maintmpl.ExecuteTemplate(output, "main", obj)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
	output.Close()

}
