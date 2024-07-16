package main

import (
	"os"
	"log"
)

type HttpFileBody struct {
	file *os.File
}

func (fileBody HttpFileBody) ToString() string {
	return string(fileBody.file)
}

func (fileBody HttpFileBody) ContentLength() int {
	fi, err := fileBody.file.Stat()
	if err != nil {
		log.Panic("Could not obtain stat, handle error")
	}
	
	return fi.Size()
}

func (fileBody HttpFileBody) ContentType() string {
	return "application/octet-stream"
}
