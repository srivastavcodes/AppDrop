package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (b *backend) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", b.healthcheckHandler)
	return router
}
