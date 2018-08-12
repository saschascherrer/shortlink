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

func Redirector() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		responseString := "<html><body>Hello World</body></html>"
		w.Write([]byte(responseString))
	}
}

func main() {
	fmt.Println("+--------------------------------------------+")
	fmt.Println("| Shortlink v0.1.0 (alpha)         GNU GPLv3 |")
	fmt.Println("| by Sascha Scherrer <dev@saschascherrer.de> |")
	fmt.Println("+--------------------------------------------+")

	var socket string
	flag.StringVar(&socket, "socket", ":4242", "The socket to listen on")
	flag.Parse()

	log.Printf("Starting Shortlink Server on %s\n", socket)

	http.HandleFunc("/", Redirector())
	log.Fatalln(http.ListenAndServe(socket, nil))
}
