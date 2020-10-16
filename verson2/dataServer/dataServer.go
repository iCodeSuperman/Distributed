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
	// 这里与单机版的main函数，多了两个goroutine
	// 每5s向apiServers exchange发送本服务节点的监听地址
	go heartbeat.StartHeartbeat()
	// 监听定位消息
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
