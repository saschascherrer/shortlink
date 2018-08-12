package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
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
		target := resolver.Resolve("key")
		if target != "" {
			http.Redirect(w, r, target, http.StatusTemporaryRedirect)
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, http.StatusText(http.StatusNotFound))
		}
	}
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

	var resolver ResolverFunc
	resolver = func(key string) string {
		return "resolved"
	}

	http.HandleFunc("/", Redirector(resolver))

	log.Printf("Starting Shortlink Server on %s\n", socket)
	log.Fatalln(http.ListenAndServe(socket, nil))
}
