package conf

//kafka配置
type KafkaConf struct {
	Address string		`ini:"address"`
}

type Etcd struct {
	Address string		`ini:"address"`
	Timeout int 		`ini:"timeout"`
}

type LogAgentConf struct {
	KafkaConf		`ini:"kafka"`
	Etcd			`ini:"etcd"`
}