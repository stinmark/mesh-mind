package handlers

import (
	"gateway-go/internal/templates"
	"html/template"
	"net/http"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(templates.FS, "index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
