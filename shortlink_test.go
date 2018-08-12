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

func TestRedirectorKeyNotFound(t *testing.T) {
	var resolver ResolverFunc
	resolver = func(key string) string {
		return ""
	}

	ts := httptest.NewServer(Redirector(resolver))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode, "wrong status")
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, []byte(http.StatusText(http.StatusNotFound)+"\n"), body)
}

func TestRedirectorExistingTarget(t *testing.T) {
	var resolver ResolverFunc
	resolver = func(key string) string {
		return "http://saschascherrer.de"
	}

	ts := httptest.NewServer(Redirector(resolver))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode, "wrong status")
}

func TestRedirectorNonexistingTarget(t *testing.T) {
	var resolver ResolverFunc
	resolver = func(key string) string {
		return "http://example.invalid"
	}

	ts := httptest.NewServer(Redirector(resolver))
	defer ts.Close()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := client.Get(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode, "wrong status")
	location, err := res.Location()
	assert.NoError(t, err)
	assert.Equal(t, "http://example.invalid", location.String())
}
