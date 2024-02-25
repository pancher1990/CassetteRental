package response

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk    = "Ok"
	StatusError = "Error"
)

func Ok() Response {
	return Response{Status: StatusOk}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errors validator.ValidationErrors) Response {
	var errMsgs []string
	for _, err := range errors {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}
	return Error(strings.Join(errMsgs, ", "))
}
func BadRequest(w http.ResponseWriter, message string) {
	sendError(w, http.StatusBadRequest, message)
}

func InternalServerError(w http.ResponseWriter, message string) {
	sendError(w, http.StatusInternalServerError, message)
}

func StatusNotFound(w http.ResponseWriter, message string) {
	sendError(w, http.StatusNotFound, message)

}

func StatusConflict(w http.ResponseWriter, message string) {
	sendError(w, http.StatusConflict, message)

}
func sendError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Error(message))
}
