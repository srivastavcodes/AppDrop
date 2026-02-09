package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (b *backend) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/healthcheck", b.healthcheckHandler)

	router.NotFound = http.HandlerFunc(b.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(b.methodNotAllowedResponse)

	// Page routes
	router.HandlerFunc(http.MethodPost, "/pages", b.createPageHandler)
	router.HandlerFunc(http.MethodGet, "/pages/:id", b.showPageHandler)
	router.HandlerFunc(http.MethodGet, "/pages", b.listPagesHandler)
	router.HandlerFunc(http.MethodPut, "/pages/:id", b.updatePageHandler)
	router.HandlerFunc(http.MethodDelete, "/pages/:id", b.deletePageHandler)

	// Widget routes
	router.HandlerFunc(http.MethodPost, "/pages/:id/widgets", b.createWidgetHandler)
	router.HandlerFunc(http.MethodPost, "/pages/:id/widgets/reorder", b.reorderWidgetsHandler)
	router.HandlerFunc(http.MethodPut, "/widgets/:id", b.updateWidgetHandler)
	router.HandlerFunc(http.MethodDelete, "/widgets/:id", b.deleteWidgetHandler)

	return b.recoverPanic(b.enableCors(b.logRequest(router)))
}

func (b *backend) healthcheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("healthy and running"))
}
