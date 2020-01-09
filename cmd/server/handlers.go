package main

import "net/http"

import "strconv"

import "fmt"

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	var err error = nil
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
	w.Write([]byte("Hello go-memcached"))
}

func (app *application) getValue(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":key"))
	if err != nil {
		app.notFound(w)
		return
	}
	w.Write([]byte(fmt.Sprintf("val for key %d", id)))
}
