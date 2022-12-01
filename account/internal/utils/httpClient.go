package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//FormRPCCall sends form urlencoded to a server
func FormRPCCall(url string, authToken string, verb string, payload url.Values) (body []byte, httpstatus int, err error) {
	req, err := http.NewRequest(verb, url, strings.NewReader(payload.Encode()))
	req.Header.Add("accept", "")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if authToken != "" {
		req.Header.Set(strings.ToLower("otp"), authToken)
	}
	lowerCaseHeader := make(http.Header)
	for key, value := range req.Header {
		lowerCaseHeader[strings.ToLower(key)] = value
	}
	req.Header = lowerCaseHeader

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)

	fmt.Println("response Status:", resp.Status)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return body, resp.StatusCode, nil
	}

	switch resp.StatusCode {
	case 400:
		return body, resp.StatusCode, GetValidationError("خطا در مقدار ورودی اطلاعات.")
	case 401:
		return body, resp.StatusCode, GetUnAuthorizedError("Who are you?!")
	case 403:
		return body, resp.StatusCode, GetForbiddenError("No Access.")
	case 404:
		return body, resp.StatusCode, GetNotFoundError("Not Found.")
	}

	return body, resp.StatusCode, errors.New("Server returned not success status code")

}

//RPCCall sends data to a server
func RPCCall(url string, authToken string, verb string, payload []byte) (body []byte, httpstatus int, err error) {
	req, err := http.NewRequest(verb, url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	if authToken != "" {
		req.Header.Set("Authorization", authToken)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println(err)
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)

	//fmt.Println("response Status:", resp.Status)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return body, resp.StatusCode, nil
	}

	switch resp.StatusCode {
	case 400:
		return body, resp.StatusCode, GetValidationError("خطا در مقدار ورودی اطلاعات.")
	case 401:
		return body, resp.StatusCode, GetUnAuthorizedError("Who are you?!")
	case 403:
		return body, resp.StatusCode, GetForbiddenError("No Access.")
	case 404:
		return body, resp.StatusCode, GetNotFoundError("Not Found.")
	}

	return body, resp.StatusCode, errors.New("Server returned not success status code")

}

// RPCCallNONESSL RPCCall sends data to a server
func RPCCallNONESSL(url string, authToken string, verb string, payload []byte) (body []byte, httpstatus int, err error) {
	req, err := http.NewRequest(verb, url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	if authToken != "" {
		req.Header.Set("Authorization", authToken)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println(err)
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)

	//fmt.Println("response Status:", resp.Status)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return body, resp.StatusCode, nil
	}

	switch resp.StatusCode {
	case 400:
		return body, resp.StatusCode, GetValidationError("خطا در مقدار ورودی اطلاعات.")
	case 401:
		return body, resp.StatusCode, GetUnAuthorizedError("Who are you?!")
	case 403:
		return body, resp.StatusCode, GetForbiddenError("No Access.")
	case 404:
		return body, resp.StatusCode, GetNotFoundError("Not Found.")
	}

	return body, resp.StatusCode, errors.New("Server returned not success status code")

}

// SMSRPCCall sending request without tsl secure connection
func SMSRPCCall(url string, authToken string, verb string, payload []byte) (body []byte, httpstatus int, err error) {
	req, err := http.NewRequest(verb, url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	if authToken != "" {
		req.Header.Set("Authorization", authToken)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println(err)
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)

	//fmt.Println("response Status:", resp.Status)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return body, resp.StatusCode, nil
	}

	switch resp.StatusCode {
	case 400:
		return body, resp.StatusCode, GetValidationError("خطا در مقدار ورودی اطلاعات.")
	case 401:
		return body, resp.StatusCode, GetUnAuthorizedError("Who are you?!")
	case 403:
		return body, resp.StatusCode, GetForbiddenError("No Access.")
	case 404:
		return body, resp.StatusCode, GetNotFoundError("Not Found.")
	}

	return body, resp.StatusCode, errors.New("Server returned not success status code")

}

// FormRequest OpenApiCallForm sends data to a server
func FormRequest(path string, verb string, headers map[string]string, payload url.Values) (body []byte, status int) {
	req, err := http.NewRequest(verb, path, strings.NewReader(payload.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, 0
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	return body, resp.StatusCode
}

//JsonRequest sends data to a server
func JsonRequest(url string, verb string, headers map[string]string, payload []byte) (body []byte, status int) {
	req, err := http.NewRequest(verb, url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, 0
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	return body, resp.StatusCode
}

func CustomDecoder(body io.Reader, i interface{}) error {
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&i)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// Catch any syntax errors in the JSON and send an error message
		// which interpolates the location of the problem to make it
		// easier for the client to fix.
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return errors.New(msg)

		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return errors.New(msg)

		// Catch any type errors, like trying to assign a string in the
		// JSON request body to a int field in our Person struct. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return errors.New(msg)

		// Catch the error caused by extra unexpected fields in the request
		// body. We extract the field name from the error message and
		// interpolate it in our custom error message. There is an open
		// issue at https://github.com/golang/go/issues/29035 regarding
		// turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.ReplaceAll(strings.TrimPrefix(err.Error(), "json: unknown field "), "\"", "")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return errors.New(msg)

		// An io.EOF error is returned by Decode() if the request body is
		// empty.
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return errors.New(msg)

		// Catch the error caused by the request body being too large. Again
		// there is an open issue regarding turning this into a sentinel
		// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return errors.New(msg)

		// Otherwise, default to logging the error and sending a 500 Internal
		// Server Error response.
		default:
			log.Println(err.Error())
			return errors.New(err.Error())
		}
	}

	return nil
}
