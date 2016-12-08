package main

import (
	"strings"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"math/rand"
	//	"protocol"
	"server"
	

	"github.com/donnie4w/go-logger/logger"
)

import "net/url"
/*

|						10byte				|
|4(body size)|2(proto version)|4(proto id)  |
[					header					] [body...]

*/

/*
go run main.go -tcpaddr=10.29.70.42:6842 -redisaddr=121.14.28.25:6379 -redispwd=probeKing+5688 -lookupaddr=10.29.70.42:6840 -idsvraddr=121.14.28.25:6841 -name=dage
*/

var (
	flagSet    = flag.NewFlagSet("abime", flag.ExitOnError)
	tcpaddr    = flagSet.String("tcpaddr", ":6842", "tcpaddr(<addr>:<port>) to listen on for TCP clients")
	redisaddr  = flagSet.String("redisaddr", ":6379", "redisaddr(redis://<addr>:<port>) to connect to redis")
	redispwd   = flagSet.String("redispwd", "pwd", "password for redis auth")
	lookupaddr = flagSet.String("lookupaddr", ":6840", "lookup addr(<addr>:<port>) ")
	idsvraddr  = flagSet.String("idsvraddr", ":6841", "idsvr addr(<addr>:<port>) ")
	abimename  = flagSet.String("name", "", "abime server name")
)

func init() {
	logger.SetConsole(true)

	logger.SetLevel(logger.INFO)
}

func main() {
	
	flagSet.Parse(os.Args[1:])

	logFileName := "abime_default.log"
	if (*abimename) != "" {
		logFileName = fmt.Sprintf("abime_%s.log", (*abimename))
	}

	logger.SetRollingDaily("./log", logFileName)

	logger.Info("Abime begin init...")

	logger.Info("Abime Name:", (*abimename), " tcpaddr:", (*tcpaddr), " redisaddr:", (*redisaddr), " lookupaddr:", (*lookupaddr), " idServerraddr:", (*idsvraddr))

	//	fmt.Println("==tcpaddr:", *tcpaddr)
	//	fmt.Println("==redisaddr:", *redisaddr)

	svrOptions := server.NewAbimeOption()
	if (*tcpaddr) != "" {
		svrOptions.TcpAddr = *tcpaddr
	}

	if (*redisaddr) != "" {
		svrOptions.RedisAddr = *redisaddr
	}

	if (*redispwd) != "" {
		svrOptions.RedisPassword = *redispwd
	}

	if (*lookupaddr) != "" {
		svrOptions.LookupAddr = *lookupaddr
	}

	if (*abimename) != "" {
		svrOptions.AbimeName = *abimename
	}

	if (*idsvraddr) != "" {
		svrOptions.IDSeverAddr = *idsvraddr
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

	logger.Info("Abime:", (*abimename), " Begin run loop ...")
	svr := server.NewServer(svrOptions)

	svr.Main()
	<-signalChan
	svr.Exit()
}
