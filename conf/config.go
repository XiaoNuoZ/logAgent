package conf

import (
	"net"
	"strings"
)

//kafka配置
type KafkaConf struct {
	Address string		`ini:"address"`
	MaxSize int 		`ini:"chan_max_size"`
}

type Etcd struct {
	Address string		`ini:"address"`
	Timeout int 		`ini:"timeout"`
	FilePathKey string 	`ini:"colletc_log_key"`
}

type LogAgentConf struct {
	KafkaConf		`ini:"kafka"`
	Etcd			`ini:"etcd"`
}

// 获取本地对外IP
func GetOutboundIP()(ip string,err error){
	//连接google获取连接
	conn,err:=net.Dial("udp","8.8.8.8:80")
	if err != nil {
		return
	}
	defer conn.Close()
	//获取ip,以：分割的第0个
	localAddr:=conn.LocalAddr().(*net.UDPAddr)
	//fmt.Println(localAddr.String())
	ip=strings.Split(localAddr.IP.String(),":")[0]
	return
}