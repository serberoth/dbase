package main

import (
	"html/template"
	"net/http"
)

// RenderTemplate is responsible for rendering the requested HTML template using the
// standard Go lang HTML template library with the provided data
func RenderTemplate(writer http.ResponseWriter, templateName string, data interface{}) {
	// Parse the HTML template and return an error response when the parse fails
	tmpl, err := template.ParseFiles(templateName)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template and return an error status with the render fails
	if err = tmpl.Execute(writer, data); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
