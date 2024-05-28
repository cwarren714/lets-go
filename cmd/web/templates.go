package main

import (
	"chandler.letsgo/internal/models"
	"html/template"
	"path/filepath"
	"time"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
// to it as the build progresses.
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
	Flash       string
}

// humanDate returns a human-readable version of a time.Time object
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// custom template functions can return only one value
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {

	// new map to act as cache
	cache := map[string]*template.Template{}

	// glob returns a slice of all pages that match the pattern
	pages, err := filepath.Glob("./ui/html/pages/*")
	if err != nil {
		return nil, err
	}

	// loop through the pages
	for _, page := range pages {

		// extract the file name from the full path
		name := filepath.Base(page)

		// parse base template
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// parse partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// add the template set to the cache,
		// using the name of the page as the key
		cache[name] = ts
	}

	return cache, nil
}
