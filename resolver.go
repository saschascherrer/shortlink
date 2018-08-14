package main

import (
	"fmt"
	"net/http"
	"strings"
)

// Resolver is the interface for all Resolver
type Resolver interface {
	Resolve(key string) string
}

// ResolverFunc makes it possible to pass a function as Resolver
type ResolverFunc func(key string) string

// Resolve makes the ResolverFunc compatible to the Resolver interface.
// It executes the ResolverFunc when Resolve is called on a ResolverFunc.
func (rf ResolverFunc) Resolve(key string) string {
	return rf(key)
}

// HTTPAction is the interface for an action called by the Redirector
// if the target is a non-empty string.
type HTTPAction interface {
	Execute(w http.ResponseWriter, r *http.Request, target string)
}

// HTTPRedirect is a SuccessAction that performs an HTTP redirect
type HTTPRedirect struct{}

// Execute performs a redirect to the target
func (HTTPRedirect) Execute(w http.ResponseWriter, r *http.Request, target string) {
	http.Redirect(w, r, target, http.StatusTemporaryRedirect)
}

// HTTPPrint is a SuccessAction that writes the target to the HTTP body
type HTTPPrint struct{}

// Execute prints the target to the HTTP body
func (HTTPPrint) Execute(w http.ResponseWriter, r *http.Request, target string) {
	fmt.Fprint(w, target)
}

// Redirector extracts the short-URL (key) from the URL Path and
// uses the provided Resolver to determine the corresponding long URL
// (target) and does an HTTP Redirect (Temporary Redirect) to that target.
// If the lookup through the resolver fails (i.e. result is the empty string),
// it returns an HTTP 404 Not found.
func Redirector(resolver Resolver, action HTTPAction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var target string

		// Extract Key
		key := strings.TrimPrefix(r.URL.EscapedPath(), "/r/")
		if key == "/" {
			target = ""
		} else {
			target = resolver.Resolve(key)
		}

		// Redirect if found, error otherwise
		if target != "" {
			action.Execute(w, r, target)
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, http.StatusText(http.StatusNotFound))
		}
	}
}

// DatabaseResolver takes a database and returns a Resolver that
// looks up the key in the provided database returning the value
// returned from db.Get(key) or an empty string if an error occurred.
func DatabaseResolver(db ShortlinkDB) Resolver {
	var resolver ResolverFunc
	resolver = func(key string) string {
		val, err := db.Get(key)
		if err != nil {
			return ""
		}
		return val
	}
	return resolver
}
