package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

type HttpSender struct {
	conn net.Conn
}

func (sender *HttpSender) SendAll(response *HttpResponse) {
	var bodyToSend IHttpBody

	if response.body == nil {
		bodyToSend = &HttpTextBody{
			text: "",
		}
	} else {
		bodyToSend = response.body
	}

	if response.request.AceeptsEncoding("gzip") {
		log.Println("Compressing body...")
		response.SetHeader("Content-Encoding", "gzip")
		bodyToSend = NewCompressedBody(bodyToSend)
	}

	log.Println("Sending response...")
	statusStr := sender.sendStatus(response)
	fmt.Print(*statusStr)
	headersStr := sender.sendHeaders(response, &bodyToSend)
	fmt.Print(*headersStr)
	bodyStr := sender.sendBody(bodyToSend)
	fmt.Println(*bodyStr)
	log.Println("Response sent.")
}

func (sender *HttpSender) sendStatus(response *HttpResponse) *string {
	statusStr := fmt.Sprintf("HTTP/1.1 %s\r\n", response.code)

	_, err := sender.conn.Write([]byte(statusStr))

	if err != nil {
		log.Panic("Error sending response status: ", err.Error())
	}

	return &statusStr
}

func (sender *HttpSender) sendHeaders(response *HttpResponse, bodyToSend *IHttpBody) *string {
	response.SetHeader("Content-Type", (*bodyToSend).ContentType())

	// Use the body that will actually be sent (might be compressed)
	if sizedBody, ok := (*bodyToSend).(IHttpBodyDefinedLength); ok {
		response.SetHeader("Content-Length", fmt.Sprintf("%d", sizedBody.ContentLength()))
	}

	headersStr := fmt.Sprintf("%s\r\n", response.GetHeaders().String())

	// Debug log all headers
	log.Printf("Sending headers:\n%s", headersStr)

	_, err := sender.conn.Write([]byte(headersStr))

	if err != nil {
		log.Panic("Error sending response headers: ", err.Error())
	}

	return &headersStr
}

func (sender *HttpSender) sendBody(body interface{}) *string {
	switch typedBody := body.(type) {
	case io.Reader:
		return sender.SendBodyAsStream(typedBody)
	case fmt.Stringer:
		log.Println("Sending body as text...")
		return sender.SendBodyAsText(typedBody)
	default:
		log.Panic("Unsupported body type")
	}
	return nil
}

func (sender *HttpSender) SendBodyAsText(body fmt.Stringer) *string {
	bodyStr := body.String()
	_, err := sender.conn.Write([]byte(bodyStr))

	if err != nil {
		log.Panic(err)
	}

	return &bodyStr
}

func (sender *HttpSender) SendBodyAsStream(body io.Reader) *string {
	// Create custom buffer with specific size
	buf := make([]byte, 1024)

	_, err := io.CopyBuffer(sender.conn, body, buf)
	if err != nil {
		log.Panic(err)
	}

	bodyStr := string(buf)

	return &bodyStr
}
