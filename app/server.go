package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"flag"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type HttpResponse struct {
	StatusCode int
	StatusText string
	Headers map[string]string
	Body    string
}
type HttpRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
}
func parseRequest(request string) (req HttpRequest , err error){
	requestParts := strings.Split(request,"\r\n\r\n")
	headersAndFirstLine := strings.Split(requestParts[0],"\r\n") 
	methodAndPath := strings.Split(headersAndFirstLine[0]," ")
	fmt.Println("mmmmmm",methodAndPath)
	method := methodAndPath[0]
	path := methodAndPath[1]
	fmt.Println("pppppppp",len(headersAndFirstLine))
	headers := []string{}
	if(len(headersAndFirstLine) > 1){
		headers = headersAndFirstLine[1:]
	}
	
	headersMap := make(map[string]string , len(headers)) 
	for _ , v := range headers {
		headerKeyValuePair := strings.Split(v,":")
		headersMap[headerKeyValuePair[0]] = strings.TrimSpace(headerKeyValuePair[1]) 
	}


	var body string 
	if(len(requestParts) >= 2){
		body = requestParts[1]
	}else {
		body = ""
	}


	req.Method = method
	req.URL = path
	req.Body = body
	req.Headers = headersMap
	return 
}

func stringfyResponse(res HttpResponse) (resText string){
	headers := []string{}
	for k ,v := range res.Headers {
		headers = append(headers, fmt.Sprintf("%v: %v",k,v))
	}
	
	resText = fmt.Sprintf("HTTP/1.1 %v %v\r\n%v\r\n\r\n%v",
					res.StatusCode,res.StatusText,strings.Join(headers,"\r\n"),res.Body) 
	return
}

func OK(conn net.Conn,params ...interface{}){
	defer conn.Close()
	var res HttpResponse 
	res.StatusCode = 200;
	res.StatusText = "OK"
	for index, val := range params{
		switch index {
            case 0: //the first mandatory param
               headers , _ := val.(map[string]string)
			   res.Headers = headers
            case 1: // age is optional param
                body, _ := val.(string)
				res.Body = body
        }
	}
	response := stringfyResponse(res)
	_ , err := conn.Write([]byte(response))
	if(err != nil){
		fmt.Println("err : ",err)
	}
}

func NotFound(conn net.Conn,params ...interface{}){
	defer conn.Close()
	var res HttpResponse 
	res.StatusCode = 404;
	res.StatusText = "Not Found"
	for index, val := range params{
		switch index {
            case 0: //the first mandatory param
               headers , _ := val.(map[string]string)
			   res.Headers = headers
            case 1: // age is optional param
                body, _ := val.(string)
				res.Body = body
        }
	}
	response := stringfyResponse(res)
	_ , err := conn.Write([]byte(response))
	if(err != nil){
		fmt.Println("err : ",err)
	}
}

var DIRFLAG = ""
func main() {
	var dirFlag = flag.String("directory", ".", "directory to serve files from")
 	flag.Parse()
	DIRFLAG = *dirFlag
	
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

		
		
		go handleRequest(connection)

		if err != nil {
			fmt.Println("Error handling request: ", err.Error())
		}
	}
}


func handleRequest(conn net.Conn) {
	buffer := make([]byte , 1024) 
	reqLen , err := conn.Read(buffer)
	if(err != nil){
		fmt.Println("err while reading req:" , err)
	}
	req := string(buffer[:reqLen])
	fmt.Println("rrrrrrrrr",req)
	request , err := parseRequest(req);

	headers := map[string]string{}
	//var response string
	if request.URL == "/"{
		 OK(conn)
	}else if strings.HasPrefix(request.URL,"/echo"){
		body := strings.TrimPrefix(request.URL,"/echo/")
		bodyLen := len(body)
		headers["Content-Type"] = "text/plain"
		headers["Content-Length"] = strconv.Itoa(bodyLen)
		OK(conn,headers,body)
	}else if request.URL == "/user-agent"{
		userAgent := request.Headers["User-Agent"]
		bodyLen := len(userAgent)
		headers["Content-Type"] = "text/plain"
		headers["Content-Length"] = strconv.Itoa(bodyLen)
		OK(conn,headers,userAgent)
	}else if strings.HasPrefix(request.URL,"/files"){
		fileName := strings.TrimPrefix(request.URL,"/files/")
		filePath := filepath.Join(DIRFLAG,fileName)
		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			NotFound(conn)
		}else{
			
			content, err := os.ReadFile(filePath)
			if(err != nil){
				fmt.Println("err: " , err)
			}
			body := string(content)
			bodyLen := len(body)
			headers["Content-Length"] = strconv.Itoa(bodyLen)
			headers["Content-Type"] = "application/octet-stream"
			OK(conn,headers,body)
		}
	}else{
		NotFound(conn)
	}

	return
}


