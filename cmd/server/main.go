package main

import "net/http"

import "log"

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello go-memcached"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	log.Println("starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
