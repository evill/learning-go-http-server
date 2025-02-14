package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
)

type CompressedBody struct {
	origin IHttpBody
	text   io.Reader
	length int
}

func (body *CompressedBody) Read(p []byte) (n int, err error) {
	return body.text.Read(p)
}

func (body *CompressedBody) ContentLength() int {
	return body.length
}

func (body *CompressedBody) ContentType() string {
	return body.origin.ContentType()
}

func NewCompressedBody(body IHttpBody) IHttpBody {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	switch typedBody := (body).(type) {
	case io.Reader:
		_, err := io.Copy(zw, typedBody)
		if err != nil {
			log.Panicf("Couldn't compress request content via gzip: %s", err)
		}
	case fmt.Stringer:
		_, err := zw.Write([]byte(typedBody.String()))
		if err != nil {
			log.Panicf("Couldn't compress request content via gzip: %s", err)
		}
	default:
		log.Panic("Unsupported body type for compression")
	}

	// Close writer before getting size
	if err := zw.Close(); err != nil {
		log.Panicf("Couldn't close gzip writer: %s", err)
	}

	log.Printf("Compressed body size: %d", buf.Len())

	return &CompressedBody{
		origin: body,
		text:   &buf,
		length: buf.Len(),
	}
}
