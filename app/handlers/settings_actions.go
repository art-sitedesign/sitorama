package handlers

import (
	"html/template"
	"net/http"

	"github.com/art-sitedesign/sitorama/app/core/settings"
)

func SettingsApp(tmpl *template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		appSettings, err := settings.NewApp()
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}

		if r.Method == "POST" {
			projectsRoot := r.FormValue("projectsRoot")
			appSettings.ProjectsRoot = projectsRoot
			err = appSettings.Save()
			if err != nil {
				writeErr(tmpl, w, err)
				return
			}

			http.Redirect(w, r, "/", 302)
			return
		}

		data := map[string]string{"ProjectsRoot": appSettings.ProjectsRoot}
		err = tmpl.ExecuteTemplate(w, "settings-app.html", data)
		if err != nil {
			writeErr(tmpl, w, err)
			return
		}
	}
}
