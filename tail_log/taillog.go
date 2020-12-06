//从各种服务中读取最新日志，引用第三方包tailf
package tail_log

import (
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	_ "github.com/hpcloud/tail"
	"studyGo/logagent/kafka"
	"time"
)


//TailTask : 一个具体的日志收集的任务（因为每个日志都需要对应一个*tail.Tail，不能再使用全局定义）
type Tailtask struct {
	Path string
	Topic string
	Instance *tail.Tail
	//为了实现退出t.RunSendMsg()
	ctx context.Context
	cancelFunc context.CancelFunc
}

//创建TailTask的构造函数，初始化一个包含了path,topic和*tail.Tail的结构体，这样子就可以通过开子携程来同时监听多个日志文件
func NewTailTask(path,topic string) (tailObj *Tailtask){
	ctx,cancel:=context.WithCancel(context.Background())
	tailObj=&Tailtask{
		Path: path,
		Topic: topic,
		ctx: ctx,
		cancelFunc: cancel,
	}
	tailObj.init()		//根据路径去打开对应的日志
	return
}

//日志收集初始化
func (t *Tailtask) init(){
	config:=tail.Config{
		ReOpen: true, //重新打开
		Follow: true,//是否跟随
		Location: &tail.SeekInfo{Offset: 0,Whence: 2},//从文件的哪个地方开始读（当文件关闭后重新打开从上次位置开始读）
		MustExist: false,	//如果文件不存在并不报错，它会等待fileName的出现，如果出现了就按照ReOpen的值来重新打开文件开始读取内容
		Poll: true,		//poll的作用是把当前的文件指针挂到等待队列，针对linux
	}
	var err error
	t.Instance,err=tail.TailFile(t.Path,config)	//以config配置打开fileName
	if err!=nil{
		fmt.Println("tail file failed,err",err)
		return
	}
	go t.RunSendMsg()	//直接去采集日志发送到kafka
}

//往kafka中发送数据
func (t *Tailtask)RunSendMsg(){
	//1.读取日志
	for{
		select {
		//从ReadLog返回的管道中获取每行的内容，如果没获取到就休息一秒
			case line := <-t.Instance.Lines:
			//2.发送到Kafka
			//kafka.SendMsg(t.Topic,line.Text)  //这是函数调函数，即有一条日志则调用一次，调用完成再循环进行下一次，可以使用通道优化成异步
			//先把日志数据发到一个通道中，这边只管发，取让kafka包自己取,这个函数调用因为只是往通道中发送数据，因此速度会比之前快很多
			kafka.SendToChan(t.Topic,line.Text)
			case <-t.ctx.Done():
				return
		default:
			time.Sleep(time.Second)
		}
	}
}
