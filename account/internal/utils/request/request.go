package request

import (
	"errors"
	"net/http"

	"github.com/parnurzeal/gorequest"
)

// This block defines the global parameters.
const (
	ZeroString = ""
	ZeroInt    = 0
)

// This block defines a variety of Method.
const (
	MethodGet     = "GET"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH"
	MethodDelete  = "DELETE"
	MethodHead    = "HEAD"
	MethodOptions = "OPTIONS"
)

// This block defines the parameters necessary for Header of HTTP request.
const (
	HeaderAuthorization = "Authorization"
	headerContentType   = "Content-Type"
	headerAccept        = "Accept"
	appJson             = "application/json"

	HeaderParam  = 0
	HeaderValue  = 1
	HeaderLength = 2
)

var ErrRespBody = errors.New("respBody is error")

type Request struct {
	Method string
	Url    string
	Header []string
	Token  string
	Query  string
	Body   interface{}
}

func (r *Request) Send() (string, error) {
	agent := gorequest.New().CustomMethod(r.Method, r.Url)
	if r.Method == MethodPost || r.Method == MethodPatch || r.Method == MethodPut {
		agent = agent.Set(headerContentType, appJson).Set(headerAccept, appJson)
	}

	if len(r.Header) == HeaderLength {
		agent = agent.Set(r.Header[HeaderParam], r.Header[HeaderValue])
	}

	if r.Token != ZeroString {
		agent = agent.Set(HeaderAuthorization, r.Token)
	}

	if r.Query != ZeroString {
		agent = agent.Query(r.Query)
	}

	if r.Body != nil {
		agent = agent.Send(r.Body)
	}

	response, respBody, errs := agent.End()
	if len(errs) != ZeroInt {
		return respBody, errs[ZeroInt]
	}

	if response.StatusCode != http.StatusOK {
		return respBody, ErrRespBody
	}
	return respBody, nil
}
