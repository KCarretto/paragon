package http

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/kcarretto/paragon/pkg/script"
)

// NewRequest creates a new Request object to be passed around.
//
// @callable:	http.NewRequest
// @param:		url	@String
// @retval:		err @Error
//
// @usage:		r = http.NewRequest(CDN_URL+"/l/nomnom")
func NewRequest(url string) (Request, error) {
	return Request{
		Url:    url,
		Method: "GET",
	}, nil
}

func newRequest(parser script.ArgParser) (script.Retval, error) {
	url, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}
	return NewRequest(url)
}

// SetMethod sets the http method on the request object.
//
// @callable:	http.SetMethod
// @param:		r		@Request
// @param:		method	@String
// @retval:		err 	@Error
//
// @usage:		http.SetMethod(r, "POST")
func SetMethod(r *Request, method string) error {
	r.Method = method
	return nil
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
	return nil, SetMethod(r, method)
}

// SetHeader sets the http header to the value passed on the request object.
//
// @callable:	http.SetHeader
// @param:		r		@Request
// @param:		header	@String
// @param:		value	@String
// @retval:		err 	@Error
//
// @usage:		http.SetHeader(r, "Content-Type", "application/json")
func SetHeader(r *Request, header string, value string) error {
	r.Headers[header] = value
	return nil
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
	return nil, SetHeader(r, header, value)
}

// SetBody sets the http body to the value passed on the request object.
//
// @callable:	http.SetBody
// @param:		r		@Request
// @param:		value	@String
// @retval:		err 	@Error
//
// @usage:		http.SetBody(r, "{key: value}")
func SetBody(r *Request, value string) error {
	r.Body = value
	return nil
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
	return nil, SetBody(r, value)
}

// Exec sends the passed request object.
//
// @callable:	http.Exec
// @param:		r			@Request
// @retval:		response	@String
// @retval:		err 		@Error
//
// @usage:		resp = http.Exec(r)
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
	return Exec(r)
}
