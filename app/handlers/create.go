package handlers

import (
	"context"
	"html/template"
	"net/http"

	"github.com/art-sitedesign/sitorama/app/core"
)

func Create(tmpl *template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			name := r.PostFormValue("domain")

			cr, err := core.NewCore()
			if err != nil {
				writeErr(tmpl, w, err)
				return
			}

			err = cr.CreateProject(context.Background(), name)
			if err != nil {
				writeErr(tmpl, w, err)
				return
			}

			http.Redirect(w, r, "/", 302)
			return
		}

		data := struct{}{}
		err := tmpl.ExecuteTemplate(w, "create.html", data)
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}
	}
}
