// TemplateTest project main.go
package main

import (
	"fmt"
	"io"
	"os"
	"text/template"
)

type Metadata struct {
	AttrKey   string
	AttrValue string
	Content   map[string]string
}

type Page struct {
	Title        string
	TemplateName string
	Metadata
	Data interface{}
}

func NewPage(title string, template string) *Page {
	p := new(Page)
	p.Title = title
	p.TemplateName = template
	return p
}

func (p Page) RenderTemplate(w io.Writer, s ...string) {
	if len(s) > 0 {
		if err := tstore[p.TemplateName].ExecuteTemplate(w, s[0], p); err != nil {
			panic(err)
		}
	} else {
		if err := tstore[p.TemplateName].ExecuteTemplate(w, "base.html", p); err != nil {
			panic(err)
		}
	}
}

type TemplateStore map[string]*template.Template

type PageStore map[string]*Page

var tstore TemplateStore

func main() {
	tstore = make(TemplateStore)
	tstore["list"] = template.Must(template.ParseFiles("template/base.html", "template/list.html"))
	tstore["row"] = template.Must(template.ParseFiles("template/base.html", "template/row.html"))

	p := NewPage("Monday", "list")

	w, err := os.Create("test.txt")
	defer w.Close()
	if err != nil {
		panic(err)
	}

	p.RenderTemplate(w)
	fmt.Fprintln(w, "------------")
	p.RenderTemplate(w, "content")
	/*
		fmt.Fprintln(w, tstore["list"].DefinedTemplates())
		tstore["list"].ExecuteTemplate(w, "base.html", "")
		fmt.Fprintln(w, tstore["row"].DefinedTemplates())
		tstore["row"].ExecuteTemplate(w, "base.html", "")
	*/
}
