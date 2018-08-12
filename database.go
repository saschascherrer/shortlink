package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sync"
)

type ShortlinkDB interface {
	Load(filename string) error
	Save(filename string) error
	Add(key, value string) error
	Get(key string) (string, error)
}

type Database struct {
	filename string
	mux      sync.RWMutex
	entries  map[string]string
}

func NewDatabase(filename string) (*Database, error) {
	if filename == "" {
		return nil, errors.New("Empty Filename")
	}

	db := &Database{}
	db.filename = filename
	db.entries = make(map[string]string)
	return db, nil
}

func (db *Database) Load(filename string) error {
	// Fallback to stored filename, if empty filename is provided
	if filename == "" {
		filename = db.filename
	}

	// Abort if filename is still empty
	if filename == "" {
		return errors.New("Empty Filename")
	}

	db.mux.Lock()
	defer db.mux.Unlock()

	db.filename = filename
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	db.entries = make(map[string]string)
	err = json.Unmarshal(data, &db.entries)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) Save(filename string) error {
	// Fallback to stored filename, if empty filename is provided
	if filename == "" {
		filename = db.filename
	}

	// Abort if filename is still empty
	if filename == "" {
		return errors.New("Empty Filename")
	}

	db.mux.RLock()
	defer db.mux.RUnlock()

	data, err := json.Marshal(db.entries)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0755)
}

func (db *Database) Add(key, value string) error {
	if key == "" && value == "" {
		return errors.New("Provided Key and Value are empty")
	} else if key == "" {
		return errors.New("Provided Key is empty")
	} else if value == "" {
		return errors.New("Provided Value is empty")
	}

	db.mux.Lock()
	defer db.mux.Unlock()

	val := db.entries[key]
	if val == "" {
		db.entries[key] = value
	} else {
		return errors.New("Key is already in use")
	}

	return nil
}

func (db *Database) Get(key string) (string, error) {
	if key == "" {
		return "", errors.New("Provided Key is empty")
	}

	db.mux.RLock()
	val := db.entries[key]
	db.mux.RUnlock()

	if val == "" {
		return "", errors.New("No value found for provided key")
	}

	return val, nil
}
