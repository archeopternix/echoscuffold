package controller

import (
	"net/http"

	"github.com/labstack/echo"
)

type Page struct {
	Title  string            // Page title within the browser
	Slug   string            // is the plural of the entities (for breadcrumb and to mark active menu entry)
	Data   interface{}       // the data that will be handed over to the templates
	Errors map[string]string // holds the errors on a page
}

func NewPage(title string, slug string) *Page {
	p := &Page{Title: title, Slug: slug}
	return p
}

func Dashboard(c echo.Context) error {
	p := NewPage("Dashboard", "Dashboard")
	return c.Render(http.StatusOK, "dashboard", p)
}
