package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBanner(t *testing.T) {
	banner()
}

func TestResolverFunc(t *testing.T) {
	var rf ResolverFunc
	rf = func(key string) string {
		return "value"
	}
	val := rf.Resolve("key")
	assert.Equal(t, "value", val)
}

func TestRedirector(t *testing.T) {
	ts := httptest.NewServer(Redirector())
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode, "wrong status")
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, []byte("Hello, World\n"), body)
}
