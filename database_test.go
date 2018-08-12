package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testdbpath string = "./test.db"
const filecontent string = "{\"key1\":\"value1\",\"key2\":\"value2\"}"

func TestNewOK(t *testing.T) {
	db, err := NewDatabase(testdbpath)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.Equal(t, testdbpath, db.filename)
	assert.Equal(t, 0, len(db.entries))
}

func TestNewEmptyPath(t *testing.T) {
	db, err := NewDatabase("")
	assert.EqualError(t, err, "Empty Filename")
	assert.Nil(t, db)
}

func TestAddErrors(t *testing.T) {
	db, err := NewDatabase(testdbpath)
	assert.NoError(t, err)

	err = db.Add("", "")
	assert.EqualError(t, err, "Provided Key and Value are empty")

	err = db.Add("", "value")
	assert.EqualError(t, err, "Provided Key is empty")

	err = db.Add("key", "")
	assert.EqualError(t, err, "Provided Value is empty")

	err = db.Add("key", "value")
	assert.NoError(t, err)
}

func TestAddDuplicate(t *testing.T) {
	db, err := NewDatabase(testdbpath)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(db.entries))

	err = db.Add("key", "first_value")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(db.entries))
	assert.Equal(t, "first_value", db.entries["key"])

	err = db.Add("key", "second_value")
	assert.EqualError(t, err, "Key is already in use")
	assert.Equal(t, 1, len(db.entries))
	assert.Equal(t, "first_value", db.entries["key"])
}

func TestAddOk(t *testing.T) {
	db, err := NewDatabase(testdbpath)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(db.entries))

	err = db.Add("key", "value")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(db.entries))
	assert.Equal(t, "value", db.entries["key"])
}

func TestGetErrors(t *testing.T) {
	db, err := NewDatabase(testdbpath)
	assert.NoError(t, err)
	db.Add("key", "value")

	val, err := db.Get("")
	assert.Error(t, err)
	assert.Equal(t, "", val)

	val, err = db.Get("fail")
	assert.EqualError(t, err, "No value found for provided key")
	assert.Equal(t, "", val)
}

func TestGetOk(t *testing.T) {
	db, err := NewDatabase(testdbpath)
	assert.NoError(t, err)
	db.Add("key", "value")

	val, err := db.Get("key")
	assert.NoError(t, err)
	assert.Equal(t, "value", val)
}

func TestSaveEmptyFilename(t *testing.T) {
	db := Database{filename: ""}
	assert.Equal(t, "", db.filename)

	err := db.Save("")
	assert.EqualError(t, err, "Empty Filename")

	db = Database{filename: testdbpath}
	assert.Equal(t, testdbpath, db.filename)

	err = db.Save("")
	assert.NoError(t, err)

	if _, err := os.Stat(testdbpath); err == nil {
		err := os.Remove(testdbpath)
		assert.NoError(t, err)
	} else {
		t.Error("File not written")
	}
}

func TestSave(t *testing.T) {
	db, err := NewDatabase(testdbpath)
	assert.NoError(t, err)
	db.Add("key1", "value1")
	db.Add("key2", "value2")
	db.Save("")

	if _, err := os.Stat(testdbpath); err == nil {

		actual, err := ioutil.ReadFile(testdbpath)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, []byte(filecontent), actual)

		err = os.Remove(testdbpath)
		assert.NoError(t, err)

	} else {
		t.Error("File not written")
	}
}

func TestLoad(t *testing.T) {
	err := ioutil.WriteFile(testdbpath, []byte(filecontent), 0754)
	assert.NoError(t, err)

	db, err := NewDatabase(testdbpath)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	err = db.Load(testdbpath)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(db.entries))

	if _, err := os.Stat(testdbpath); err == nil {
		err := os.Remove(testdbpath)
		assert.NoError(t, err)
	} else {
		t.Error("File not written")
	}

}
