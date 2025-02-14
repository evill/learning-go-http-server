package main

import (
	"fmt"
	"strings"
)

type HttpHeaders map[string]string

type HttpResponseHeaders HttpHeaders

func (headers HttpResponseHeaders) String() string {
	if len(headers) == 0 {
		return ""
	}

	headersStrArray := make([]string, 0, len(headers))
	for headerName, headerValue := range headers {
		headersStrArray = append(headersStrArray, fmt.Sprintf("%s: %s", headerName, headerValue))
	}

	return strings.Join(headersStrArray[:], "\r\n") + "\r\n"
}

type HttpRequestHeaders HttpHeaders

func (headers *HttpRequestHeaders) Get(key string) string {
	value, ok := (*headers)[strings.ToLower(key)]

	if ok {
		return value
	}

	return (*headers)[key]
}

type AcepptedEcoding struct {
	name    string
	quality float64
}

func (headers HttpRequestHeaders) GetAceeptedEncodings() []AcepptedEcoding {
	accpetedEncoding := headers.Get("Accept-Encoding")
	// log.Printf("Accepted encoding: %w", headers)

	if accpetedEncoding == "" {
		return []AcepptedEcoding{}
	}

	encodingsData := strings.Split(accpetedEncoding, ",")
	encodings := make([]AcepptedEcoding, 0, len(encodingsData))
	for _, encoding := range encodingsData {
		encoding = strings.TrimSpace(encoding)
		if encoding == "" {
			continue
		}
		encodingParts := strings.Split(encoding, ";")
		encodingName := strings.TrimSpace(encodingParts[0])
		if encodingName == "" {
			continue
		}
		encodingQuality := 1.0
		if len(encodingParts) > 1 {
			qValue := strings.Split(encodingParts[1], "=")[1]
			fmt.Sscanf(qValue, "%f", &encodingQuality)
		}
		encodings = append(encodings, AcepptedEcoding{encodingName, encodingQuality})
	}
	return encodings
}
