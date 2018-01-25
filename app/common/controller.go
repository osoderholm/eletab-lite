package common

import (
	"net/http"
	"encoding/json"
)

// Default controller, to be implemented in more specific controllers
type Controller struct {
}

// Internal structs for better json structure
type jsonError struct {
	Error string `json:"error"`
}
type jsonErrors struct {
	Errors []jsonError `json:"errors"`
}

// Sends json encoded data via httpResponseWriter with provided code
func (c *Controller) SendJSON(w http.ResponseWriter, code int, data interface{})  {
	marshaled, err := json.Marshal(data)

	if err != nil {
		c.SendErrorsJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(marshaled)
}

// Send json encoded error message
func (c *Controller) SendErrorsJSON(w http.ResponseWriter, code int, errors... string)  {
	errs := jsonErrors{Errors:make([]jsonError, 0)}

	for _, e := range errors {
		errs.Errors = append(errs.Errors, jsonError{Error:e})
	}

	marshaled, err := json.Marshal(errs)

	if err != nil {
		marshaled = []byte("{\"errors\":[{\"error\":\"Internal server error\"}]}")
		code = http.StatusInternalServerError
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(marshaled)
}