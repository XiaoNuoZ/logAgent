package conf

//kafka配置
type KafkaConf struct {
	Address string		`ini:"address"`
	Topic string		`ini:"topic"`
}

//taillog配置
type TaillogConf struct {
	Filename string		`ini:"filename"`
}

type LogAgentConf struct {
	KafkaConf		`ini:"kafka"`
	TaillogConf		`ini:"taillog"`
}