package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sync"
)

// ShortlinkDB is an interface for the four methods a Database as is
// implemented in this file MUST support
type ShortlinkDB interface {
	Load(filename string) error
	Save(filename string) error
	Add(key, value string) error
	Get(key string) (string, error)
}

// The Database struct bundles the variables we need to keep track of
// for the database
type Database struct {
	filename string
	rwmutex  sync.RWMutex
	entries  map[string]string
}

// NewDatabase creates and returns a new instance of a ShortlinkDB.
// It initializes the fileds specified in the Database type struct.
func NewDatabase(filename string) (*Database, error) {
	if filename == "" {
		return nil, errors.New("Empty Filename")
	}

	db := &Database{}
	db.filename = filename
	db.entries = make(map[string]string)
	return db, nil
}

// Load the contents of the database from a file. The path to this file
// is provided by the filename parameter. If this parameter is empty, we
// attempt to use the filename specified when NewDatabase() was called.
func (db *Database) Load(filename string) error {
	// Fallback to stored filename, if empty filename is provided
	if filename == "" {
		filename = db.filename
	}

	// Abort if filename is still empty
	if filename == "" {
		return errors.New("Empty Filename")
	}

	db.rwmutex.Lock()
	defer db.rwmutex.Unlock()

	db.filename = filename
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	db.entries = make(map[string]string)
	return json.Unmarshal(data, &db.entries)
}

// Save the database to a JSON-File. The path to this file
// is provided by the filename parameter. If this parameter is empty, we
// attempt to use the filename specified when NewDatabase() was called.
func (db *Database) Save(filename string) error {
	// Fallback to stored filename, if empty filename is provided
	if filename == "" {
		filename = db.filename
	}

	// Abort if filename is still empty
	if filename == "" {
		return errors.New("Empty Filename")
	}

	db.rwmutex.RLock()
	defer db.rwmutex.RUnlock()

	data, err := json.Marshal(db.entries)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0755)
}

// Add a key-value-pair to the database.
// This does NOT automatically save to file.
func (db *Database) Add(key, value string) error {
	if key == "" && value == "" {
		return errors.New("Provided Key and Value are empty")
	} else if key == "" {
		return errors.New("Provided Key is empty")
	} else if value == "" {
		return errors.New("Provided Value is empty")
	}

	db.rwmutex.Lock()
	defer db.rwmutex.Unlock()

	val := db.entries[key]
	if val == "" {
		db.entries[key] = value
	} else {
		return errors.New("Key is already in use")
	}

	return nil
}

// Get the value that belongs to a key.
func (db *Database) Get(key string) (string, error) {
	if key == "" {
		return "", errors.New("Provided Key is empty")
	}

	db.rwmutex.RLock()
	val := db.entries[key]
	db.rwmutex.RUnlock()

	if val == "" {
		return "", errors.New("No value found for provided key")
	}

	return val, nil
}
