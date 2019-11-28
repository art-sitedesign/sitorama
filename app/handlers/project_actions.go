package handlers

import (
	"context"
	"html/template"
	"net/http"

	"github.com/art-sitedesign/sitorama/app/core"
)

func ProjectStart(tmpl *template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		cr, err := core.NewCore()
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		name := r.FormValue("name")
		err = cr.StartProject(context.Background(), name)
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		http.Redirect(w, r, "/", 302)
	}
}

func ProjectStop(tmpl *template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		cr, err := core.NewCore()
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		name := r.FormValue("name")
		err = cr.StopProject(context.Background(), name)
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		http.Redirect(w, r, "/", 302)
	}
}

func ProjectRemove(tmpl *template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		cr, err := core.NewCore()
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		name := r.FormValue("name")
		err = cr.RemoveProject(context.Background(), name)
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		http.Redirect(w, r, "/", 302)
	}
}
