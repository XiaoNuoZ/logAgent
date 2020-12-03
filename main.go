package main

import (
	//"context"
	"fmt"
	"gopkg.in/ini.v1"
	"studyGo/logagent/conf"
	"studyGo/logagent/etcd"
	"studyGo/logagent/kafka"
	//"studyGo/logagent/tail_log"
	"time"
)

//logAgen入口程序

//func run(){
//	//1.读取日志
//	for{
//		select {
//			//从ReadLog返回的管道中获取每行的内容，如果没获取到就休息一秒
//			case line:=<-tail_log.ReadLog():
//				//2.发送到Kafka
//				kafka.SendMsg(cfg.KafkaConf.Topic,line.Text)
//		default:
//			time.Sleep(time.Second)
//		}
//	}
//
//}

//接受配置文件的key和value
var cfg =new(conf.LogAgentConf)

func main(){
	//0.读取配置文件,引入第三方包ini
	err:=ini.MapTo(cfg,"./conf/config.ini")
	if err!=nil{
		fmt.Println("ini-config init faild,err:",err)
		return
	}

	//1.1 初始化kafka连接
	err=kafka.Init([]string{cfg.KafkaConf.Address})
	if err!=nil{
		fmt.Println("kafka init faild,err:",err)
		return
	}
	fmt.Println("kafka初始化成功")
	//1.1 初始化etcd连接
	err=etcd.Init(cfg.Etcd.Address,time.Duration(cfg.Etcd.Timeout)*time.Second)
	if err!=nil{
		fmt.Println("etcd init faild,err:",err)
		return
	}
	fmt.Println("etcd初始化成功")

	LogAgentConf,err:= etcd.GetConf("fileName")
	if err!=nil{
		fmt.Println("get etcd faild,err:",err)
		return
	}
	fmt.Println("get conf from etcd success")
	fmt.Println(LogAgentConf)
	//2.1 从etcd中获取日志收集项的配置信息，这样的好处是可以进行热更新，每往etcd中put一个新的文件路径，程序就可以同时读取它的日志内容
	//派一个哨兵监视fileName key的变化
	//logAgentConf,err:=etcd.GetConf("fileName")
	//if err != nil {
	//	fmt.Println("etcd.GetConf faild,err:",err)
	//	return
	//}
	//rch:=etcd.Cli.Watch(context.Background(),"fileName")
	////从通道尝试取值，range会在没值的时候阻塞，有值后执行内部代码，然后循环回来继续阻塞等待下一个值
	//for wresp:=range rch{
	//	for _,ev:=range wresp.Events{
	//		//type可以取到类型，是删除修改还是新增
	//		fmt.Println("type:",string(ev.Type),"----key:",string(ev.Kv.Key),"----value:",string(ev.Kv.Value))
	//		err=tail_log.Init(string(ev.Kv.Value))
	//		if err!=nil{
	//			fmt.Println("taillog init faild,err:",err)
	//			return
	//		}
	//		fmt.Println("taillog初始化成功")
	//	}
	//}
	//run()
}

