package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"
	"path"
	"strings"
)

type ServerConfig struct {
	filesDirectory *string
	port           *int
}

type Server struct {
	config ServerConfig
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	config := ServerConfig{
		filesDirectory: flag.String("directory", "", "Directory with files for endpoint /files"),
		port:           flag.Int("port", 4221, "Port to listen on"),
	}

	flag.Parse()

	server := Server{
		config: config,
	}

	startServer(&server)
}

func startServer(server *Server) {
	port := *(server.config.port)
	address := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("Failed to bind to port %d\n", port)
		os.Exit(1)
	}

	fmt.Printf("Listening %s...\n", address)

	// Ensure we teardown the server when the program exits
	defer listener.Close()

	for {
		// Block until we receive an incoming connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
		}

		// Handle client connection
		go handleConn(conn, *server)
	}
}

func handleConn(conn net.Conn, server Server) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Unhandled error in connection: %v", err)
		}
	}()
	defer conn.Close()

	buf := make([]byte, 1024)
	input, err := conn.Read(buf)

	if err != nil {
		fmt.Println("Error reading input: ", err.Error())
		return
	}

	inputStr := string(buf[:input])
	log.Printf("Received request with %d bytes: \n%s", input, inputStr)

	request := newRequest(&server, inputStr)
	sender := HttpSender{conn: conn}

	response := &HttpResponse{
		sender:  sender,
		request: request,
	}

	routeRequest(request, response)
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
	if request.method == "GET" {
		getFileRoute(request, response)
	} else if request.method == "POST" {
		postFileRoute(request, response)
	} else {
		response.Status404().Text("Method not allowed")
	}
}

func getFile(request *HttpRequest, response *HttpResponse) {
	fileName, _ := strings.CutPrefix(request.path, "/files/")

	if fileName == "" {
		response.Status404().Text("Name of file is not passed in URL")
		return
	}

	filesDirectory := *request.server.config.filesDirectory
	files, err := os.ReadDir(filesDirectory)
	if err != nil {
		log.Print(err)
		response.Status500().Text("File server feature is not available!")
		return
	}

	var targetFile fs.DirEntry
	for _, file := range files {
		if file.Name() == fileName && !file.IsDir() {
			targetFile = file
			break
		}
	}

	if targetFile == nil {
		log.Printf("Requested file '%s' doesn't exists in folder '$s'", fileName, filesDirectory)
		response.Status404().Text("Requested file not found")
		return
	}

	fullFilePath := path.Join(filesDirectory, fileName)

	response.SetHeader("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", fileName))
	response.Status200().LocalFile(fullFilePath)
}
