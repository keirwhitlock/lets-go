package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	static_dir := flag.String("static", "./ui/static/", "Directory where static files are located.")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fileserver := http.FileServer(http.Dir(*static_dir))

	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
