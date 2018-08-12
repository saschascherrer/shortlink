package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBanner(t *testing.T) {
	banner()
}

type TestResolver struct {
	value string
	err   error
}

func (r TestResolver) Resolve(key string) (string, error) {
	return r.value, r.err
}

// createServer is a convenience method to start a testing service
// with a provided Authenticator that just prints "Hello World" if
// the Authenticator returns boolean true.
func createServer(resolver Resolver) *httptest.Server {
	return httptest.NewServer(Redirector(resolver))
}

func TestRedirector(t *testing.T) {
	validResolver := TestResolver{"resolver.example/path?key=value#id", nil}
	ts := httptest.NewServer(Redirector(validResolver))
	defer ts.Close()

	http.Get(ts.URL)

}
