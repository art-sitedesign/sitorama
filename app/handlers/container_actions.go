package handlers

import (
	"context"
	"html/template"
	"net/http"

	"github.com/art-sitedesign/sitorama/app/core"
)

func ContainerRestart(tmpl *template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		cr, err := core.NewCore()
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		cID := r.FormValue("cid")
		err = cr.ContainerRestart(context.Background(), cID)
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		http.Redirect(w, r, "/", 302)
	}
}

func ContainerStop(tmpl *template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		cr, err := core.NewCore()
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		cID := r.FormValue("cid")
		err = cr.ContainerStop(context.Background(), cID)
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		http.Redirect(w, r, "/", 302)
	}
}

func ContainerStart(tmpl *template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		cr, err := core.NewCore()
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		cID := r.FormValue("cid")
		err = cr.ContainerStart(context.Background(), cID)
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		http.Redirect(w, r, "/", 302)
	}
}

func ContainerRemove(tmpl *template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		cr, err := core.NewCore()
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		cID := r.FormValue("cid")
		err = cr.ContainerRemove(context.Background(), cID)
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		http.Redirect(w, r, "/", 302)
	}
}
