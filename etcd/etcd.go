package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var cli *clientv3.Client

type logEntry struct {
	Topic string `json:"topic"`
	Path string `json:"path"`

}

func Init(address string,timeout time.Duration)(err error){
	//etcd连接
	cli,err=clientv3.New(clientv3.Config{
		Endpoints: []string{address},	//连接对象
		DialTimeout: timeout,		//设置超时时间，没连接上就超时返回err
	})
	if err!=nil{
		fmt.Println("connect to etcd faild,err:",err)
		return
	}
	return
}

//从etcd中根据key获取配置信息
func GetConf(key string)(logEntryConf []*logEntry,err error){
	ctx,cancel:=context.WithTimeout(context.Background(),time.Second)
	resp,err:=cli.Get(ctx,key)
	cancel()
	if err!=nil{
		fmt.Println("get to etcd faild,err:",err)
		return
	}
	for _,ev:=range resp.Kvs{
		//通过json的反序列化将etcd中的值对应填充到logEntryConf,因此etcd的value得是json格式的字符串（因为value需要存储topic和fileName，所以需要使用反序列化）
		//json传输如果用cmd put的话会过滤掉双引号，因此需要用个go程序调用clientv3的put传输
		err=json.Unmarshal(ev.Value,&logEntryConf)
		if err!=nil{
			fmt.Println("Unmarshal etcd value faild,err:",err)
			return
		}
	}
	return
}