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

func Redirector() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World\n")
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

	http.HandleFunc("/", Redirector())

	log.Printf("Starting Shortlink Server on %s\n", socket)
	log.Fatalln(http.ListenAndServe(socket, nil))
}
