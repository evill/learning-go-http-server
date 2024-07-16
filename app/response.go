package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type HttpResponse struct {
	code    string
	body    IHttpBody
	headers map[string]string
	sender  HttpSender
}

// func (response HttpResponse) toStringResponse() string {
// 	statusStr := fmt.Sprintf("HTTP/1.1 %s\r\n", response.code)

// 	headersStr := response.headersToString()
// 	bodyStr := response.bodyToString()

// 	return fmt.Sprintf("%s%s\r\n%s", statusStr, headersStr, bodyStr)
// 	//
// }

func (response HttpResponse) headersToString() string {
	if response.headers == nil {
		return ""
	}

	headersStrArray := make([]string, 0, len(response.headers))
	for headerName, headerValue := range response.headers {
		headersStrArray = append(headersStrArray, fmt.Sprintf("%s: %s", headerName, headerValue))
	}

	return strings.Join(headersStrArray[:], "\r\n") + "\r\n"
}

// func (response HttpResponse) bodyToString() string {
// 	if response.body == nil {
// 		return ""
// 	}

// 	return response.body.ToString()
// }

func (response *HttpResponse) SetHeader(name string, value string) *HttpResponse {
	if response.headers == nil {
		response.headers = make(map[string]string)
	}
	formattedName := strings.ToLower(strings.Trim(name, " "))
	response.headers[formattedName] = (strings.Trim(value, " "))

	return response
}

func (response *HttpResponse) Status404() *HttpResponse {
	response.code = "404 Not Found"
	return response
}

func (response *HttpResponse) Status500() *HttpResponse {
	response.code = "500 Internal Server Error"
	return response
}

func (response *HttpResponse) Status200() *HttpResponse {
	response.code = "200 OK"
	return response
}

func (response *HttpResponse) Body(body IHttpBody) *HttpResponse {
	response.body = body
	// response.SetHeader("Content-Length", fmt.Sprintf("%d", body.ContentLength()))
	// response.SetHeader("Content-Type", body.ContentType())

	return response
}

func (response *HttpResponse) Text(text string) {
	body := HttpTextBody{text: text}
	response.Body(body)
	response.Send()
}

func (response *HttpResponse) LocalFile(pathToFile string) {
	file, err := os.Open(pathToFile)

	if err != nil {
		log.Panicf("Couldn't open a file '%s'.", pathToFile)
	}

	defer file.Close()

	body := HttpFileBody{
		file: file,
	}

	response.Body(body)
	response.Send()
}

func (response *HttpResponse) Send() {
	response.sender.SendAll(response)
}

type IHttpBody interface {
	// ToString() string
	ContentLength() int
	ContentType() string
}

type IHttpStringBody interface {
	ToString() string
}
