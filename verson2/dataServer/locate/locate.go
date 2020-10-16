package locate

import (
	"os"
	"strconv"
)

// 用于实际定位对象
func Locate(name string) bool{
	// name为文件名，存在则定位成功
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// 用于监听定位消息
func StartLocate()  {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	// 绑定dataServers exchange
	q.Bind("dataServers")
	// 返回一个Go语言的Channel(用于后续遍历)
	c := q.Consume()
	// 接口服务发送过来，需要做定位的对象名字
	for msg := range c{
		object, e := strconv.Unquote(string(msg.Body))
		if e != nil{
			panic(e)
		}
		// 检查文件是否存在
		if Locate(os.Getenv("STORAGE_ROOT") + "/objects/" + object){
			// 向消息的发送方返回本服务节点的监听地址(表示该对象存在于本服务节点上)
			q.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}
