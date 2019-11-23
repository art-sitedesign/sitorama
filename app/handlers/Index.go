package handlers

import (
	"html/template"
	"net/http"
)

func Index(tmpl *template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct{}{}
		err := tmpl.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			_, _ = w.Write([]byte("Error!"))
		}
	}
}
