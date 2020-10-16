package heartbeat

import (
	"../../rabbitmq"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

//整个包内可见
var dataServers = make(map[string]time.Time)
//自带的标准包，互斥访问
var mutex sync.Mutex

func ListenHeartbeat(){
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	//绑定apiServers exchange
	q.Bind("apiServers")
	//返回一个Go语言的Channel(用于后续遍历)
	c := q.Consume()
	go removeExpiredDataServer()
	for msg := range c {
		dataServer, e := strconv.Unquote(string(msg.Body))
		if e != nil{
			panic(e)
		}
		mutex.Lock()
		// key-数据服务节点的监听地址
		dataServers[dataServer] = time.Now()
		// value-收到消息的时间
		mutex.Unlock()
	}

}

//每5s扫描一遍dataServers，并清除其中超过10s没有收到消息的数据服务节点
func removeExpiredDataServer(){
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for s, t := range dataServers{
			if t.Add(10 * time.Second).Before(time.Now()){
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}
}

//遍历dataServers并返回
func GetDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()
	ds := make([]string, 0)
	for s, _ := range dataServers{
		ds = append(ds, s)
	}
	return ds
}

//在当前所有的数据服务节点中随机选出一个结点并返回。
//若结点为空，则返回空字符串
func ChooseRandomDataServer() string{
	ds := GetDataServers()
	n := len(ds)
	if n == 0{
		return ""
	}
	return ds[rand.Intn(n)]
}






































