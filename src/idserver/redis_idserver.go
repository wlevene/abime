package idserver

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/donnie4w/go-logger/logger"
)

type RedisIDServer struct {
	appkey    string
	idChannel chan string
}

func NewRedisIDServer(appkey string, startID string) *RedisIDServer {

	if appkey == "" {
		return nil
	}

	p := &RedisIDServer{
		appkey: appkey,
	}

	p.idChannel = make(chan string, 1)

	go func(appkey string) {
		var counter uint64 = 1
		if startID != "" {
			t := strings.Split(startID, "_")
			if len(t) == 2 {
				lastIDString := t[1]
				lastID, err := strconv.ParseUint(lastIDString, 16, 64)
				if err == nil && lastID > 1 {
					counter = lastID
				}
			}
		}
		logger.Info("====APPKEY:", p.appkey, "START ID:", counter)
		for {
			p.idChannel <- fmt.Sprintf("%s_%x", appkey, counter)
			counter += 1
		}
	}(p.appkey)

	return p
}

func (p *RedisIDServer) NewID() string {
	return <-p.idChannel
}
