package gousuchi

import (
	"fmt"
	"net/http"

	"github.com/indece-official/go-gousu"
)

type ResponseError struct {
	StatusCode    int
	PublicMessage string
	DetailedError error
}

func (r *ResponseError) Write(w http.ResponseWriter) {
	w.WriteHeader(r.StatusCode)
	fmt.Fprintf(w, r.PublicMessage)
}

func (r *ResponseError) Log(req *http.Request, log *gousu.Log) {
	if r.StatusCode >= 500 {
		log.Errorf("%s %s - %d %s", req.Method, req.RequestURI, r.StatusCode, r.DetailedError)
	} else {
		log.Warnf("%s %s - %d %s", req.Method, req.RequestURI, r.StatusCode, r.DetailedError)
	}
}

func InternalServerError(detailedError error) *ResponseError {
	return &ResponseError{
		StatusCode:    http.StatusInternalServerError,
		PublicMessage: "Internal server error",
		DetailedError: detailedError,
	}
}

func NotFound(detailedError error) *ResponseError {
	return &ResponseError{
		StatusCode:    http.StatusNotFound,
		PublicMessage: "Not found",
		DetailedError: detailedError,
	}
}

func BadRequest(detailedError error) *ResponseError {
	return &ResponseError{
		StatusCode:    http.StatusBadRequest,
		PublicMessage: "Bad request",
		DetailedError: detailedError,
	}
}
