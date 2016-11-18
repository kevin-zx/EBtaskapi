package main

import (
	"./taskResult"
	"fmt"
	"io"
	"log"
	"net/http"
)

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}
func HelloServer2(w http.ResponseWriter, req *http.Request) {
	// fmt.Println(req.URL)
	// fmt.Println(req.PostForm)
	// `area`,`device`
	// fmt.Println(req.URL.Query().Get("task_id", "success_status", "ip", "port", "elapsed"))

}
func main() {
	http.HandleFunc("/hello", HelloServer)
	http.HandleFunc("/hello2", HelloServer2)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
