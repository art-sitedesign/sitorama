package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func writeErr(tmpl *template.Template, w http.ResponseWriter, err error) {
	errText := fmt.Sprintf("error: %v", err)

	execErr := tmpl.ExecuteTemplate(w, "error.html", map[string]string{"Error": errText})
	if execErr != nil {
		log.Fatalf("error execute template: %v", execErr)
	}

	log.Println(errText)
}
