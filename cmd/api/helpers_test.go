package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func TestBackend_readIdParam(t *testing.T) {
	b := &backend{}

	t.Run("valid id", func(t *testing.T) {
		want := uuid.New()

		req := httptestRequestWithParams(t, httprouter.Params{
			{Key: "id", Value: want.String()}, // arbitrary bytes are allowed in a Go string
		})
		got, err := b.readIdParam(req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got != want {
			t.Fatalf("expected %v, got %v", want, got)
		}
	})
	t.Run("invalid id (wrong length)", func(t *testing.T) {
		req := httptestRequestWithParams(t, httprouter.Params{
			{Key: "id", Value: "not-16-bytes"},
		})
		got, err := b.readIdParam(req)
		if err == nil {
			t.Fatalf("expected error, got nil (uuid=%v)", got)
		}
		if got != uuid.Nil {
			t.Fatalf("expected uuid.Nil, got %v", got)
		}
	})
	t.Run("invalid id (nil uuid)", func(t *testing.T) {
		req := httptestRequestWithParams(t, httprouter.Params{
			{Key: "id", Value: string(make([]byte, 16))}, // 16 zero bytes => uuid.Nil
		})
		got, err := b.readIdParam(req)
		if err == nil {
			t.Fatalf("expected error, got nil (uuid=%v)", got)
		}
		if got != uuid.Nil {
			t.Fatalf("expected uuid.Nil, got %v", got)
		}
	})
}

func httptestRequestWithParams(t *testing.T, ps httprouter.Params) *http.Request {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)

	ctx := context.WithValue(req.Context(), httprouter.ParamsKey, ps)
	return req.WithContext(ctx)
}
