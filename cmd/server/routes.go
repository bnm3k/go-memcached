package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	middleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/get/:key", http.HandlerFunc(app.getValue))
	return middleware.Then(mux)
}
