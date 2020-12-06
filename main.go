package main

import (
	//"context"
	"fmt"
	"gopkg.in/ini.v1"
	"studyGo/logagent/conf"
	"studyGo/logagent/etcd"
	"studyGo/logagent/kafka"
	"studyGo/logagent/tail_log"
	"sync"

	"time"
)


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
	err=kafka.Init([]string{cfg.KafkaConf.Address},cfg.KafkaConf.MaxSize)
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
	//2. 获取日志文件路径
	//拼接ip,用于区分每个服务器的key
	//ipStr,err:=conf.GetOutboundIP()
	//cfg.Etcd.FilePathKey=fmt.Sprintf(cfg.Etcd.FilePathKey,ipStr)
	LogEntryConf,err:= etcd.GetConf(cfg.Etcd.FilePathKey)
	if err!=nil{
		fmt.Println("get etcd faild,err:",err)
		return
	}

	//3. 收集日志发往kafka
	tail_log.Init(LogEntryConf)
	//因为NewConfChan()访问了tskMgr的newConfChan,这个channel是在init中初始化的，因此要把init()放在之前
	//4. 派一个后台的哨兵去监视日志收集项的变化（有变化及时通知程序，实现热加载配置）
	var wg sync.WaitGroup
	wg.Add(1)
	newConfChan:=tail_log.NewConfChan()		//从taillog包中获取对外暴露的通道
	go etcd.WatchConf(cfg.Etcd.FilePathKey,newConfChan)		//哨兵发现最新的配置信息会通知上面的通道
	wg.Wait()
}