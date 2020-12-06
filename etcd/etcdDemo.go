package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)
func main(){
	//etcd连接
	cli,err:=clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},	//连接对象
		DialTimeout: 5*time.Second,		//5秒没连接上就超时返回err
	})
	if err!=nil{
		fmt.Println("connect to etcd faild,err:",err)
		return
	}
	defer cli.Close()

	//etcd Put
	ctx,cancel:=context.WithTimeout(context.Background(),time.Second)
	value:=`[{"topic":"web_log","path":"d:/xxx.log"},{"topic":"redis_log","path":"d:/redis.log"},{"topic":"nginx_log","path":"d:/nginx.log"}]`
	_,err=cli.Put(ctx,"fileName",value)
	//这里没开携程，因此会put执行完才执行cancel()，put如果一秒后还没有成功就会自动结束上下文，cancel()是避免等待直到垃圾回收结束它。
	cancel()
	if err!=nil{
		fmt.Println("put to etcd faild,err:",err)
		return
	}

	////etcd Get
	//ctx,cancel=context.WithTimeout(context.Background(),time.Second)
	////Get(ctx context.Context, key string, opts ...OpOption)中的第三个选项是clientv3.WithPrefix()，它可以给key定义一个前缀，方便整理
	//resp,err:=cli.Get(ctx,"xiaonuo")
	//cancel()
	//if err!=nil{
	//	fmt.Println("get to etcd faild,err:",err)
	//	return
	//}
	//for _,ev:=range resp.Kvs{
	//	fmt.Println("key:",string(ev.Key),"----value:",string(ev.Value))
	//}
	//
	////etcd Watch，watch会派一个哨兵一直监听着key的value的变化(新增、修改、删除)
	////watch返回的是一个channel,<-chan WatchResponse
	////也可以设置WithTimeout()来让它监视一个固定时间或设置条件让其取消监视（不过还是建议让其一直监视着）
	//rch:=cli.Watch(context.Background(),"xiaonuo")
	////从通道尝试取值，range会在没值的时候阻塞，有值后执行内部代码，然后循环回来继续阻塞等待下一个值
	//for wresp:=range rch{
	//	for _,ev:=range wresp.Events{
	//		//type可以取到类型，是删除修改还是新增
	//		fmt.Println("type:",string(ev.Type),"----key:",string(ev.Kv.Key),"----value:",string(ev.Kv.Value))
	//	}
	//}
}
