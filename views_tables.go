package main

import (
	"log"
	"net/http"
)

func HandleTables(writer http.ResponseWriter, request *http.Request) {
	log.Println("Requested URL[", request.Method, "]: ", request.URL.Path)

	keys := make([]string, 0, len(Tables))
	for key := range Tables {
		keys = append(keys, key)
	}

	if "GET" == request.Method {
		RenderTemplate(writer, "tables_view.html", keys)
	} else if "DELETE" == request.Method {
		Tables = make(map[string]Table)
		// PersistTables()
		http.Redirect(writer, request, "/tables", http.StatusFound)
	} else {
		http.Error(writer, "Invalid request", http.StatusBadRequest)
	}
}
