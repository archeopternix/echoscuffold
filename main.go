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
}

type Config struct {
	ApplicationPath string
}

func (c Config) CreateTargetApp() {
	/*_, err := os.Lstat(c.ApplicationPath + "/model")
	if err != nil {
		err = os.Mkdir(c.ApplicationPath+"/model", os.ModeDir)
		if err != nil {
			log.Fatalf("Subdir created: %s", err)
		}
	}
	*/
	_, err := os.Lstat(c.ApplicationPath + "/view")
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

var config Config

func identifyLookups(list []Entity) []Entity {
	_, es := getAllEntities()
	for _, entity := range es {
		for _, field := range entity.Fields {
			if field.Type == 4 {
				lk := NewEntity()
				lk.Name = field.Object
				lk.EntityType = 1
				lk.addField(Field{Name: "Text", Required: true, Type: 1})
				lk.addField(Field{Name: "Order", Type: 2})
				list = append(list, *lk)
			}
		}
	}
	return list
}

func main() {
	var err error
	var obj ObjectModel
	config.ApplicationPath = "/Users/A.Eisner/go/src/Test"
	config.CreateTargetApp()

	_, obj.Entities = getAllEntities()
	obj.Entities = identifyLookups(obj.Entities)
	fmt.Println(obj.Entities)
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

		output, err = os.Create(config.ApplicationPath + "/" + strings.ToLower(e.Name) + ".go")
		if err != nil {
			log.Fatalf("File creation: %s", err)
		}

		err = modeltmpl.ExecuteTemplate(output, "model", obj.Object)
		if err != nil {
			log.Fatalf("template execution: %s", err)
		}

		fmt.Println(output.Name())
	}

	// Copy the view components
	err = copyFile("template/base.html", config.ApplicationPath+"/view/base.html")
	if err != nil {
		log.Fatalf("Copy of base: %s", err)
	}
	err = copyFile("template/dashboard.html", config.ApplicationPath+"/view/dashboard.html")
	if err != nil {
		log.Fatalf("Copy of dashboard: %s", err)
	}
	err = copyFile("template/copy/_footer.html", config.ApplicationPath+"/view/_footer.html")
	if err != nil {
		log.Fatalf("Copy of _footer: %s", err)
	}
	err = copyFile("template/copy/_header.html", config.ApplicationPath+"/view/_header.html")
	if err != nil {
		log.Fatalf("Copy of _header: %s", err)
	}
	err = copyFile("template/copy/_hero.html", config.ApplicationPath+"/view/_hero.html")
	if err != nil {
		log.Fatalf("Copy of _hero: %s", err)
	}
	err = copyFile("template/copy/_mainnav.html", config.ApplicationPath+"/view/_mainnav.html")
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

	output, err = os.Create(config.ApplicationPath + "/view/sidenav.html")
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

		output, _ = os.Create(config.ApplicationPath + "/view/" + strings.ToLower(e.Name) + "list.html")

		err = listtmpl.ExecuteTemplate(output, "list", e)
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

		output, _ = os.Create(config.ApplicationPath + "/view/" + strings.ToLower(e.Name) + "detail.html")

		err = detailtmpl.ExecuteTemplate(output, "detail", e)
		if err != nil {
			log.Fatalf("template execution: %s", err)
		}
	}

	// Application classes:
	// Copy of database.go
	err = copyFile("database.go", config.ApplicationPath+"/database.go")
	if err != nil {
		log.Fatalf("Copy of .go files: %s", err)
	}

	// Rendering the main.go
	maintmpl, err := template.New("main").Funcs(funcMap).ParseFiles("template/main.tmpl")
	if err != nil {
		log.Fatalf("Parse main.go template: %s", err)
	}

	output, err = os.Create(config.ApplicationPath + "/main.go")
	if err != nil {
		log.Fatalf("File creation: %s", err)
	}
	err = maintmpl.ExecuteTemplate(output, "main", obj)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
	output.Close()

}
