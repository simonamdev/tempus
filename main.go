package main

import (
	"errors"
	"io"
	"math"
	"net/http"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

// Define the template registry struct
type TemplateRegistry struct {
	templates map[string]*template.Template
}

// Implement e.Renderer interface
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}

	return tmpl.ExecuteTemplate(w, "base.html", data)
}

// func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
// 	return t.templates.ExecuteTemplate(w, name, data)
// }

func main() {

	funcs := make(map[string]interface{})
	funcs["isNil"] = func(value interface{}) bool { return value == nil }
	funcs["formatDate"] = func(t time.Time) string {
		return t.Format("Jan 2 15:04:05")
	}
	funcs["round"] = func(f float64) float64 {
		return math.Round(f*100) / 100
	}

	funcMap := template.FuncMap(funcs)

	// Echo instance
	e := echo.New()
	e.Debug = true

	parseFuncs := func() *template.Template {
		return template.New("").Funcs(funcMap)
	}

	// Ref: https://gist.github.com/rand99/808e6e9702c00ce64803d94abff65678
	templates := make(map[string]*template.Template)
	templates["index"] = template.Must(parseFuncs().ParseFiles("public/views/index.html", "public/views/base.html"))

	e.Renderer = &TemplateRegistry{
		templates: templates,
	}

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})
	e.Logger.Fatal(e.Start(":1323"))
}