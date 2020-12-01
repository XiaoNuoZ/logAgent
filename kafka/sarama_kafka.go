//往kafka写日志的模块
package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

//基于sarama第三方库开发的kafka client
var(
	client sarama.SyncProducer
)

//kafka连接初始化
func Init(addrs []string)(err error){
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
	return
}

//发送消息
func SendMsg(topic,data string){
	//构造一个消息
	msg:=&sarama.ProducerMessage{}
	msg.Topic=topic
	msg.Value=sarama.StringEncoder(data)
	//发送消息
	pid,offset,err:=client.SendMessage(msg)
	if err!=nil{
		fmt.Println("send message faild,err:",err)
		return
	}
	fmt.Printf("pid:%v offset:%v\n",pid,offset)
}
