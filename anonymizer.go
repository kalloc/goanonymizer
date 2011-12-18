package main

import (
    "io"
    "strings"
    "io/ioutil"
    "os"
    "fmt"
    "http"
)

func root(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Index!")
}
var ReplacemendContentType = map[string] bool {
    "text/html": true,
    "text/javascript": true,
    "application/x-javascript": true,
    "application/javascript": true,
    "text/plain": true,
    "text/css": true,
}

func handle_http(responseWrite http.ResponseWriter, request *http.Request) {
    var proxyResponse *http.Response;
    var proxyResponseError os.Error;
    fmt.Printf("%s %s %s %s\n", request.Method, request.RawURL, request.RemoteAddr)
    url := "http://"+request.URL.Path[6:]
    proxyRequest,_ := http.NewRequest(request.Method, url, nil) 
    proxy := &http.Client{} 
    proxyResponse, proxyResponseError = proxy.Do(proxyRequest)
    if proxyResponseError != nil {
        http.NotFound(responseWrite, request)
        return
    } 
    for header := range proxyResponse.Header {
        responseWrite.Header().Add(header, proxyResponse.Header.Get(header))
    }
    contentType := strings.Split(proxyResponse.Header.Get("Content-Type"), ";")[0]
    if proxyResponseError != nil {
        fmt.Fprintf(responseWrite, "pizda\n")
    } else if ReplacemendContentType[contentType] {
        body,_ := ioutil.ReadAll(proxyResponse.Body)
        defer proxyResponse.Body.Close()
        bodyString := strings.Replace(string(body), "http://", string("http://"+request.Host+"/http/"), -1)
        //bodyString = strings.Replace(bodyString, "https://", string("http://"+request.Host+"/https/"), -1)
        responseWrite.Write([]byte(bodyString))

    } else {
        io.Copy(responseWrite, proxyResponse.Body)
    }
}

func main() {
    http.HandleFunc("/http/", handle_http)
    http.HandleFunc("/", root)
    http.ListenAndServe(":8080", nil)
}
