package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// The name of the persistent database file
const FILENAME = ".data.base.json"

// A simple table structure containing the
// table name and a map data structure
type Table struct {
	Name string
	Data map[string]interface{}
}

// Data structure containing the information for
// any table entry
type EntryInfo struct {
	TableName string
	EntryName string
	Value     interface{}
}

// Search data structure containing the search
// query and results information
type Search struct {
	Query   string
	Results []EntryInfo
}

// The global tables data
var Tables = map[string]Table{}

// Decode the data from the provided Reader instance to a
// table entry and return any error that occured during
// the decode process.  The data is expected to be in
// JSON format.
func DecodeEntry(reader io.Reader) (interface{}, error) {
	var obj interface{}

	dec := json.NewDecoder(reader)
	if err := dec.Decode(&obj); err != nil {
		return nil, err
	}

	return obj, nil
}

// Decode the data from the provided Reader instance to
// a Table instance and return any error that occured
// during the decode process. The data is expected to
// be in JSON object format.
func DecodeTable(reader io.Reader) (Table, error) {
	var t Table

	dec := json.NewDecoder(reader)
	if err := dec.Decode(&t.Data); err != nil {
		return Table{}, err
	}

	return t, nil
}

// Initialize the global tables with the content from
// the persistant database file storage.
func InitializeTables() error {
	// Check if the file exists, if not create it with an empty JSON object
	if _, err := os.Stat(FILENAME); err != nil {
		log.Println("Database persistance file does not exist")
		file, err := os.Create(FILENAME)

		if err != nil {
			return err
		}
		file.WriteString("{}")
		file.Close()
	}

	// Open the file for reading
	file, err := os.Open(FILENAME)
	if err != nil {
		return err
	}

	defer file.Close()

	dec := json.NewDecoder(file)

	// Decode the JSON object contained in the file
	if err = dec.Decode(&Tables); err != nil {
		return err
	}

	return nil
}

// Persist the global tables back into the file
// that will be read on application load.
func PersistTables() error {
	// Create the file, truncating any existing content
	writer, err := os.Create(FILENAME)
	if err != nil {
		return err
	}

	defer writer.Close()

	enc := json.NewEncoder(writer)

	// Encode the JSON object contained into the file
	if err = enc.Encode(&Tables); err != nil {
		return err
	}

	return nil
}
