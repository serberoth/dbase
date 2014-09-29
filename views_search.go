package main

import (
	"log"
	"net/http"
	"regexp"
)

// HandleSearch handles the /search service request.  This method handles
// both GET and POST for the search requests.
func HandleSearch(writer http.ResponseWriter, request *http.Request) {
	log.Println("Requested URL: ", request.URL)

	// Get the query parameter
	q := request.FormValue("q")
	log.Println("Query: ", q)

	// Handle GET and POST methods
	if "GET" == request.Method || "POST" == request.Method {
		// Render no search results for a blank query string
		if "" == q {
			RenderTemplate(writer, "search_view.html", Search{})
			return
		}

		// Compile the regular expression and return an error
		// response if the compilation failed
		regex, err := regexp.Compile(q)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		search := Search{Query: q, Results: make([]EntryInfo, 0)}

		// Search the tables for keys that match the given regex
		for _, table := range Tables {
			for name, entry := range table.Data {
				if match := regex.FindStringSubmatch(name); match != nil {
					search.Results = append(search.Results, EntryInfo{TableName: table.Name, EntryName: name, Value: entry})
				}
			}
		}

		// Render the search template with the results
		RenderTemplate(writer, "search_view.html", search)
	} else {
		// Return an invalid request response
		http.Error(writer, "Invalid request", http.StatusBadRequest)
	}
}
