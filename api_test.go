package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseAPI(t *testing.T) {
	db, err := NewDatabase(testdbpath)
	assert.NoError(t, err)

	ts := httptest.NewServer(DatabaseAPI(db))
	defer ts.Close()

	// Get method
	res, err := http.Get(ts.URL + "/manage/")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Empty body
	res, err = http.Post(ts.URL+"/manage/", "text/json", bytes.NewBuffer([]byte("")))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, []byte(http.StatusText(http.StatusUnprocessableEntity)+"\n"), body)
	assert.Equal(t, 0, len(db.entries))

	// Body misses Target
	res, err = http.Post(ts.URL+"/manage/", "text/json", bytes.NewBuffer([]byte("{\"Target\":\"Value1\"}")))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
	body, err = ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, []byte(http.StatusText(http.StatusUnprocessableEntity)+"\n"), body)
	assert.Equal(t, 0, len(db.entries))

	// Body misses Target
	res, err = http.Post(ts.URL+"/manage/", "text/json", bytes.NewBuffer([]byte("{\"Key\":\"Key1\"}")))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
	body, err = ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, []byte(http.StatusText(http.StatusUnprocessableEntity)+"\n"), body)
	assert.Equal(t, 0, len(db.entries))

	// Body OK
	res, err = http.Post(ts.URL+"/manage/", "text/json", bytes.NewBuffer([]byte("{\"Key\":\"Key1\", \"Target\":\"Value1\"}")))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	body, err = ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, []byte(""), body)
	assert.Equal(t, 1, len(db.entries))
	storedValue, err := db.Get("Key1")
	assert.NoError(t, err)
	assert.Equal(t, "Value1", storedValue)

	// Duplicate
	res, err = http.Post(ts.URL+"/manage/", "text/json", bytes.NewBuffer([]byte("{\"Key\":\"Key1\", \"Target\":\"AnotherValue2\"}")))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusPreconditionFailed, res.StatusCode)
	body, err = ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, []byte("Key is already in use\n"), body)
	assert.Equal(t, 1, len(db.entries))
	storedValue, err = db.Get("Key1")
	assert.NoError(t, err)
	assert.Equal(t, "Value1", storedValue)
}
