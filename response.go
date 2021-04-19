package gousuchi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/indece-official/go-gousu"
)

type ResponseType string

const (
	ResponseTypeJSON ResponseType = "json"
	ResponseTypeText ResponseType = "text"
)

type Response struct {
	StatusCode      int
	Header          http.Header
	Type            ResponseType
	Body            interface{}
	DetailedMessage string
	DisableLogging  bool
}

func (r *Response) Write(w http.ResponseWriter) *ResponseError {
	var err error
	respData := []byte{}
	contentType := "text/plain"

	if r.Body != nil {
		switch r.Type {
		case ResponseTypeJSON:
			respData, err = json.Marshal(r.Body)
			if err != nil {
				return InternalServerError(fmt.Errorf("can't json encode response: %s", err))
			}
			contentType = "application/json"
		case ResponseTypeText:
			respDataStr, ok := r.Body.(string)
			if !ok {
				return InternalServerError(fmt.Errorf("response is not of type string"))
			}

			respData = []byte(respDataStr)

			contentType = "text/plain"
		default:
			return InternalServerError(fmt.Errorf("unsupported content type '%s'", r.Type))
		}
	}

	if r.Header != nil {
		for field, values := range r.Header {
			w.Header()[field] = values
		}
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(r.StatusCode)
	w.Write(respData)

	return nil
}

func (r *Response) Log(req *http.Request, log *gousu.Log) {
	message := r.DetailedMessage
	if message == "" {
		message = "OK"
	}

	log.Infof("%s %s - %d %s", req.Method, req.RequestURI, r.StatusCode, message)
}

// JSON creates a new RestResponse of type application/json
func JSON(obj interface{}) *Response {
	return &Response{
		StatusCode: http.StatusOK,
		Type:       ResponseTypeJSON,
		Body:       obj,
	}
}

// Text creates a new RestResponse of type text/plain
func Text(obj interface{}) *Response {
	return &Response{
		StatusCode: http.StatusOK,
		Type:       ResponseTypeText,
		Body:       obj,
	}
}
