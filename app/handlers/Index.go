package handlers

import (
	"context"
	"html/template"
	"net/http"

	"github.com/art-sitedesign/sitorama/app/core"
	"github.com/art-sitedesign/sitorama/app/core/settings"
)

func Index(tmpl *template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		cr, err := core.NewCore()
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		appSettings, err := settings.NewApp()
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		data := make(map[string]interface{})
		data["AppSettings"] = appSettings
		data["State"], err = cr.State(context.Background())
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		err = tmpl.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}
	}
}
