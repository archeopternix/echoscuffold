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
	_, err := os.Lstat(c.ApplicationPath + "/model")
	if err != nil {
		err = os.Mkdir(c.ApplicationPath+"/model", os.ModeDir)
		if err != nil {
			log.Fatalf("Subdir created: %s", err)
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

func main() {
	var err error
	var obj ObjectModel
	config.ApplicationPath = "/Users/Andreas Eisner/go/src/testcrud"
	config.CreateTargetApp()

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

		output, err = os.Create(config.ApplicationPath + "/model/" + e.Name + ".go")
		if err != nil {
			log.Fatalf("File creation: %s", err)
		}

		err = modeltmpl.ExecuteTemplate(output, "model", obj.Object)
		if err != nil {
			log.Fatalf("template execution: %s", err)
		}

		fmt.Println(output.Name())
	}

	// Rendering the view components
	err = copyFile("template/base.html", config.ApplicationPath+"/view/base.html")
	if err != nil {
		log.Fatalf("Copy of base: %s", err)
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
	listtmpl, _ := template.New("content").Funcs(funcMap).ParseFiles("template/list.html")
	if err != nil {
		log.Fatalf("Parse list template: %s", err)
	}
	for _, e := range obj.Entities {
		var output *os.File
		defer output.Close()

		output, _ = os.Create(config.ApplicationPath + "/view/" + strings.ToLower(e.Name) + "list.html")

		err = listtmpl.ExecuteTemplate(output, "content", e)
		if err != nil {
			log.Fatalf("template execution: %s", err)
		}
	}

	// Application classes:
	// Copy of database.go
	err = copyFile("database.go", config.ApplicationPath+"database.go")
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
