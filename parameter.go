package gousuchi

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gopkg.in/guregu/null.v4"
)

// QueryParamInt64 loads a parameter from the url's query and parses it as int64
//
// If the parameter is not a valid int64 a BadRequest-Response is returned, else
// the response is nil
func QueryParamInt64(request *http.Request, name string) (int64, IResponse) {
	valueStr := request.URL.Query().Get(name)
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0, BadRequest(request, "Invalid query param %s (value: '%s'): %s", name, valueStr, err)
	}

	return value, nil
}

// OptionalQueryParamInt64 loads a parameter from the url's query and parses it as int64
//
// If the parameter is not empty and not a valid int64 a BadRequest-Response is returned, else
// the response is nil
func OptionalQueryParamInt64(request *http.Request, name string) (null.Int, IResponse) {
	valueStr := request.URL.Query().Get(name)
	if valueStr == "" {
		return null.Int{}, nil
	}

	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return null.Int{}, BadRequest(request, "Invalid query param %s (value: '%s'): %s", name, valueStr, err)
	}

	return null.IntFrom(value), nil
}

// QueryParamBool loads a parameter from the url's query and parses it as bool
//
// If the parameter is not a valid bool a BadRequest-Response is returned, else
// the response is nil. Accepted values are 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
func QueryParamBool(request *http.Request, name string) (bool, IResponse) {
	valueStr := request.URL.Query().Get(name)
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return false, BadRequest(request, "Invalid query param %s (value: '%s'): %s", name, valueStr, err)
	}

	return value, nil
}

// OptionalQueryParamBool loads a parameter from the url's query and parses it as bool
//
// If the parameter is not empty and not a valid bool a BadRequest-Response is returned, else
// the response is nil. Accepted values are 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
func OptionalQueryParamBool(request *http.Request, name string) (null.Bool, IResponse) {
	valueStr := request.URL.Query().Get(name)
	if valueStr == "" {
		return null.Bool{}, nil
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return null.Bool{}, BadRequest(request, "Invalid query param %s (value: '%s'): %s", name, valueStr, err)
	}

	return null.BoolFrom(value), nil
}

// URLParamInt64 loads a parameter from the url and parses it as int64
//
// If the parameter is not a valid int64 a BadRequest-Response is returned, else
// the response is nil
func URLParamInt64(request *http.Request, name string) (int64, IResponse) {
	valueStr := chi.URLParam(request, name)
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0, BadRequest(request, "Invalid url param %s (value: '%s'): %s", name, valueStr, err)
	}

	return value, nil
}

// URLParamBool loads a parameter from the url and parses it as bool
//
// If the parameter is not a valid bool a BadRequest-Response is returned, else
// the response is nil. Accepted values are 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
func URLParamBool(request *http.Request, name string) (bool, IResponse) {
	valueStr := chi.URLParam(request, name)
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return false, BadRequest(request, "Invalid url param %s (value: '%s'): %s", name, valueStr, err)
	}

	return value, nil
}
