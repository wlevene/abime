package main

import (
	"os/exec"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"runtime"
	"net"
	"net/rpc"
	"fmt"
	"strings"
	
	"time"
    "math/rand"

	"gateway"
	"github.com/donnie4w/go-logger/logger"
)

/*
	gateway listen 端口6848

	6848 此端口用于rpc通信  
	由abime主动请求过来，验证app是否合法、是否需要上报Metric数据等
	
	目前只做简单的规则，且规则存放到内存中，后期再将规则外置到外部介质中
*/

var apppool *gateway.AppPool

type Gateway int

func (t *Gateway) ConnectGatewayAndVerifyApp(appKey string, reply *bool) error {
	//	*reply = "Gateway:" + args
	
	if appKey == "" {
		logger.Warn("appkey is nil")
		*reply = false
		return nil
	}

	appConfig := apppool.GetAppConfig(appKey)
	
	logger.Info("appConfig VailedApp", appConfig.VailedApp)
	
	*reply = appConfig.VailedApp
  
	return nil
}

func (t *Gateway) VerifyMetric(appkey string, reply *bool) error {

	ret := apppool.VailedMetric(appkey)
	logger.Info("VerifyMetric", appkey, ret)
	*reply = ret
	return nil
}

func init() {
	logger.SetConsole(false)
	logger.SetRollingDaily("./log", "gateway.log")
	logger.SetLevel(logger.INFO)
}


func post() {
	
	var count uint64
	
	fmt.Println(time.Now().Format("2006-01-02 15:04:05")," -- ", count)

	for {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))	
		i := r.Intn(100)
	
		str := fmt.Sprintf("curl -i -XPOST 'http://10.20.77.104:8086/write?db=testdb' --data-binary 'table,host=server01,region=us-west value=%d'", i)
	
		cmd := exec.Command("/bin/sh", "-c", str)
    	_, err := cmd.Output()
	 	if err != nil {
        	fmt.Println(err.Error())
    	}
		
		count++
		if count % 10000 == 0 {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), " -- ",count)
		}
		
	}
		
	return
	client := &http.Client{}

	reqest, err := http.NewRequest("POST", "http://10.20.77.104:8086/write?db=mydb", 
				strings.NewReader("test:host=server02 value=0.67"))
				
	//	reqest, err := http.NewRequest("GET", "http://www.baidu.com", nil)

	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		return
	}
	

//	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Linux; U; Android 4.4.4; zh-cn; YQ601 Build/KTU84P) AppleWebKit/533.1 (KHTML, like Gecko)Version/4.0 MQQBrowser/5.4 TBS/025478 Mobile Safari/533.1 MicroMessenger/6.3.7.52_rbb7fa12.660 NetType/WIFI Language/zh_CN")
//	reqest.Header.Add("Cookie", "CNZZDATA1252947662=1547195292-1447664434-%7C1447664434; PHPSESSID=oa8mu1kh598n95f7dlcj9hvca1 ")

//	reqest.Header.Add("Referer", "http://www.iweizhuli.com/index.php?&g=Wap&m=VbzhuliNew&a=zhuli&token=agkiyl1447206122&wecha_id=o6AXJjmzaK8uhcPMUd9rNTzR61h1&vbid=237&vbmid=31335&from=timeline&isappinstalled=0")
//	reqest.Header.Add("Host", "www.iweizhuli.com")

//	reqest.Header.Add("Accept-Charset", "utf-8, iso-8859-1, utf-16, *;q=0.7")
//	reqest.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

//	reqest.Header.Add("Accept-Language", "zh-CN")
//	reqest.Header.Add("Origin", "http://www.iweizhuli.com")

//	reqest.Header.Add("Accept", "*/*")
//	reqest.Header.Add("Content-Length", "13")

//	reqest.Header.Add("X-Requested-With", "XMLHttpRequest")
//	reqest.Header.Add("Accept-Encoding", "gzip")

//	reqest.Header.Add("Connection", "")

	response, err := client.Do(reqest)

	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		return
	}

	defer response.Body.Close()

	if response.StatusCode == 200 {
		fmt.Println("Forward success")
	}
	
	fmt.Println(response.StatusCode)
	result, err := ioutil.ReadAll(response.Body) 
  	response.Body.Close() 
  	if err != nil { 
        return 
  	} 
  	fmt.Printf("%s", result)
}

func main() {


	go post()
	go post()
	go post()
	go post()
	
	post()

	return;
	runtime.GOMAXPROCS(runtime.NumCPU())

	logger.Info("Gateway init...")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

	apppool = gateway.NewAppPool()

	go func() {
		rpc.Register(new(Gateway))
		rpc.HandleHTTP()
		l, e := net.Listen("tcp", ":6848")
		if e != nil {
			logger.Fatal("Gateway listen error:", e)
		}
		http.Serve(l, nil)
	}()
	
	// wait for signal...
	<-signalChan	
}

