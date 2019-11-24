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
		data := struct{}{}
		err := tmpl.ExecuteTemplate(w, "create.html", data)
		if err != nil {
			_, _ = w.Write([]byte("Error!"))
		}
	}
}

func CreateSave(tmpl *template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PostFormValue("domain")

		cr, err := core.NewCore()
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		err = cr.CreateSite(context.Background(), name)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		http.Redirect(w, r, "/", 302)
	}
}
