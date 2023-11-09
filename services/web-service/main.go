package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type IndexPageViewModel struct {
	PdfTemplate []PdfTemplate
}

type PdfTemplate struct {
	UpdatedAt time.Time
	CreatedAt time.Time
	Name      string
}

func requestPdfTemplates() ([]PdfTemplate, error) {
	// define TEMPLATE_ENDPOINT in docker-compose.yml
	endpoint := fmt.Sprintf("http://%s/api/templates", os.Getenv("TEMPLATE_ENDPOINT"))
	log.Println(endpoint)

	req, _ := http.NewRequest("GET", endpoint, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	log.Println(res)
	log.Println(err)

	var templates []PdfTemplate
	if err := json.NewDecoder(res.Body).Decode(&templates); err != nil {
		return nil, err
	}

	return templates, nil
}

func indexHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templates, err := requestPdfTemplates()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		viewModel := IndexPageViewModel{
			PdfTemplate: templates,
		}

		tmpl.ExecuteTemplate(w, "index", viewModel)
	}
}

func main() {
	tmpl := template.Must(template.ParseGlob("templates/*.gohtml"))

	router := http.NewServeMux()
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))
	router.HandleFunc("/", indexHandler(tmpl))

	log.Fatal(http.ListenAndServe("0.0.0.0:3000", router))
}
