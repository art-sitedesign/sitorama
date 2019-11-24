package handlers

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/art-sitedesign/sitorama/app/core"
)

func ProjectStop(tmpl *template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		cr, err := core.NewCore()
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		name := r.FormValue("name")
		err = cr.StopProject(context.Background(), name)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		http.Redirect(w, r, "/", 302)
	}
}
