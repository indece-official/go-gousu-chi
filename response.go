package gousuchi

import (
	"encoding/json"
	"net/http"

	"github.com/indece-official/go-gousu"
)

type IResponse interface {
	GetRequest() *http.Request
	Write(w http.ResponseWriter) IResponse
	Log(log *gousu.Log)
}

type ContentType string

const (
	ContentTypeApplicationJSON ContentType = "application/json"
	ContentTypeTextPlain       ContentType = "text/plain"
)

type Response struct {
	Request         *http.Request
	StatusCode      int
	Header          http.Header
	ContentType     ContentType
	Body            interface{}
	DetailedMessage string
	DisableLogging  bool
}

var _ IResponse = (*Response)(nil)

func (r *Response) GetRequest() *http.Request {
	return r.Request
}

func (r *Response) Write(w http.ResponseWriter) IResponse {
	var err error
	respData := []byte{}

	if r.Body != nil {
		switch r.ContentType {
		case ContentTypeApplicationJSON:
			respData, err = json.Marshal(r.Body)
			if err != nil {
				return InternalServerError(r.Request, "Can't json encode response: %s", err)
			}
		case ContentTypeTextPlain:
			respDataStr, ok := r.Body.(string)
			if !ok {
				return InternalServerError(r.Request, "Response is not of type string")
			}

			respData = []byte(respDataStr)
		default:
			var ok bool

			respData, ok = r.Body.([]byte)
			if !ok {
				return InternalServerError(r.Request, "Response is not of type bytes")
			}
		}
	}

	if r.Header != nil {
		for field, values := range r.Header {
			w.Header()[field] = values
		}
	}

	w.Header().Set("Content-Type", string(r.ContentType))
	w.WriteHeader(r.StatusCode)
	w.Write(respData)

	return nil
}

func (r *Response) Log(log *gousu.Log) {
	message := r.DetailedMessage
	if message == "" {
		message = "OK"
	}

	log.Infof("%s %s - %d %s", r.Request.Method, r.Request.RequestURI, r.StatusCode, message)
}

// JSON creates a new RestResponse of type application/json
func JSON(request *http.Request, obj interface{}) *Response {
	return &Response{
		Request:     request,
		StatusCode:  http.StatusOK,
		ContentType: ContentTypeApplicationJSON,
		Body:        obj,
	}
}

// Text creates a new RestResponse of type text/plain
func Text(request *http.Request, obj interface{}) *Response {
	return &Response{
		Request:     request,
		StatusCode:  http.StatusOK,
		ContentType: ContentTypeTextPlain,
		Body:        obj,
	}
}
