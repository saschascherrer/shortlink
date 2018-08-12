package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

type Resolver interface {
	Resolve(key string) (string, error)
}

type ResolverFunc func(key string) (string, error)

func (rf ResolverFunc) Resolve(key string) (string, error) {
	return rf(key)
}

func Redirector(resolver Resolver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path
		target, err := resolver.Resolve(key)
		if err == nil {
			http.Redirect(w, r, target, http.StatusTemporaryRedirect)
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Entry for '%s' not found. Error Message: %s", key, err.Error())
	}
}

type LocalResolver struct{}

func (r LocalResolver) Resolve(key string) (string, error) {
	return "localhost:4242/#resolved=" + key, nil
}

func banner() {
	fmt.Println("+--------------------------------------------+")
	fmt.Println("| Shortlink v0.1.0 (alpha)         GNU GPLv3 |")
	fmt.Println("| by Sascha Scherrer <dev@saschascherrer.de> |")
	fmt.Println("+--------------------------------------------+")
}

func main() {
	banner()

	var socket string
	flag.StringVar(&socket, "socket", ":4242", "The socket to listen on")
	flag.Parse()

	http.HandleFunc("/", Redirector(LocalResolver{}))

	log.Printf("Starting Shortlink Server on %s\n", socket)
	log.Fatalln(http.ListenAndServe(socket, nil))
}
