package locate

import (
	"../../rabbitmq"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request)  {
	m := r.Method
	if m != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	info := Locate(strings.Split(r.URL.EscapedPath(), "/")[2])
	if len(info) == 0{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, _ := json.Marshal(info)
	w.Write(b)
}

// name为需要定位对象的名字
func Locate(name string) string {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	//群发这个对象名字的定位信息
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		//超时机制，避免无限等待
		time.Sleep(time.Second)
		// if - 1s后没有任何反馈，消息队列关闭，收到长度为0的消息
		// else - 返回该数据服务节点的监听地址
		q.Close()
	}()
	msg := <-c
	s, _ := strconv.Unquote(string(msg.Body))
	return s
}

//通过检查Locate结果是否为空字符串来判定对象是否存在
func Exit(name string) bool {
	return Locate(name) != ""
}
































