package httputils

import (
	"encoding/json"
	"io"
	"net/http"
)

// HTTPError - contains data that returned in http response body when encounter error
type HTTPError struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

// SendTextResponse - send http response with text body and given status code
func SendTextResponse(w http.ResponseWriter, code int, text string) error {
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(code)

	if _, err := io.WriteString(w, text); err != nil {
		return err
	}

	return nil
}

// SendJSONResponse - send http response with json body and given status code
func SendJSONResponse(w http.ResponseWriter, code int, data interface{}) error {
	w.Header().Set("content-type", "application/json")

	jsonData, err := json.Marshal(data)

	if err != nil {
		return err
	}

	w.WriteHeader(code)

	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}

// SendJSONErrorResponse - send http response with body with structure HTTPError and given status code
func SendJSONErrorResponse(w http.ResponseWriter, statusCode int, error string, errorCode int) error {
	return SendJSONResponse(w, statusCode, HTTPError{Error: error, Code: errorCode})
}

// SendRedirectResponse - send http response with redirect status code and Location header
func SendRedirectResponse(w http.ResponseWriter, location string) {
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// SendStatusCode - send http response with given status code
func SendStatusCode(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}
