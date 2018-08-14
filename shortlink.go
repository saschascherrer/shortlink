package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

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

	router.HandleFunc("/r/", Redirector(DatabaseResolver(db), HTTPRedirect{})) // redirect from short-URL
	router.HandleFunc("/s/", Redirector(DatabaseResolver(db), HTTPPrint{}))    // show target of short-URL
	router.HandleFunc("/list/", http.NotFound)                                 // show all short-URLs and their target
	router.Handle("/manage/", DatabaseAPI(db))                                 // manage short-URLs and their target (API)
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
