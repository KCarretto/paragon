package http

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/kcarretto/paragon/pkg/script"
)

// NewRequest creates a new Request object to be passed around.
//
//go:generate go run ../gendoc.go -lib http -func newRequest -param url@String -retval request@Request -doc "NewRequest creates a new Request object to be passed around."
//
// @callable:	http.NewRequest
// @param:		url		@String
// @retval:		request @Request
//
// @usage:		r = http.NewRequest(CDN_URL+"/l/nomnom")
func NewRequest(url string) *Request {
	return &Request{
		Url:    url,
		Method: "GET",
	}
}

func newRequest(parser script.ArgParser) (script.Retval, error) {
	url, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}
	return NewRequest(url), nil
}

// SetMethod sets the http method on the request object.
//
//go:generate go run ../gendoc.go -lib http -func setMethod -param r@Request -param method@String -doc "SetMethod sets the http method on the request object."
//
// @callable:	http.SetMethod
// @param:		r		@Request
// @param:		method	@string
//
// @usage:		http.SetMethod(r, "POST")
func SetMethod(r *Request, method string) {
	r.Method = method
}

func setMethod(parser script.ArgParser) (script.Retval, error) {
	r, err := ParseParam(parser, 0)
	if err != nil {
		return nil, err
	}
	method, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}
	SetMethod(r, method)
	return nil, nil
}

// SetHeader sets the http header to the value passed on the request object.
//
//go:generate go run ../gendoc.go -lib http -func setHeader -param r@Request -param header@String -param value@String -doc "SetHeader sets the http header to the value passed on the request object."
//
// @callable:	http.SetHeader
// @param:		r		@Request
// @param:		header	@string
// @param:		value	@string
//
// @usage:		http.SetHeader(r, "Content-Type", "application/json")
func SetHeader(r *Request, header string, value string) {
	r.Headers[header] = value
}

func setHeader(parser script.ArgParser) (script.Retval, error) {
	r, err := ParseParam(parser, 0)
	if err != nil {
		return nil, err
	}
	header, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}
	value, err := parser.GetString(2)
	if err != nil {
		return nil, err
	}
	SetHeader(r, header, value)
	return nil, nil
}

// SetBody sets the http body to the value passed on the request object.
//
//go:generate go run ../gendoc.go -lib http -func setBody -param r@Request -param value@String -doc "SetBody sets the http body to the value passed on the request object."
//
// @callable:	http.SetBody
// @param:		r		@Request
// @param:		value	@string
//
// @usage:		http.SetBody(r, "{key: value}")
func SetBody(r *Request, value string) {
	r.Body = value
}

func setBody(parser script.ArgParser) (script.Retval, error) {
	r, err := ParseParam(parser, 0)
	if err != nil {
		return nil, err
	}
	value, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}
	SetBody(r, value)
	return nil, nil
}

// Exec sends the passed request object.
//
//go:generate go run ../gendoc.go -lib http -func exec -param r@Request -retval response@String -retval err@Error -doc "Exec sends the passed request object."
//
// @callable:	http.Exec
// @param:		r			@Request
// @retval:		response	@string
// @retval:		err 		@Error
//
// @usage:		resp, err = http.Exec(r)
func Exec(r *Request) (string, error) {
	client := &http.Client{}
	var req *http.Request
	var err error
	if r.Body != "" {
		req, err = http.NewRequest(r.Method, r.Url, bytes.NewBufferString(r.Body))
	} else {
		req, err = http.NewRequest(r.Method, r.Url, nil)
	}
	if err != nil {
		return "", err
	}
	for k, v := range r.Headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func exec(parser script.ArgParser) (script.Retval, error) {
	r, err := ParseParam(parser, 0)
	if err != nil {
		return nil, err
	}
	retVal, retErr := Exec(r)
	return script.WithError(retVal, retErr), nil
}
