package main

import (
	"log"
	"net/http"
	"regexp"
)

// HandleTable handles the /tables/:table(/:key) requests.  This
// function is responsible for handling table views as well as
// entry views.
func HandleTable(writer http.ResponseWriter, request *http.Request) {
	log.Println("Requested URL: ", request.URL.Path)

	// Compile a regex to match the URL data (:table and :key)
	regex, err := regexp.Compile("^/tables(/[^/]+)(/.*)?$")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	// Pull the possible table name and entry name from the URL path
	// removing any '/' character from the beginning of the string.
	tableName, entryName := "", ""
	if match := regex.FindStringSubmatch(request.URL.Path); match != nil {
		if len(match[2]) >= 1 {
			tableName, entryName = match[1][1:], match[2][1:]
		} else {
			tableName, entryName = match[1][1:], match[2]
		}
	}

	log.Println("Requested Table: ", tableName)
	log.Println("Requested Entry: ", entryName)

	// Handle either a table view or an entry view
	if entryName == "" {
		HandleTableView(tableName, writer, request)
	} else {
		HandleTableEntry(tableName, entryName, writer, request)
	}
}

// HandleTableView handles the /tables/:table service requests.
func HandleTableView(tableName string, writer http.ResponseWriter, request *http.Request) {
	// If the table name is blank and the request is a GET request
	// redirect to the /tables view otherwise display an error
	if tableName == "" {
		if "GET" == request.Method {
			http.Redirect(writer, request, "/tables", http.StatusFound)
		} else {
			http.Error(writer, "Invalid request", http.StatusBadRequest)
		}
		return
	}

	if "GET" == request.Method {
		// Handle GET requests by displaying the table contents
		// or a message stating the table does not exist.
		if table, exists := Tables[tableName]; exists {
			RenderTemplate(writer, "table_view.html", table)
		} else {
			RenderTemplate(writer, "table_none.html", tableName)
		}
	} else if "POST" == request.Method {
		// Handle POST request by updating the table contents
		// with the JSON object data in the POST body display
		// an error message if the decode failed.  The new table
		// content overrites any existing table content.
		table, err := DecodeTable(request.Body)
		if err != nil {
			http.Error(writer, "Server error", http.StatusInternalServerError)
		} else {
			table.Name = tableName
			Tables[tableName] = table
			PersistTables()
			http.Redirect(writer, request, "/tables/"+tableName, http.StatusFound)
		}
	} else if "DELETE" == request.Method {
		// Handle DELETE reques by removing the table entry from
		// the global table map.
		delete(Tables, tableName)
		PersistTables()
		http.Redirect(writer, request, "/tables", http.StatusFound)
	} else {
		// Display an invalid request message
		http.Error(writer, "Invalid request", http.StatusBadRequest)
	}
}

// HandleTableEntry handles the /tables/:table/:key service requests.
func HandleTableEntry(tableName string, entryName string, writer http.ResponseWriter, request *http.Request) {
	if "GET" == request.Method {
		// Display either the entry view or if the requested table does not exist
		// the non-existant table message.
		if table, exists := Tables[tableName]; exists {
			RenderTemplate(writer, "entry_view.html", EntryInfo{TableName: tableName, EntryName: entryName, Value: table.Data[entryName]})
		} else {
			RenderTemplate(writer, "table_none.html", tableName)
		}
	} else if "PUT" == request.Method {
		// Update the table with the provided PUT body content.
		// The body content is expected to be in JSON format.
		// The new table entry overwrites any existing content.
		v, err := DecodeEntry(request.Body)
		if err != nil {
			http.Error(writer, "Server error", http.StatusInternalServerError)
		} else {
			Tables[tableName].Data[entryName] = v
			PersistTables()
			http.Redirect(writer, request, "/tables/"+tableName+"/"+entryName, http.StatusFound)
		}
	} else if "DELETE" == request.Method {
		// Remove the table entry from the table or render a bad request message
		// if the table does not exist.
		if table, exists := Tables[tableName]; exists {
			delete(table.Data, entryName)
			PersistTables()
			http.Redirect(writer, request, "/tables/"+tableName, http.StatusFound)
		} else {
			http.Error(writer, "Invalid request", http.StatusBadRequest)
		}
	} else {
		// Display an invalid request message
		http.Error(writer, "Invalid request", http.StatusBadRequest)
	}
}
