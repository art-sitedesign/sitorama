package handlers

import (
	"context"
	"html/template"
	"net/http"

	"github.com/art-sitedesign/sitorama/app/core"
)

func Init(tmpl *template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		cr, err := core.NewCore()
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		err = cr.Init(context.Background())
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		http.Redirect(w, r, "/", 302)
	}
}
