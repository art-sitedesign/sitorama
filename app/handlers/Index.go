package handlers

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/art-sitedesign/sitorama/app/core"
)

func Index(tmpl *template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		cr, err := core.NewCore()
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		data := make(map[string]interface{})
		data["State"], err = cr.State(context.Background())
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		err = tmpl.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			_, _ = w.Write([]byte("Error!"))
		}
	}
}
