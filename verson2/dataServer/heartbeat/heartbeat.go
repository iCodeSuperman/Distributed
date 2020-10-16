package heartbeat

import (
	"../../rabbitmq"
	"os"
	"time"
)

// 注意，本函数在goroutine中执行，即使不返回也不会影响其他功能
func StartHeartbeat() {
	// 创建了一个rabbitmq.RabbitMQ结构体
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	//死循环
	for {
		// 向apiServer exchange发送本节点的监听地址
		// os.Getenv("key") 获取系统环境变量
		q.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(5 * time.Second)
	}
}
