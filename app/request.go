package main

import (
	"strings"
)

type HttpRequest struct {
	path    string
	body    string
	headers map[string]string
}

func (request HttpRequest) GetHeader(name string) string {
	return request.headers[strings.ToLower(name)]
}

func newRequest(rawRequest string) HttpRequest {
	requestPieces := strings.Split(rawRequest, "\r\n\r\n")
	requestMetadata := requestPieces[0]
	requestBody := requestPieces[1]
	requestMetadataPieces := strings.Split(requestMetadata, "\r\n")

	requestHeadersRequestMetadataPieces := requestMetadataPieces[1:]
	requestHeaders := make(map[string]string)
	for _, rawHeader := range requestHeadersRequestMetadataPieces {
		headerPair := strings.Split(rawHeader, ":")
		headerName := strings.ToLower(strings.Trim(headerPair[0], " "))
		requestHeaders[headerName] = strings.Trim(headerPair[1], " ")
	}

	requestLine := requestMetadataPieces[0]
	requestPath := strings.Split(requestLine, " ")[1]
	return HttpRequest{
		path:    requestPath,
		body:    requestBody,
		headers: requestHeaders,
	}
}
