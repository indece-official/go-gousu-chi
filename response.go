package gousuchi

import (
	"encoding/json"
	"io"
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
	ContentTypeApplicationJSON        ContentType = "application/json"
	ContentTypeApplicationOctetStream ContentType = "application/octet-stream"
	ContentTypeApplicationPDF         ContentType = "application/pdf"
	ContentTypeTextPlain              ContentType = "text/plain"
	ContentTypeTextHTML               ContentType = "text/html"
	ContentTypeTextCSV                ContentType = "text/csv"
	ContentTypeImagePNG               ContentType = "image/png"
	ContentTypeImageJPEG              ContentType = "image/jpeg"
	ContentTypeImageBMP               ContentType = "image/bmp"
)

type Response struct {
	responseError   *ResponseError
	Request         *http.Request
	StatusCode      int
	Header          http.Header
	ContentType     ContentType
	Body            []byte
	BodyReader      io.Reader
	DetailedMessage string
	DisableLogging  bool
}

var _ IResponse = (*Response)(nil)

func (r *Response) GetRequest() *http.Request {
	return r.Request
}

func (r *Response) Write(w http.ResponseWriter) IResponse {
	if r.responseError != nil {
		return r.responseError
	}

	if r.Header != nil {
		for field, values := range r.Header {
			w.Header()[field] = values
		}
	}

	w.Header().Set("Content-Type", string(r.ContentType))
	w.WriteHeader(r.StatusCode)

	if r.BodyReader != nil {
		io.Copy(w, r.BodyReader)
	} else {
		w.Write(r.Body)
	}

	return nil
}

func (r *Response) Log(log *gousu.Log) {
	message := r.DetailedMessage
	if message == "" {
		message = "OK"
	}

	log.Infof("%s %s - %d %s", r.Request.Method, r.Request.RequestURI, r.StatusCode, message)
}

func NewResponse(
	request *http.Request,
	statusCode int,
	contentType ContentType,
	body []byte,
) *Response {
	return &Response{
		Request:     request,
		StatusCode:  statusCode,
		ContentType: contentType,
		Body:        body,
	}
}

func NewStreamResponse(
	request *http.Request,
	statusCode int,
	contentType ContentType,
	bodyReader io.Reader,
) *Response {
	return &Response{
		Request:     request,
		StatusCode:  statusCode,
		ContentType: contentType,
		BodyReader:  bodyReader,
	}
}

// JSON creates a new RestResponse of type application/json
func JSON(request *http.Request, obj interface{}) *Response {
	body, err := json.Marshal(obj)
	if err != nil {
		return &Response{
			Request:       request,
			responseError: InternalServerError(request, "Can't json encode response: %s", err),
		}
	}

	return &Response{
		Request:     request,
		StatusCode:  http.StatusOK,
		ContentType: ContentTypeApplicationJSON,
		Body:        body,
	}
}

// Text creates a new RestResponse of type text/plain
func Text(request *http.Request, body string) *Response {
	return &Response{
		Request:     request,
		StatusCode:  http.StatusOK,
		ContentType: ContentTypeTextPlain,
		Body:        []byte(body),
	}
}

// HTML creates a new RestResponse of type text/html
func HTML(request *http.Request, body string) *Response {
	return &Response{
		Request:     request,
		StatusCode:  http.StatusOK,
		ContentType: ContentTypeTextHTML,
		Body:        []byte(body),
	}
}
