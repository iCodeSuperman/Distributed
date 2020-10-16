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
	/**
	1. http.HandleFunc注册HTTP处理函数objects.Handler，若有客户端访问本机的HTTP服务且URL
	   以"/objects/"开头，那么该请求将由objects.Handler负责处理。
	2. 处理函数注册成功后，调用http.ListenAndServe正式开始监听端口
	3. objects.Handler
	 */
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
