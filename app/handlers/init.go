package handlers

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"pet-projects/sitorama/app/core"
)

func Init(tmpl *template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		cr, err := core.NewCore()
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		err = cr.Init(context.Background())
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		http.Redirect(w, r, "/", 301)
	}
}
