package main

import (
	"fmt"
	"log"
	"net"
	"io"
)

type HttpSender struct {
	conn net.Conn
}

// func (sender *HttpSender) SendAll(response *HttpResponse) {
// 	outputStr := response.toStringResponse()
// 	log.Printf("Send the following data: \n%s", outputStr)

// 	output, err := sender.conn.Write([]byte(outputStr))

// 	if err != nil {
// 		fmt.Println("Error writing output: ", err.Error())
// 		return
// 	}

// 	log.Printf("sent %d bytes", output)
// }

func (sender *HttpSender) SendAll(response *HttpResponse) {
	if response.body == nil {
		response.body = HttpTextBody{
			text: "",
		}
	}

	sender.SendMetadata(response)
	sender.SendBody(response)
}

func (sender *HttpSender) SendMetadata(response *HttpResponse) {
	statusStr := fmt.Sprintf("HTTP/1.1 %s\r\n", response.code)

	headersStr := response.headersToString()

	metadata := fmt.Sprintf("%s%s\r\n", statusStr, headersStr)

	log.Printf("Sending response metadata: \n%s", metadata)
	output, err := sender.conn.Write([]byte(metadata))

	if err != nil {
		fmt.Println("Error sending response metadata: ", err.Error())
		panic(err)
	}

	log.Printf("Response metadata has been sent (%d bytes)", output)
}

func (sender *HttpSender) SendBody(response *HttpResponse) {
	return sender.SendBodyAsText(response.body)
}

func (sender *HttpSender) SendBodyAsText(body IHttpStringBody) {
	bodyStr := body.ToString()

	log.Printf("Sending response body: \n%s", bodyStr)
	output, err := sender.conn.Write([]byte(bodyStr))

	if err != nil {
		fmt.Println("Error sending response body: ", err.Error())
		panic(err)
	}

	log.Printf("Response body has been sent (%d bytes)", output)
}
