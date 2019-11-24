package handlers

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/art-sitedesign/sitorama/app/core"
)

func Create(tmpl *template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		cr, err := core.NewCore()
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		err = cr.CreateSite(context.Background(), "test.loc")
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		http.Redirect(w, r, "/", 302)
	}
}
