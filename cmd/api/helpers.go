package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// envelope is a helper type for wrapping responses
type envelope map[string]any

// readIdParam reads and parses a UUID from URL parameters
func (b *backend) readIdParam(r *http.Request) (uuid.UUID, error) {
	params := httprouter.ParamsFromContext(r.Context())
	uid, err := uuid.Parse(params.ByName("id"))
	if err != nil || uid == uuid.Nil {
		return uuid.Nil, errors.New("invalid id parameter")
	}
	return uid, nil
}

// writeJson writes JSON response with headers
func (b *backend) writeJson(w http.ResponseWriter, status int, data any, headers http.Header) error {
	jsonB, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	jsonB = append(jsonB, '\n')
	for key, val := range headers {
		w.Header()[key] = val
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(jsonB); err != nil {
		b.logger.Info("Failed writing response", "err", err)
	}
	return nil
}

// readJson decodes the request body and populates the given dst field.
func (b *backend) readJson(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return b.decodeJsonError(err)
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("body must contain exactly one Json object")
	}
	return nil
}

func (b *backend) decodeJsonError(err error) error {
	var unmarshalTypeError *json.UnmarshalTypeError
	var syntaxError *json.SyntaxError
	var invalidUnmarshalError *json.InvalidUnmarshalError
	var maxBytesError *http.MaxBytesError

	switch {
	case errors.As(err, &syntaxError):
		return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
	case errors.Is(err, io.ErrUnexpectedEOF):
		return fmt.Errorf("body contains badly-formed JSON")
	case errors.Is(err, io.EOF):
		return fmt.Errorf("body must not be empty")
	case errors.As(err, &unmarshalTypeError):
		if unmarshalTypeError.Field != "" {
			return fmt.Errorf(
				"body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
		}
		return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		return fmt.Errorf("body contains unknown JSON field for key %s", fieldName)
	case errors.As(err, &maxBytesError):
		return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
	case errors.As(err, &invalidUnmarshalError):
		panic(err)
	default:
		return err
	}
}
