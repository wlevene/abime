package archivedata

import (
	"fmt"
	"strings"
)
import "github.com/garyburd/redigo/redis"
import "github.com/donnie4w/go-logger/logger"
//import "time"
import "net/url"

const (
	REDIS_KEY_QUEUE =  1000000
	MAX_PUSH_DATA_COUNT = 1000000
	MAX_REDIS_CONN_NUM = 200
)

type ArchiveRedis struct {
	
	IArchive
	
	keyProviderChan chan string
	pushDataChan	chan string
	opts            *RedisOption
	
//	redisPool *redis.Pool
}

func NewArchiveRidesEx(options *RedisOption) *ArchiveRedis {
	
	p := &ArchiveRedis{}

	p.opts = options
	
	p.keyProviderChan = make(chan string, REDIS_KEY_QUEUE)
	p.pushDataChan =  make(chan string, MAX_PUSH_DATA_COUNT)
	
	p.makekeys()
	
//	p.redisPool = p.newPool(p.opts.RedisAddr, p.opts.RedisPassword)
	
	p.archiveData()
		
	return p
} 

func (p * ArchiveRedis ) makekeys() {
	go func() {
		/*
			每 REDIS_KEY_QUEUE 多产生一个新的随机数种子，然后根据此随机数种子
			产生REDIS_KEY_QUEUE多条id, 再重新产生一个新的随机数种子
		*/
		u := fmt.Sprintf("%s", RandomAlphabetic(10))
		var counter uint64 = 1
		
		for {
			p.keyProviderChan <- fmt.Sprintf("%s_%x", u, counter)
			counter += 1
			
			if (counter > REDIS_KEY_QUEUE) {
				u = fmt.Sprintf("%s", RandomAlphabetic(8))
				counter = 1
			}
		}
	}()
}

func (p * ArchiveRedis ) refreshANewKey() string {
	return <- p.keyProviderChan
}

func (p * ArchiveRedis ) archiveData() {
	
	var watchDog4RedisConnectionChan chan int
	watchDog4RedisConnectionChan =  make(chan int, MAX_REDIS_CONN_NUM)
	
	// run watch dog for redis connection
	go func() {
		for _ = range watchDog4RedisConnectionChan {
			logger.Debug("Watch dog run redis connection")
			p.connectRedisAndWork(watchDog4RedisConnectionChan)
		}
	}()
	
	for i := 0; i < MAX_REDIS_CONN_NUM; i++ {
		p.connectRedisAndWork(watchDog4RedisConnectionChan)
	}
}

func (p * ArchiveRedis ) connectRedisAndWork(watchDog4RedisConnectionChan chan int) {
	// connect redis
	
	redis_url, err := url.Parse(p.opts.RedisAddr)
	
	if err != nil {
		logger.Fatal(err)
		return
	}
	
//	fmt.Println("host:", redis_url.Host)
//	fmt.Println("path:", redis_url.Path)
//	fmt.Println("Scheme:", redis_url.Scheme)
//	fmt.Println("RawQuery:", redis_url.RawQuery)
	
	redis_pwd := redis_url.Query().Get("pwd")
	
	c, err := redis.Dial("tcp",redis_url.Host)
	if err != nil {
		logger.Fatal(err)
		watchDog4RedisConnectionChan <- 1
	    return
	}
	
//	defer c.Close()

	if redis_pwd != "" {
		if _, err := c.Do("AUTH", redis_pwd); err != nil {
			c.Close()
			logger.Fatal(err)
	        return
	    }
	}
	
	// select table
	if _, err := c.Do("SELECT", strings.Trim(redis_url.Path,"/")); err != nil {
		logger.Fatal(err)
		c.Close()
		watchDog4RedisConnectionChan <- 1
	    return
	}
			
	go func(conn redis.Conn) {
		
		ok := false
		var err error
		
		for pushData := range p.pushDataChan {	
					
			key := p.refreshANewKey()

			if key == "" {
				continue
			}
			
			if ok, err = redis.Bool(c.Do("HSET", "hash", key, pushData)); ok {
				c.Do("LPUSH", "keys", key)
			} else {
				c.Close()
				logger.Fatal(err)
				watchDog4RedisConnectionChan <- 1
				return
			}
		}
	}(c)	
}


func (p * ArchiveRedis ) PushData(data string) {
	
	if data == "" {
		return
	}
	
	p.pushDataChan <- data
}

func (p * ArchiveRedis ) Exit() {

}
