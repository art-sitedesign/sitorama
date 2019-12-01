package handlers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/art-sitedesign/sitorama/app/core"
	"github.com/art-sitedesign/sitorama/app/core/builder"
	"github.com/art-sitedesign/sitorama/app/models"
)

func ProjectCreate(tmpl *template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			d := r.PostFormValue("domain")
			ep := r.PostFormValue("entryPoint")
			ws, _ := strconv.Atoi(r.PostFormValue("webServer"))
			db, _ := strconv.Atoi(r.PostFormValue("database"))
			c, _ := strconv.Atoi(r.PostFormValue("cache"))
			model := models.NewProjectCreate(d, ep, ws, db, c)

			cr, err := core.NewCore()
			if err != nil {
				writeErr(tmpl, w, err)
				return
			}

			builders := cr.CreateBuilders(model)

			data := make(map[string]builder.Config)
			for _, b := range builders {
				conf, err := b.PrepareConfig()
				if err != nil {
					writeErr(tmpl, w, err)
					return
				}
				data[b.Name()] = conf
			}

			err = tmpl.ExecuteTemplate(w, "confirm.html", map[string]interface{}{"Model": model, "Config": data})
			if err != nil {
				writeErr(tmpl, w, err)
				return
			}

			http.Redirect(w, r, "/", 302)
			return
		}

		data := map[string]interface{}{
			"WebserverBuilders": builder.WebserverBuilders,
			"DatabaseBuilders":  builder.DatabaseBuilders,
			"CacheBuilders":     builder.CacheBuilders,
		}
		err := tmpl.ExecuteTemplate(w, "create.html", data)
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}
	}
}

func ProjectCreateConfirm(tmpl *template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		ctx := context.Background()

		form := r.PostForm
		d := form.Get("model[domain]")
		ep := form.Get("model[entryPoint]")
		ws, _ := strconv.Atoi(form.Get("model[webServer]"))
		db, _ := strconv.Atoi(form.Get("model[database]"))
		ch, _ := strconv.Atoi(form.Get("model[cache]"))
		model := models.NewProjectCreate(d, ep, ws, db, ch)

		cr, err := core.NewCore()
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		builders := cr.CreateBuilders(model)
		for _, b := range builders {
			builderConfig := make(builder.Config)
			configNames := b.ConfigNames()
			for _, cName := range configNames {
				key := fmt.Sprintf("config[%s][%s]", b.Name(), cName)
				if _, ok := form[key]; ok {
					builderConfig[cName] = form.Get(key)
				}
			}

			b.SetConfig(builderConfig)
		}

		err = cr.CreateProject(ctx, model, builders)
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}
	}
}

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
