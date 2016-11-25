package abime_lookup

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/donnie4w/go-logger/logger"
)

type IdEjector interface {
	NewID() string
}

type RedisIdEjector struct {
	IdEjector
	Appkey    string
	idChannel chan string
}

func NewRedisIDEjector(appkey string, startID string) *RedisIdEjector {

	p := &RedisIdEjector{}

	p.Appkey = appkey
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

		logger.Info("START ID:", counter, " appkey :", appkey)
		for {
			p.idChannel <- fmt.Sprintf("%s_%x", appkey, counter)
			counter += 1
			
			// reset counter to 1  every 50000000 record
			if counter >= 50000000 {
				counter = 1;
			}			
		}
	}(p.Appkey)

	return p
}

func (p *RedisIdEjector) NewID() string {
	return <-p.idChannel
}
