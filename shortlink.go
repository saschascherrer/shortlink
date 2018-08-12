package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Resolver interface {
	Resolve(key string) string
}

type ResolverFunc func(key string) string

func (rf ResolverFunc) Resolve(key string) string {
	return rf(key)
}

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

func flags() (string, string) {
	var socket, dbfile string
	flag.StringVar(&socket, "socket", ":4242", "The socket to listen on")
	flag.StringVar(&dbfile, "dbfile", "./shortlink.db", "Database File")
	flag.Parse()
	return socket, dbfile
}

func start(socket string) {
	log.Printf("Starting Shortlink Server on %s\n", socket)
	log.Fatalln(http.ListenAndServe(socket, nil))
}

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

func main() {
	banner()
	socket, dbfile := flags()
	http.Handle("/", Server(dbfile))
	start(socket)
}
