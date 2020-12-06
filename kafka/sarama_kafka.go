//往kafka写日志的模块
package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

type logData struct {
	topic string
	data string
}

//基于sarama第三方库开发的kafka client
var(
	client sarama.SyncProducer
	logDataChan chan *logData	//存结构体指针比存储结构体本身所耗费的内存更少（指针只有十六位）
)

//kafka连接初始化
func Init(addrs []string,maxSize int)(err error){
	config:=sarama.NewConfig()
	//tailf包的使用,Producer为生产者
	config.Producer.RequiredAcks=sarama.WaitForAll	//发送完数据需要leader和follow都确认
	config.Producer.Partitioner=sarama.NewRandomPartitioner //选择生产者写往Kafka的分区模式（这里是轮询模式）
	config.Producer.Return.Successes=true//成功交付的消息将在success channel返回

	//连接kafka
	client,err=sarama.NewSyncProducer(addrs,config)
	if err!=nil{
		fmt.Println("producer closed,err:",err)
		return err
	}
	logDataChan=make(chan *logData,maxSize)
	//起一个后台任务等待通道中有数据然后发送给kafka，如果起多个协程的话，日志顺序可能会混乱
	go sendMsg()
	return
}

// 给外部暴露的一个函数，该函数只把日志数据发送到一个内部的channel中
func SendToChan(topic,data string){
	logDataChan<-&logData{topic: topic,data: data}
}

//通过通道往Kafka发送消息
func sendMsg(){
	for  {
		select {
		case logMsg:=<-logDataChan:
			//构造一个消息
			msg:=&sarama.ProducerMessage{}
			msg.Topic=logMsg.topic
			msg.Value=sarama.StringEncoder(logMsg.data)
			//发送消息
			pid,offset,err:=client.SendMessage(msg)
			if err!=nil{
				fmt.Println("send message faild,err:",err)
				return
			}
			fmt.Printf("pid:%v offset:%v\n",pid,offset)
		default:
			time.Sleep(time.Millisecond*50)
		}
	}
}

