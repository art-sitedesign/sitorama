package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/art-sitedesign/sitorama/app/handlers"
)

//todo: запилить нормальные ошибки и логирование
//todo: годоки дописать
//todo: конструктор приложения
//todo: логи сервисов прокидывать наружу или в stdout?

func main() {
	tmpl := template.Must(template.ParseFiles("app/templates/html/index.html"))

	http.HandleFunc("/", handlers.Index(tmpl))
	http.HandleFunc("/init", handlers.Init(tmpl))
	http.HandleFunc("/create", handlers.Create(tmpl))

	fmt.Println("Open GUI: http://127.0.0.1:8085")
	err := http.ListenAndServe(":8085", nil)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return
}
