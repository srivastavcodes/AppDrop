package main

import (
	"net/http"
)

func (b *backend) healthcheckHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("url working"))
}
