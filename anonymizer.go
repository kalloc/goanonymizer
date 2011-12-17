package main

import (
    "io"
    "os"
    "fmt"
    "http"
)

func root(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Index!")
}
func handle_http(w http.ResponseWriter, r *http.Request) {
    var url string;
    var response *http.Response;
    var responseError os.Error;
    fmt.Printf("%s %s %s\n", r.Method, r.RawURL,r.RemoteAddr)
    url  = "http://"+r.URL.Path[6:]
    newRequest,_ := http.NewRequest(r.Method, url, nil) 
    c := &http.Client{} 
    response, responseError = c.Do(newRequest)
    for header := range response.Header {
        w.Header().Add(header, response.Header.Get(header))
    }
    if responseError != nil {
        fmt.Fprintf(w, "pizda\n")
    } else {
        io.Copy(w, response.Body)
    }
}

func main() {
    http.HandleFunc("/http/", handle_http)
    http.HandleFunc("/", root)
    http.ListenAndServe(":8080", nil)
}
