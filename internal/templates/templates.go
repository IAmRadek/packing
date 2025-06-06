package templates

import (
	"embed"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

//go:embed layouts/* pages/*
var fs embed.FS

type Templates struct {
	templates map[string]*template.Template
}

func NewTemplates() (*Templates, error) {
	// TODO: possible improvement to not use embed.FS when in development environment.
	cache := map[string]*template.Template{}

	// Read all template files from the embedded filesystem
	pages, err := fs.ReadDir("pages")
	if err != nil {
		return nil, err
	}

	// Read the base template
	baseContent, err := fs.ReadFile("layouts/base.gotempl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		if !page.IsDir() && strings.HasSuffix(page.Name(), ".gotempl") {
			pageContent, err := fs.ReadFile(filepath.Join("pages", page.Name()))
			if err != nil {
				return nil, err
			}

			tpl, err := template.New("").Parse(string(baseContent))
			if err != nil {
				return nil, err
			}

			tpl, err = tpl.Parse(string(pageContent))
			if err != nil {
				return nil, err
			}

			cache[page.Name()] = tpl
		}
	}

	return &Templates{
		templates: cache,
	}, nil
}

func (t *Templates) Render(w http.ResponseWriter, r *http.Request, page string, data any) {
	tt, ok := t.templates[page+".gotempl"]
	if !ok {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tt.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *Templates) Templates() []string {
	out := make([]string, 0, len(t.templates))
	for page := range t.templates {
		out = append(out, page)
	}
	return out
}
