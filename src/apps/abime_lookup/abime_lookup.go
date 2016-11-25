package main

import (
	"strings"
	"encoding/json"
	//	"encoding/json"
	"fmt"
	"net/http"
//	"runtime"
//	"strings"
	"time"
	"sync"
//	"math/rand"

	//	"protocol"

	"abime_lookup"
	"config"

	"github.com/donnie4w/go-logger/logger"
)

/*
	lookup listen 端口6840 6841

	6840客户端连接 并分配abime给client, 使用短连接(http)
	各个abime会通过此端口向lookup注册、反注册、上报自身的性能情况 使用短连接(http)

	6841是id ejector使用的rpc进程间通信
*/


const (
	ClientMaxCountAtAbime = 200000
)

var abimePool *abime_lookup.AbimePool
var dataGateway * abime_lookup.DataGateway

//var id_ejector_pool *abime_lookup.IdEjectorPool

//type IDEjector int

//func (t *IDEjector) ConnectIDEjector(appkey_startID string, reply *bool) error {
//	//	*reply = "IDEjector:" + args
//	strArray := strings.Split(appkey_startID, ":")
//	if len(strArray) != 2 {
//		*reply = false
//		return nil
//	}

//	*reply = id_ejector_pool.CreateIdEjector(strArray[0], strArray[1])

//	return nil
//}

//func (t *IDEjector) NewID(appkey string, reply *string) error {

//	*reply = id_ejector_pool.NewID(appkey)
//	return nil
//}

func init() {
	logger.SetConsole(true)
	logger.SetRollingDaily("./log", "abime_lookup.log")
	logger.SetLevel(logger.DEBUG)
}

func main() {
	// go 1.5之后不需要此设置
	// runtime.GOMAXPROCS(runtime.NumCPU())

	logger.Info("AbimeLookup begin init...")
//	abime_lookup.StartIdEjectorServer()

	go checkAbimeStatus()

	abimePool = abime_lookup.NewAbimePool()
	dataGateway = abime_lookup.NewDataGateway()
	

	http.HandleFunc("/registeabime", registerabime)
	http.HandleFunc("/reportbime", reportbime)

	// 由客户端请求此api, 返回可用的abime server给客户端
	http.HandleFunc("/shake", shake)
	http.HandleFunc("/status", status)
		
	// supper both http and https in a single instance
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err := http.ListenAndServe(config.LookUpHTTPADDR, nil); err != nil {
			logger.Fatal("ListenAndServe", config.LookUpHTTPADDR, err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		if err := http.ListenAndServeTLS(config.LookUpHTTPADDR_HTTPS, "cert.pem", "key.pem", nil); err != nil {
			logger.Fatal("ListenAndServeTLS", config.LookUpHTTPADDR_HTTPS, err)
		}
	}()
	
	wg.Wait()
}


// 每3s 将不合法或者过期的abime server 从列表中去掉
func checkAbimeStatus() {
	for _ = range time.Tick(3 * time.Second) {

		count := abimePool.Count()
		if count <= 0 {
			continue
		}

		logger.Debug("abimePool:", abimePool.Abimes())

		abimes := abimePool.Abimes()
		now_time := time.Now().Unix()
		
		/*
			不合法的abime定义为:
			1：10分析未上报数据的abime
			2：ip是特殊ip地址的abime
			3：ip地址异常的abime
		*/

		for _, ab := range abimes {
			if ab.ReportTime < now_time-10*60 {
				logger.Info("abimePool UnRegisterAbime by time", ab)
				abimePool.UnRegisterAbime(ab)
			}
			
			if strings.Contains(ab.Addr, "0.0.0.0") ||
				strings.Contains(ab.Addr, "127.0.0.1") {
				logger.Info("abimePool UnRegisterAbime by addr error", ab.Addr)
				abimePool.UnRegisterAbime(ab)
			}

			if len(ab.Addr) < 7 {
				logger.Info("abimePool UnRegisterAbime by data invaild", ab)
				abimePool.UnRegisterAbime(ab)
			}
		}
	}
}

var shake_count int64
func shake(w http.ResponseWriter, r *http.Request) {

	/*
		两种策略
			1 将连接数据平均的分配的各个abme 以保证每个abime的性能开销都最小
			2 将一个abime填满后，才将新的连接分配到下一个 abime, 让abime都集中在最少量的
				abime中减少维护的工作量
			3 如果每个 abime 都满的时候 可以考虑分配给最小的那个abime
				
		目前使用的是第二条策略
		
		测试的结果是每个server 15w 应该是没问题的 到18w多就开始会出现超时
		
		修改了一版本后 每个server测试到20w完全ok，上限因为测试机器不足，没有测试出上限，预估
		可以到40w
	*/ 
	
//	abime := abimePool.LessHeavyestAbime()

	params := r.URL.Query()
	
	var appkey string
	
	if len(params["appkey"]) > 0 {
		appkey = params["appkey"][0]
	}
	
	var device_id string
	if len(params["device"]) > 0 {
		device_id = params["device"][0]
	}
		
	fmt.Println("appkey :", appkey)
	fmt.Println("device_id :", device_id)

	abime := abimePool.LessHeavyestAbimeAboveAt(ClientMaxCountAtAbime)
	
	retMap := make(map[string]interface{})
	retMap["rand"] = 100	
//	if appkey == "6ac9e810" {
//		retMap["rand"] = 10
//	}
	retMap["svr"] = abime.Addr	
	
	/*
		START: 将部分数据引流到新的系统上  start ==================== 
	*/
//	if appkey == "6ac9e810" {
		
//		if shake_count % 85 != 0 {
//			retMap["svr"] = "121.14.30.235:6842"
//		} else {
//			retMap["svr"] = "121.14.30.242:6842"
//		}
		
//		shake_count++
		
//		if shake_count >= 10000000 {
//			shake_count = 0
//		}
//	}
	
	/*
		END: 将部分数据引流到新的系统上   end   ==================== 
	*/
	
	retMap["battery"] = false	
	
	
	/* fps 返回结构 {"fps":true,"normal_fps":true, "calamity_fps":12}
		fps	总开关
		normal_fps normal的fps是否上报，规则可同电量一样每天的前200台
		calamity_fps 定义为calamity的fps一般指非常低的fps, 可和crash同样理解
	*/
	
	// fps 总开关一般都为true
	retMap["fps"] = true
	retMap["normal_fps"] = false
	retMap["calamity_fps"] = 12
	
	
	if abime.Addr == "" {
		retMap["ret"] = -1
	} else {
		retMap["ret"] = 0
		
		// 查询此appkey已收集到多少battery信息
		if device_id != "" {			
			if dataGateway.IsPassBatteryGateway(appkey) {
				retMap["battery"] = true
			}			
			dataGateway.ReportDeviceToBatteryGateway(appkey, device_id)
			
			if dataGateway.IsPassFpsGateway(appkey) {
				retMap["normal_fps"] = true
			}			
			dataGateway.ReportDeviceToFpsGateway(appkey, device_id)
		}	
	}
		
	
	// 大于 (比如:20w) 服务器繁忙, 暂不接受新的连接
	if abime.ClientCount >= ClientMaxCountAtAbime {
		logger.Warn("abime ", abime.AbimeName, abime.Addr, abime.ClientCount, " is full")
		retMap["ret"] = -1
		retMap["rand"] = 0
		retMap["svr"] = ""
		retMap["battery"] = false
		retMap["fps"] = false
		retMap["normal_fps"] = false
		retMap["calamity_fps"] = -1		
	}
	
	if b, err := json.Marshal(retMap); err == nil {
		var ret string
		ret = string(b)
		
		logger.Debug(ret)
						
		w.Write([]byte(ret))
		return
	}

	w.Write([]byte("has error"))
}

func status(w http.ResponseWriter, r *http.Request) {

	abimes := abimePool.Abimes()
	
	var msgCount uint64
	var clientCount int64
	for _, abime := range abimes {
		msgCount += abime.MessageCount
		clientCount += abime.ClientCount
	}
	
	retMap := make(map[string]interface{})
	
	retMap["abimes"] = abimes
	retMap["MsgCount"] = msgCount
	retMap["Clientcount"] = clientCount
	retMap["ClientMaxCount_AtAbime"] = ClientMaxCountAtAbime

	if b, err := json.Marshal(retMap); err == nil {
		var ret string
		ret = string(b)
		w.Write([]byte(ret))
		return
	}

	w.Write([]byte("has error"))
}

func reportbime(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	buf := make([]byte, 1024)
	l, err := r.Body.Read(buf)
	bufx := buf[:l]

	logger.Debug("reportbime abime:", string(bufx))

	var abime abime_lookup.Abime
	err = json.Unmarshal(bufx, &abime)
	if err == nil {
		abime.ReportTime = time.Now().Unix()
		_, exist := abimePool.RegisterAbime(abime)
		if exist == false {
			logger.Debug("Register abime", abime)
		} else {
			logger.Debug("Updata abime", abime)
		}
		abimePool.UpdateAbime(abime)
		w.Write([]byte("ok"))
		return
	} else {
		logger.Fatal(err, string(bufx))
	}
	w.Write([]byte("error"))
}

func registerabime(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	buf := make([]byte, 1024)
	r.Body.Read(buf)

	logger.Debug("abime:", string(buf))

	var abime abime_lookup.Abime
	if err := json.Unmarshal(buf, &abime); err == nil {
		abimePool.RegisterAbime(abime)
		w.Write([]byte("ok"))
		return
	}
	w.Write([]byte("error"))
}





