//从各种服务中读取最新日志，引用第三方包tailf
package tail_log

import (
	"fmt"
	"github.com/hpcloud/tail"
	_ "github.com/hpcloud/tail"
)
var (
	tailObj *tail.Tail
)

//日志收集初始化
func Init(fileName string) (err error){
	config:=tail.Config{
		ReOpen: true, //重新打开
		Follow: true,//是否跟随
		Location: &tail.SeekInfo{Offset: 0,Whence: 2},//从文件的哪个地方开始读（当文件关闭后重新打开从上次位置开始读）
		MustExist: false,	//如果文件不存在并不报错，它会等待fileName的出现，如果出现了就按照ReOpen的值来重新打开文件开始读取内容
		Poll: true,		//poll的作用是把当前的文件指针挂到等待队列，针对linux
	}
	tailObj,err=tail.TailFile(fileName,config)	//以config配置打开fileName
	if err!=nil{
		fmt.Println("tail file failed,err",err)
		return err
	}
	return

}

//从文件中读取日志
func ReadLog() <-chan *tail.Line{
	return tailObj.Lines
}
