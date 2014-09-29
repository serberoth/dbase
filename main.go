// The main package of the DBase sample RESTful database service
package main

import (
	// "fmt"
	"log"
	"net/http"
)

// The entry point for the application.  This method initialized the tables
// registers the HTTP service hooks and starts the HTTP server listening
// on port 8888.
func main() {
	// Initialize the tables from the persisted data
	if err := InitializeTables(); err != nil {
		log.Fatal(err)
	}

	defer PersistTables()

	// Initialize the routes
	http.HandleFunc("/tables", HandleTables)
	http.HandleFunc("/tables/", HandleTable)
	http.HandleFunc("/search", HandleSearch)
	http.HandleFunc("/search/", HandleSearch)

	http.HandleFunc("/styles.css", HandleStyles)

	// Initialize the root to redirect to the tables service
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/tables", http.StatusFound)
	})

	log.Fatal(http.ListenAndServe(":8888", nil))
}
