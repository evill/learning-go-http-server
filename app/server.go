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

	// requestLine := strings.Split(inputStr, "\n")[0]
	// requestPath := strings.Split(requestLine, " ")[1]
	request := newRequest(inputStr)
	response := routeRequest(request)
	outputStr := response.toStringResponse()
	log.Printf("Send the following data: \n%s", outputStr)

	output, err := conn.Write([]byte(outputStr))

	if err != nil {
		fmt.Println("Error writing output: ", err.Error())
		os.Exit(1)
	}

	log.Printf("sent %d bytes", output)
}

type HttpResponse struct {
	code    string
	body    string
	headers map[string]string
}

func (response HttpResponse) toStringResponse() string {
	statusStr := fmt.Sprintf("HTTP/1.1 %s", response.code)

	if response.headers == nil {
		response.headers = make(map[string]string)
	}

	response.headers["Content-Length"] = fmt.Sprintf("%d", len(response.body))
	response.headers["Content-Type"] = "text/plain"
	headersStr := response.headersToString()

	return fmt.Sprintf("%s\r\n%s\r\n\r\n%s", statusStr, headersStr, response.body)
	//
}

func (response HttpResponse) headersToString() string {
	headersStrArray := make([]string, 0, len(response.headers))
	for headerName, headerValue := range response.headers {
		headersStrArray = append(headersStrArray, fmt.Sprintf("%s: %s", headerName, headerValue))
	}

	return strings.Join(headersStrArray[:], "\r\n")
}

type HttpRequest struct {
	path    string
	body    string
	headers map[string]string
}

func (request HttpRequest) GetHeader(name string) string {
	return request.headers[strings.ToLower(name)]
}

func newRequest(rawRequest string) *HttpRequest {
	requestPieces := strings.Split(rawRequest, "\r\n\r\n")
	requestMetadata := requestPieces[0]
	requestBody := requestPieces[1]
	requestMetadataPieces := strings.Split(requestMetadata, "\r\n")

	requestHeadersRequestMetadataPieces := requestMetadataPieces[1:len(requestMetadataPieces)]
	requestHeaders := make(map[string]string)
	for _, rawHeader := range requestHeadersRequestMetadataPieces {
		headerPair := strings.Split(rawHeader, ":")
		headerName := strings.ToLower(strings.Trim(headerPair[0], " "))
		requestHeaders[headerName] = strings.Trim(headerPair[1], " ")
	}

	requestLine := requestMetadataPieces[0]
	requestPath := strings.Split(requestLine, " ")[1]
	return &HttpRequest{
		path:    requestPath,
		body:    requestBody,
		headers: requestHeaders,
	}
}

func routeRequest(request *HttpRequest) *HttpResponse {
	switch {
	case request.path == "/":
		return routeRoot(request)
	case strings.HasPrefix(request.path, "/echo/"):
		return routeEcho(request)
	case strings.HasPrefix(request.path, "/user-agent"):
		return routeUserAgent(request)
	default:
		return route404(request)
	}
}

func routeRoot(request *HttpRequest) *HttpResponse {
	return &HttpResponse{
		code: "200 OK",
		body: "",
	}
}

func routeEcho(request *HttpRequest) *HttpResponse {
	parameter, _ := strings.CutPrefix(request.path, "/echo/")
	return &HttpResponse{
		code: "200 OK",
		body: parameter,
	}
}

func routeUserAgent(request *HttpRequest) *HttpResponse {
	return &HttpResponse{
		code: "200 OK",
		body: request.GetHeader("User-Agent"),
	}
}

func route404(request *HttpRequest) *HttpResponse {
	return &HttpResponse{
		code: "404 Not Found",
		body: "",
	}
}
