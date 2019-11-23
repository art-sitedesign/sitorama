package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"pet-projects/sitorama/app/handlers"
)

func main() {
	tmpl := template.Must(template.ParseFiles("app/templates/index.html", "app/templates/test.html"))

	http.HandleFunc("/", handlers.Index(tmpl))
	http.HandleFunc("/init", handlers.Init(tmpl))

	fmt.Println("Open GUI: http://127.0.0.1:8085")
	err := http.ListenAndServe(":8085", nil)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return
}
