package main

import (
	"flag"
	"fmt"
	"log"
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

// Redirector extracts the short-URL (key) from the URL Path and
// uses the provided Resolver to determine the corresponding long URL
// (target) and does an HTTP Redirect (Temporary Redirect) to that target.
// If the lookup through the resolver fails (i.e. result is the empty string),
// it returns an HTTP 404 Not found.
func Redirector(resolver Resolver) http.HandlerFunc {
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
			http.Redirect(w, r, target, http.StatusTemporaryRedirect)
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

func banner() {
	fmt.Println("" +
		"+---------------------------------------------------------+\n" +
		"| Shortlink v0.1.1 (alpha)                      GNU GPLv3 |\n" +
		"| by Sascha Scherrer <dev@saschascherrer.de>              |\n" +
		"+---------------------------------------------------------+")
}

// flags defines and parses the flags, returning their values
func flags() (string, string) {
	var socket, dbfile string
	flag.StringVar(&socket, "socket", ":4242", "The socket to listen on")
	flag.StringVar(&dbfile, "dbfile", "./shortlink.db", "Database File")
	flag.Parse()
	return socket, dbfile
}

// start the server on the provided socket. Print out where it listens
func start(socket string) {
	log.Printf("Starting Shortlink Server on %s\n", socket)
	log.Fatalln(http.ListenAndServe(socket, nil))
}

// Server composes the parts needed for the Shortlink server.
// Reads the database file and initializes Database.
// Configures Server routes
func Server(dbfile string) http.Handler {
	db, err := NewDatabase(dbfile)
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Load(dbfile)
	if err != nil {
		log.Fatalln(err)
	}

	router := http.NewServeMux()

	router.HandleFunc("/r/", Redirector(DatabaseResolver(db))) // redirect from short-URL
	router.HandleFunc("/s/", http.NotFound)                    // show target of short-URL
	router.HandleFunc("/list/", http.NotFound)                 // show all short-URLs and their target
	router.HandleFunc("/manage/", http.NotFound)               // show all short-URLs and their target
	router.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static"))))

	return router
}

// main entry point to the Shortlink application
func main() {
	banner()                         // print banner
	socket, dbfile := flags()        // get config
	http.Handle("/", Server(dbfile)) // configure server
	start(socket)                    // start server
}
