package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"os"
	"strings"
)
const okResponse = "HTTP/1.1 200 OK\r\n%v\r\n\r\n%v\r\n"
const notFoundResponse = "HTTP/1.1 404 Not Found\r\n\r\n"
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	
	
	for  {
		connection , err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}


		_ , err = handleRequest(connection)

		if err != nil {
			fmt.Println("Error handling request: ", err.Error())
		}
	}
}


func handleRequest(conn net.Conn) (n int , err error) {
	defer conn.Close()
	buffer := make([]byte , 1024) 
	
	reqLen , err := conn.Read(buffer)
	if(err != nil){
		fmt.Println("err while reading req:" , err)
	}
	req := string(buffer[:reqLen])

	lines := strings.Split(req,"\r\n")

	path := strings.Split(lines[0]," ")[1]


	var response string
	if path == "/"{
		response = okResponse
	}else if strings.HasPrefix(path,"/echo"){
		body := strings.Split(path,"/echo/")[1]
		bodyLen := len(body)

		headers := []string { "Content-Type: text/plain" , fmt.Sprintf("Content-Length: %v",bodyLen) }
		response = fmt.Sprintf(okResponse , strings.Join(headers,"\r\n"),body)
		fmt.Println(response)
	}else{
		response = notFoundResponse
	}

	n , err = conn.Write([]byte(response))
	return
}


