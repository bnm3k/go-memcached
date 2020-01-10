package main

import (
	"log"
	"net/http"
	"os"

	"github.com/nagamocha3000/go-memcached/pkg/cache"
)

type httpAPI struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	cache    *cache.Adapter
}

func main() {

	cfg := getConfig()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	api := &httpAPI{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  api.routes(),
	}
	infoLog.Printf("starting server on %s", cfg.addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
