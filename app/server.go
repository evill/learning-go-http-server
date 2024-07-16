package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"flag"
	"io/fs"
	"path"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	startServer()
}

func startServer() {
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
		go handleConn(conn)
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

	request := newRequest(inputStr)
	sender := HttpSender{conn: conn}

	response := HttpResponse{
		sender: sender,
	}

	routeRequest(&request, &response)
}

func routeRequest(request *HttpRequest, response *HttpResponse) {
	switch {
	case request.path == "/":
		routeRoot(request, response)
	case strings.HasPrefix(request.path, "/echo/"):
		routeEcho(request, response)
	case strings.HasPrefix(request.path, "/user-agent"):
		routeUserAgent(request, response)
	case strings.HasPrefix(request.path, "/files"):
		routeFiles(request, response)
	default:
		response.Status404().Send()
	}
}

func routeRoot(request *HttpRequest, response *HttpResponse) {
	response.Status200().Send()
}

func routeEcho(request *HttpRequest, response *HttpResponse) {
	parameter, _ := strings.CutPrefix(request.path, "/echo/")
	response.Status200().Text(parameter)
}

func routeUserAgent(request *HttpRequest, response *HttpResponse) {
	response.Status200().Text(request.GetHeader("User-Agent"))
}

func routeFiles(request *HttpRequest, response *HttpResponse) {
	directoryPtr := flag.String("directory", "", "Directory with files for endpoint /files")
	fileName, _ := strings.CutPrefix(request.path, "/files/")

	if fileName == "" {
		response.Status404().Text("Name of file is not passed in URL")
		return
	}

	files, err := os.ReadDir(*directoryPtr)
	if err != nil {
		response.Status500().Text("File server feature is not available!")
		log.Print(err)
	}

	var targetFile fs.DirEntry
	for _, file := range files {
		if file.Name() == fileName && !file.IsDir() {
			targetFile = file
			break
		}
	}

	if targetFile == nil {
		log.Printf("Requested file '%s' doesn't exists in folder '$s'", fileName, *directoryPtr)
		response.Status404().Text("Requested file not found")
	}

	fullFilePath := path.Join(*directoryPtr, fileName)

	response.SetHeader("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", fileName))
	response.Status200().LocalFile(fullFilePath)
}
