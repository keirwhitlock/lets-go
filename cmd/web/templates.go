package main

import (
	"html/template"
	"path/filepath"

	"github.com/kwhitlock/lets-go-book/pkg/models"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {

	// init new map as a cache
	cache := map[string]*template.Template{}

	// get a slice of all filepaths with *.page.tmpl
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// get filename to use as a map key
		name := filepath.Base(page)

		// parse the page template into a template set
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// parse the layout template into a template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}
		// parse the partial template into a template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
