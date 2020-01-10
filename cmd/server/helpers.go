package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (api *httpAPI) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	api.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (api *httpAPI) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (api *httpAPI) notFound(w http.ResponseWriter) {
	api.clientError(w, http.StatusNotFound)
}
