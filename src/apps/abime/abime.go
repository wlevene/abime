package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	//	"protocol"
	"server"

	"github.com/donnie4w/go-logger/logger"
)

/*

[--------------------10byte	----------------]
|4(body size)|2(proto version)|4(proto id)  |
[					header					] [body...]

*/

/*

./abime_liunx_64 -tcpaddr=:6842 -httpaddr=:6843 -redis=redis://172.16.8.209:6379/4?pwd= -lookupaddr=:6840 -profile=:9000 -name=supperman

./abime_liunx_64 -tcpaddr=121.14.30.235:6842 -redisaddr=127.0.0.1:6379 -redispwd=no -lookupaddr=0.0.0.0:6840 -name=ghost_1

go run abime.go -tcpaddr=127.0.0.1:6842 -redisaddr=121.14.28.25:6379 -redispwd=no -lookupaddr=10.29.70.42:6840 -name=dage
go run abime.go -tcpaddr=127.0.0.1:6842 -redisaddr=121.14.30.235:6379 -redispwd=no -lookupaddr=121.14.30.235:6840 -name=dage1


*/

var (
	flagSet    	= flag.NewFlagSet("abime", flag.ExitOnError)
	tcpaddr    	= flagSet.String("tcpaddr", ":6842", "tcpaddr(<addr>:<port>) to listen on for TCP clients")
	httpaddr   	= flagSet.String("httpaddr", ":6843", "httpaddr(<addr>:<port>) to get abime status for HTTP ")
	redis  		= flagSet.String("redis", ":6379", "redisaddr(<addr>:<port>/<number>?pwd=pwd) to connect to redis")
//	redispwd   	= flagSet.String("redispwd", "pwd", "password for redis auth. if redis has no password, put \"no\" value")
	lookupaddr 	= flagSet.String("lookupaddr", ":6840", "lookup addr(<addr>:<port>) ")
	abimename  	= flagSet.String("name", "", "abime server name")
	profile    	= flagSet.String("profile", ":9000", "golang program profile(<addr>:<port>) ")
)

func init() {
	
	logger.SetConsole(false)
	logger.SetLevel(logger.INFO)
	
}

func main() {
	
	defer func() {
        if err := recover(); err != nil {
            logger.Fatal("============ EXIT =========== \n", err)
        }
    }()
	
	flagSet.Parse(os.Args[1:])

	logFileName := "abime_default.log"
	if (*abimename) != "" {
		logFileName = fmt.Sprintf("abime_%s.log", (*abimename))
	}

	logger.SetRollingDaily("./log", logFileName)
//	logger.SetRollingFile("./log", logFileName, 1000, 500, logger.MB)

	logger.Info("Abime begin init...")

	logger.Info("Abime Name:", (*abimename), 
		" tcpaddr:", (*tcpaddr), 
		" httpaddr:", (*httpaddr),
		" redis:", (*redis), 
		" lookupaddr:", (*lookupaddr))

	svrOptions := server.NewAbimeOption()
	
	if (*tcpaddr) == "" {
		logger.Fatal("tcpaddr param is nil")
		return
	}
	svrOptions.TcpAddr = *tcpaddr
	
	if (*httpaddr) == "" {
		logger.Fatal("httpaddr param is nil")
		return
	}
	svrOptions.HttpAddr = *httpaddr

	if (*redis) == "" {
		logger.Fatal("redis param is nil")
		return
	}
	svrOptions.RedisAddr = *redis
	
	if (*lookupaddr) == "" {
		logger.Fatal("lookupaddr param is nil")
		return
	}
	svrOptions.LookupAddr = *lookupaddr

	if (*abimename) == "" {
		logger.Fatal("abimename param is nil")
		return
	}
	svrOptions.AbimeName = *abimename

	if (*profile) != "" {
		svrOptions.Profile = *profile
	}

	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
    	<-signalChan
    	exitChan <- 1
	}()


	logger.Info("Abime:", (*abimename), " Begin run loop ...")
	svr := server.NewServer(svrOptions)


	signal.Notify(signalChan, 
		syscall.SIGINT, 
		syscall.SIGTERM, 
		syscall.SIGABRT, 
		syscall.SIGKILL,
		syscall.SIGQUIT)
	
	svr.Main()
	<-exitChan
	svr.Exit()
	logger.Info("Shutdown ok")
}
