package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type HttpResponse struct {
	request *HttpRequest
	code    string
	body    IHttpBody
	headers HttpResponseHeaders
	sender  HttpSender
}

func (response *HttpResponse) SetHeader(name string, value string) *HttpResponse {
	if response.headers == nil {
		response.headers = HttpResponseHeaders{}
	}
	formattedName := strings.ToLower(strings.Trim(name, " "))
	response.headers[formattedName] = (strings.Trim(value, " "))

	return response
}

func (response *HttpResponse) GetHeaders() HttpResponseHeaders {
	return response.headers
}

func (response *HttpResponse) Status(code int, message string) *HttpResponse {
	response.code = fmt.Sprintf("%d %s", code, message)
	return response
}

func (response *HttpResponse) Status409() *HttpResponse {
	response.code = "409 Conflict"
	return response
}

func (response *HttpResponse) Status501() *HttpResponse {
	response.code = "501 Not Implemented"
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

func (response *HttpResponse) Status400() *HttpResponse {
	response.code = "400 Bad Request"
	return response
}

func (response *HttpResponse) Body(body IHttpBody) *HttpResponse {
	response.body = body
	return response
}

func (response *HttpResponse) Text(text string) {
	body := HttpTextBody{text: text}
	response.Body(&body)
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

	response.Body(&body)
	response.Send()
}

func (response *HttpResponse) Send() {
	response.sender.SendAll(response)
}

type IHttpBody interface {
	ContentType() string
}

type IHttpBodyDefinedLength interface {
	ContentLength() int
}
