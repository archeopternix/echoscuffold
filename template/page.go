package controller

import (
	"net/http"

	"github.com/labstack/echo"
)

type Page struct {
	Title  string
	Slug   string
	Data   interface{}
	Errors map[string]string
}

func NewPage(title string, slug string) *Page {
	p := &Page{Title: title, Slug: slug}
	return p
}

func Dashboard(c echo.Context) error {
	p := NewPage("Dashboard", "Dashboard")
	return c.Render(http.StatusOK, "dashboard", p)
}
