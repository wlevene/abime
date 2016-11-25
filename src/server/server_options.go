package server

type abimeOption struct {
	TcpAddr       string `flag:"tcpaddr"`
	HttpAddr	  string `flag:"httpaddr"`
	RedisAddr     string `flag:"redisaddr"`
	RedisPassword string `flag:"redispwd"`
	LookupAddr    string `flag:"lookupaddr"` // 127.0.0.1:6840 交互压力不大 外网地址也影响不大
	IDSeverAddr   string `flag:"idsvraddr"`  // 127.0.0.1:6841 尽量是内网地址
	AbimeName     string `flag:"name"`
	GatewayAddr   string `flag:"gateway"` // 127.0.0.1:6848 尽量是内网地址
	Profile		  string `flag:"gateway"`
}

func NewAbimeOption() *abimeOption {
	p := &abimeOption{
		TcpAddr:       ":6842",
		HttpAddr:	   ":6843",
		RedisAddr:     ":6379",
		RedisPassword: "probeKing+5688",
		GatewayAddr:   ":6848",
		Profile:	   ":9000",
	}
	return p
}
