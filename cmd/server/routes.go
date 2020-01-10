package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (api *httpAPI) routes() http.Handler {
	middleware := alice.New(api.recoverPanic, api.logRequest, secureHeaders)
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(api.home))
	mux.Get("/set/:key", http.HandlerFunc(api.handleSet))
	mux.Get("/add/:key", http.HandlerFunc(api.handleAdd))
	mux.Get("/replace/:key", http.HandlerFunc(api.handleReplace))
	mux.Get("/append/:key", http.HandlerFunc(api.handleAppend))
	mux.Get("/prepend/:key", http.HandlerFunc(api.handlePrepend))
	mux.Get("/increment/:key", http.HandlerFunc(api.handleIncrement))
	mux.Get("/decrement/:key", http.HandlerFunc(api.handleDecrement))
	mux.Get("/cas/:key", http.HandlerFunc(api.handleCompareAndSwap))
	mux.Get("/get/:key", http.HandlerFunc(api.handleGet))
	mux.Get("/gets/:key", http.HandlerFunc(api.handleGetEntryPlusToken))
	mux.Get("/delete/:key", http.HandlerFunc(api.handleDelete))
	mux.Get("/clear", http.HandlerFunc(api.handleClear))
	mux.Get("/stats", http.HandlerFunc(api.handleStats))
	return middleware.Then(mux)
}
