package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBanner(t *testing.T) {
	banner()
}

func TestFlags(t *testing.T) {
	flags()
}

func TestServer(t *testing.T) {
	err := ioutil.WriteFile(testdbpath, []byte(filecontent), 0754)
	assert.NoError(t, err)
	Server(testdbpath)
}
