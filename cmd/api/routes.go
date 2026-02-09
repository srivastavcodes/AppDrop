package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (b *backend) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", b.healthcheckHandler)
	return b.recoverPanic(b.enableCors(b.logRequest(router)))
}

func (b *backend) healthcheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("healthy and running"))
}
