package main

import (
	"log"
	"os"
)

type HttpFileBody struct {
	file *os.File
}

func (fileBody *HttpFileBody) Read(p []byte) (n int, err error) {
	return fileBody.file.Read(p)
}

func (fileBody *HttpFileBody) ContentLength() int {
	fi, err := fileBody.file.Stat()
	if err != nil {
		log.Panic("Could not obtain stat, handle error")
	}

	return int(fi.Size())
}

func (fileBody *HttpFileBody) ContentType() string {
	return "application/octet-stream"
}
