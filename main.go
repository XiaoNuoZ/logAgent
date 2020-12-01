package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"studyGo/logagent/conf"
	"studyGo/logagent/kafka"
	"studyGo/logagent/tail_log"
	"time"
)

//logAgen入口程序

func run(){
	//1.读取日志
	for{
		select {
			//从ReadLog返回的管道中获取每行的内容，如果没获取到就休息一秒
			case line:=<-tail_log.ReadLog():
				//2.发送到Kafka
				kafka.SendMsg(cfg.KafkaConf.Topic,line.Text)
		default:
			time.Sleep(time.Second)
		}
	}

}

//接受配置文件的key和value
var cfg =new(conf.LogAgentConf)

func main(){
	//cfg, err := ini.Load("./conf/config.ini")
	//if err != nil {
	//	fmt.Printf("Fail to read file: %v", err)
	//	os.Exit(1)
	//}
	//0.读取配置文件,引入第三方包ini
	err:=ini.MapTo(cfg,"./conf/config.ini")
	if err!=nil{
		fmt.Println("ini-config init faild,err:",err)
		return
	}

	//1.初始化kafka连接
	err=kafka.Init([]string{cfg.KafkaConf.Address})
	if err!=nil{
		fmt.Println("kafka init faild,err:",err)
		return
	}
	fmt.Println("kafka初始化成功")
	//2.打开日志文件准备收集日志
	err=tail_log.Init(cfg.TaillogConf.Filename)
	if err!=nil{
		fmt.Println("taillog init faild,err:",err)
		return
	}
	fmt.Println("taillog初始化成功")
	run()
}

