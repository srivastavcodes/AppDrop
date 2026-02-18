package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (b *backend) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(b.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(b.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/healthcheck", b.healthcheckHandler)

	// Store routes
	router.HandlerFunc(http.MethodGet, "/stores", b.listStoresHandler)
	router.HandlerFunc(http.MethodPost, "/stores", b.createStoreHandler)
	router.HandlerFunc(http.MethodGet, "/stores/:store_id", b.showStoreHandler)
	router.HandlerFunc(http.MethodPut, "/stores/:store_id", b.updateStoreHandler)
	router.HandlerFunc(http.MethodDelete, "/stores/:store_id", b.deleteStoreHandler)

	// Page routes — nested under store
	router.HandlerFunc(http.MethodGet, "/stores/:store_id/pages", b.listPagesHandler)
	router.HandlerFunc(http.MethodPost, "/stores/:store_id/pages", b.createPageHandler)
	router.HandlerFunc(http.MethodGet, "/stores/:store_id/pages/:page_id", b.showPageHandler)
	router.HandlerFunc(http.MethodPut, "/stores/:store_id/pages/:page_id", b.updatePageHandler)
	router.HandlerFunc(http.MethodDelete, "/stores/:store_id/pages/:page_id", b.deletePageHandler)

	// Widget routes — nested under store, page_id only where semantically required
	router.HandlerFunc(http.MethodPost, "/stores/:store_id/pages/:page_id/widgets", b.createWidgetHandler)
	router.HandlerFunc(http.MethodPost, "/stores/:store_id/pages/:page_id/widgets/reorder", b.reorderWidgetsHandler)
	router.HandlerFunc(http.MethodPut, "/stores/:store_id/widgets/:id", b.updateWidgetHandler)
	router.HandlerFunc(http.MethodDelete, "/stores/:store_id/widgets/:id", b.deleteWidgetHandler)

	return b.recoverPanic(b.enableCors(b.logRequest(router)))
}

func (b *backend) healthcheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("healthy and running"))
}
