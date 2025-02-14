package main

import (
	"strings"
)

type HttpRequest struct {
	method  string
	path    string
	body    string
	headers HttpRequestHeaders
	server  *Server
}

func (request HttpRequest) GetHeader(name string) string {
	return request.headers[strings.ToLower(name)]
}

func (request HttpRequest) AceeptsEncoding(name string) bool {
	encodings := request.headers.GetAceeptedEncodings()
	for _, encoding := range encodings {
		if encoding.name == name {
			return true
		}
	}
	return false
}

func newRequest(server *Server, rawRequest string) (request *HttpRequest) {
	request = &HttpRequest{
		server: server,
	}
	requestPieces := strings.Split(rawRequest, "\r\n\r\n")
	requestMetadataPieces := strings.Split(requestPieces[0], "\r\n")
	request.headers = make(HttpRequestHeaders)

	for _, rawHeader := range requestMetadataPieces[1:] {
		headerPair := strings.Split(rawHeader, ":")
		headerName := strings.ToLower(strings.Trim(headerPair[0], " "))
		request.headers[headerName] = strings.Trim(headerPair[1], " ")
	}

	request.method = strings.Split(requestMetadataPieces[0], " ")[0]
	request.path = strings.Split(requestMetadataPieces[0], " ")[1]

	request.body = requestPieces[1]

	return
}
