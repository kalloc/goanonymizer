package main

import (
    "io"
    "regexp"
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

func replace_url(body string, requestHost string, proxyHost string) string {
    reCSS,_ := regexp.Compile("url\\((http://|https://|/)[^\\)]+\\)")
    reSRCandHREF, _ := regexp.Compile("(src|href)=((\"(http://|/|https://)[^\"]+\")|('(https://|/|http://)[^']+')|((http://|/|https://)[^ ]+))")

    FromURLtoAnonURL := func (url string) string {
        if strings.HasPrefix(url, "//") {
            url = requestHost+url[2:]
        } else if strings.HasPrefix(url, "/") {
            url = requestHost+proxyHost+url
        } else if strings.HasPrefix(url, "https://") {
            url = "/https/"+url[len("https://"):]
        } else if strings.HasPrefix(url, "http://") {
            url = "/http/"+url[len("http://"):]
        }
        return url
    }
    fnCSS := func(s string) string { 
        url := s[4:]
        if url[0] == '"' || url[0] == '\'' {
            url = url[1:len(url)-1]
        }
        return "url("+FromURLtoAnonURL(url[:len(url)-1])+")"
    }
    fnSRCandHREF := func(s string) string { 
        var isHREF bool
        var url string
        if strings.HasPrefix(s, "src=") {
            url = s[4:]
            isHREF = false
        } else {
            url = s[5:]
            isHREF = true
        }
        if url[0] == '"' || url[0] == '\'' {
            url = url[1:len(url)-1]
        }
        if isHREF{
            return "href=\""+FromURLtoAnonURL(url)+"\""
        } 
        return "src=\""+FromURLtoAnonURL(url)+"\""
    }
    body = reCSS.ReplaceAllStringFunc(body, fnCSS)
    body = reSRCandHREF.ReplaceAllStringFunc(body, fnSRCandHREF)
    body = strings.Replace(body, "\"https://", "\"/https/", -1)
    body = strings.Replace(body, "\"http://", "\"/http/", -1)
    body = strings.Replace(body, "'https://", "'/https/", -1)
    body = strings.Replace(body, "'http://", "'/http/", -1)
    return body

}

func handle_http(responseWrite http.ResponseWriter, request *http.Request) {
    var proxyResponse *http.Response;
    var proxyResponseError os.Error;

    proxyHost := strings.Split(request.URL.Path[6:], "/")[0]

    fmt.Printf("%s %s %s\n", request.Method, request.RawURL, request.RemoteAddr)
    //TODO https 
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
        bodyString := replace_url(string(body), "/http/", proxyHost)
        
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
