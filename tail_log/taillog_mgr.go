package tail_log

import (
	"fmt"
	"studyGo/logagent/etcd"
)

var tskMgr *tailLogMgr

type tailLogMgr struct {
	logEntry []*etcd.LogEntry
	tskMap map[string]*Tailtask
	newConfChan chan []*etcd.LogEntry
}

// Init : 管理所有的taillog，调用此方法则是开启一个子协程去读取文件发往kafka
//循环每一个日志收集项，创建tailObj(这里tailObj就不能使用全局变量了，因为每加一个或删除一个文件读取路径，相应的就应该在程序中更改)
func Init(LogEntryConf []*etcd.LogEntry){
	tskMgr=&tailLogMgr{
		logEntry: LogEntryConf,
		tskMap: make(map[string]*Tailtask,16),  //一台机器收集16个服务的日志
		newConfChan: make(chan []*etcd.LogEntry),	//定义一个无缓冲通道，如果没有数据就阻塞
	}
	for _,logEntry:=range LogEntryConf{
		//3.2 发往kafka,在init中直接go
		//记录启动的tailTask，方便后续新增修改和删除
		tailObj:=NewTailTask(logEntry.Path,logEntry.Topic)
		pathAndTopic:=fmt.Sprintf("%s_%s",logEntry.Path,logEntry.Topic)
		tskMgr.tskMap[pathAndTopic]=tailObj
		fmt.Println("监听成功:",logEntry.Path)
	}
	go tskMgr.run()
}

//监听自己的newConfChan,有了新的配置过来后就做对应的处理
func (t *tailLogMgr) run(){
	for{
		select {
		case newConf:=<-t.newConfChan:
			// 1. 配置新增
			for _,conf:=range newConf{
				pathAndTopic:=fmt.Sprintf("%s_%s",conf.Path,conf.Topic)
				_,ok:=t.tskMap[pathAndTopic]
				//如果原来就有，就不需要操作
				if ok{
					continue
				}else{
					//没有就新增
					tailObj:=NewTailTask(conf.Path,conf.Topic)
					tskMgr.tskMap[pathAndTopic]=tailObj
				}
			}
			// 2. 配置删除,找出原本存在但是现在不存在的配置项
			for _,confOld:=range t.logEntry{	//从原来的配置中依次拿出配置项
				isDelete:=true
				for _,confNew:=range newConf{	//去新的配置中逐一进行比较，如果存在就更改标志位并且退出内部的循环
					if confOld.Path==confNew.Path&&confOld.Topic==confNew.Topic{
						isDelete=false
						break
					}
				}
				if isDelete{
					//如果存在原配置有但是新配置没有的，则删除原本的
					pathAndTopic:=fmt.Sprintf("%s_%s",confOld.Path,confOld.Topic)
					t.tskMap[pathAndTopic].cancelFunc()
					delete(t.tskMap, pathAndTopic)

				}
			}
		}


		// 3. 配置修改
	}
}

// 向外暴露tskMgr的newConfChan
func NewConfChan() chan<- []*etcd.LogEntry{
	return tskMgr.newConfChan
}