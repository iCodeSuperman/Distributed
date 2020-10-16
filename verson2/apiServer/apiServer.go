package main

import (
	"./heartbeat"
	"./locate"
	"./objects"
	"log"
	"net/http"
	"os"
)


func main() {
	go heartbeat.ListenHeartbeat()
	//REST接口 处理URL以/objects/开头的对象请求
	http.HandleFunc("/objects/", objects.Handler)
	//locate功能 处理URL以/locate/开头的定位请求
	http.HandleFunc("/locate/", locate.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
