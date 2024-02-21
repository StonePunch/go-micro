package tools

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// defaultMaxUpload is the default max upload size (10 mb)
const defaultMaxUpload = 10485760

// Tools is the type for this package.
type Tools struct {
	MaxJSONSize int // MaxJSONSize is the maximum size of JSON file in bytes
}

// New returns a new toolbox with sensible defaults.
func New() Tools {
	return Tools{
		MaxJSONSize: defaultMaxUpload,
	}
}

// JSONResponse is the type used for sending JSON around.
type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// ReadJSON tries to convert the body of a request from JSON to a variable.
// The third parameter, data, is expected to be a pointer, so that we can read
// data into it.
func (t *Tools) ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	// Set a sensible default for the maximum payload size.
	maxBytes := defaultMaxUpload

	// If MaxJSONSize is set, use that value instead of default.
	if t.MaxJSONSize != 0 {
		maxBytes = t.MaxJSONSize
	}

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	// The decoded value here does not matter so &struct{}{} is used this logic
	// is just to check that the input only has one JSON value
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

// WriteJSON takes a response status code and arbitrary data and writes a JSON
// response to the client.
func (t *Tools) WriteJSON(w http.ResponseWriter, status int, data any, header ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(header) != 0 {
		for key, value := range header[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)

	return err
}

// ErrorJSON takes an error, and optionally a response status code, and sends
// a JSON error response.
func (t *Tools) ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) != 0 {
		statusCode = status[0]
	}

	var payload JsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return t.WriteJSON(w, statusCode, payload)
}
