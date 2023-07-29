package views

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/Masterminds/sprig/v3"
	"github.com/gin-gonic/gin/render"
)

//go:embed *.tmpl
var views embed.FS

type instance struct {
	Template *template.Template
	Data     any
}

func (i instance) Instance(name string, data any) render.Render {
	panic("unreachable")
}

func (i instance) Render(w http.ResponseWriter) error {
	i.WriteContentType(w)
	return i.Template.Execute(w, i.Data)
}

func (i instance) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{"text/html; charset=utf-8"}
	}
}

type renderer struct {
	templates map[string]*template.Template
}

func (r renderer) Instance(name string, data any) render.Render {
	t, ok := r.templates[name]
	if !ok {
		panic(fmt.Errorf("template '%s' not found", name))
	}
	return instance{Template: t, Data: data}
}

func Renderer() (render.HTMLRender, error) {
	r := renderer{templates: make(map[string]*template.Template)}
	tmplFiles, err := fs.ReadDir(views, ".")
	if err != nil {
		return nil, err
	}
	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}
		pt, err := template.
			New(tmpl.Name()).
			Funcs(sprig.FuncMap()).
			Funcs(template.FuncMap{
				"embedHTML": func(filename string) template.HTML {
					d, err := public.ReadFile(fmt.Sprintf("public/%s", filename))
					if err != nil {
						return template.HTML(err.Error()) // TODO: is there a better way?
					}
					return template.HTML(d)
				},
			}).
			ParseFS(views, tmpl.Name(), "_*.tmpl")
		if err != nil {
			return nil, err
		}
		r.templates[tmpl.Name()] = pt
	}
	return r, nil
}

func MustRenderer() render.HTMLRender {
	r, err := Renderer()
	if err != nil {
		panic(fmt.Errorf("failed to created renderer: %w", err))
	}
	return r
}
