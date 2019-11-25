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
//todo: дождаться выполнения контейнера

func main() {
	tmpl := template.Must(template.ParseFiles(
		"app/templates/html/index.html",
		"app/templates/html/create.html",
	))

	http.HandleFunc("/", handlers.Index(tmpl))
	http.HandleFunc("/init", handlers.Init(tmpl))
	http.HandleFunc("/create", handlers.Create(tmpl))

	http.HandleFunc("/create/save", handlers.CreateSave(tmpl))

	// экшены проектов
	http.HandleFunc("/project/start", handlers.ProjectStart(tmpl))
	http.HandleFunc("/project/stop", handlers.ProjectStop(tmpl))
	http.HandleFunc("/project/remove", handlers.ProjectRemove(tmpl))

	// экшены контейнеров
	http.HandleFunc("/container/restart", handlers.ContainerRestart(tmpl))
	http.HandleFunc("/container/stop", handlers.ContainerStop(tmpl))
	http.HandleFunc("/container/start", handlers.ContainerStart(tmpl))
	http.HandleFunc("/container/remove", handlers.ContainerRemove(tmpl))

	fmt.Println("Open GUI: http://127.0.0.1:8085")
	err := http.ListenAndServe(":8085", nil)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return
}
