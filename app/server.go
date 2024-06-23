package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	// Ensure we teardown the server when the program exits
	defer listener.Close()

	for {
		// Block until we receive an incoming connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
		}

		// Handle client connection
		handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	input, err := conn.Read(buf)

	if err != nil {
		fmt.Println("Error reading input: ", err.Error())
		os.Exit(1)
	}

	log.Printf("received %d bytes", input)
	inputStr := string(buf[:input])
	log.Printf("received the following data: \n%s", inputStr)
	requestLine := strings.Split(inputStr, "\n")[0]
	requestPath := strings.Split(requestLine, " ")[1]

	var responseCode string
	if requestPath == "/abcdefg" {
		responseCode = "200 OK"
	} else {
		responseCode = "404 Not Found"
	}

	response := fmt.Sprintf("HTTP/1.1 %s\r\n\r\n", responseCode)
	output, err := conn.Write([]byte(response))

	if err != nil {
		fmt.Println("Error writing output: ", err.Error())
		os.Exit(1)
	}

	log.Printf("sent %d bytes", output)
}
