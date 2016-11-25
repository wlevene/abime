package protocol

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net"
	"strings"
	"time"
	"os"

	"archivedata"
	
	"github.com/pquerna/ffjson/ffjson"
)

import "github.com/donnie4w/go-logger/logger"

type ClientV1 struct {
	net.Conn
	session SessionMessage
	network AppNetWorkStatus

	Appkey string
	ClientStartTime time.Time
	Arch        *archivedata.ArchiveRedis
	
	deviceInfoMap map[string]interface{}
	appInfoMap map[string]interface{}
}

func NewClientV1(conn net.Conn, arch *archivedata.ArchiveRedis ) *ClientV1 {

	return &ClientV1{
		Conn: conn,
		Arch : arch,
	}
}

func (c *ClientV1) Exit() {

}

func (c *ClientV1) String() string {
	addr := c.RemoteAddr().String()
	return strings.Split(addr, ":")[0]
}

func (c *ClientV1) AllString() string {
	return c.RemoteAddr().String()
}

func (c *ClientV1) appInfo2Json() (map[string]interface{}, error) {
	
	if c.appInfoMap != nil &&
	 len(c.appInfoMap) != 0 {
		return c.appInfoMap, nil
	}
		
	var result map[string]interface{}
	var err error

	if c.session.App != nil {
		data, err := json.Marshal(c.session.App)
		
		if err == nil {			
			if err = ffjson.Unmarshal(data, &result); err != nil {				
        		return result, err
    		}
			
			result["appkey"] = c.session.Appkey
		}
	}
	
	c.appInfoMap = result;
	
	return result, err
}

func (c *ClientV1) deviceInfo2Json() (map[string]interface{}, error) {
	
	if c.deviceInfoMap != nil &&
	 len(c.deviceInfoMap) != 0 {
		return c.deviceInfoMap, nil
	}
	
	var result map[string]interface{}
	var err error

	if c.session.Device != nil {
		
		data, err := json.Marshal(c.session.Device)
		
		if err == nil {
			if err = ffjson.Unmarshal(data, &result); err != nil {				
        		return result, err
    		}
		}
	}
	
	c.deviceInfoMap = result
	return result, err
}

func (c *ClientV1) archiveJsonData(dataMap map[string]interface{}, cmdType CmdType) {
	
	resultJson := make(map[string]interface{})

	appInfo, err := c.appInfo2Json()
	if err != nil {
		logger.Fatal(err)
		return
	}

	deviceInfo, err := c.deviceInfo2Json()

	if err != nil {
		logger.Fatal(err)
		return
	}
	
	resultJson["app"] = appInfo
	resultJson["device"] = deviceInfo
	resultJson["createtime"] = time.Now().UnixNano() / int64(time.Millisecond) // .UnixNano()
	resultJson["datatype"] = cmdType
	
	switch cmdType {
		
		case CmdType_WEBVIEW_DATA:
			resultJson["webview"] = dataMap
			
		case CmdType_HTTP_DATA:
			resultJson["httpdata"] = dataMap
			
		case CmdType_METHOD_DATA:
			resultJson["methoddata"] = dataMap
			
		case CmdType_BATTERY_DATA:
			resultJson["batterydata"] = dataMap
			
		case CmdType_FPS_DATA:
			resultJson["fpsdata"] = dataMap
			
		case CmdType_FPS_CALAMITY_DATA:
			resultJson["calamityfpsdata"] = dataMap
			
		default:
			return		
	}

	if c.network.NetworkName != nil {
		resultJson["network"] = (*c.network.NetworkName)
	} else {
		resultJson["network"] = ""
	}
	
	if c.network.CarrierName != nil {
		resultJson["carrier"] = (*c.network.CarrierName)
	} else {
		resultJson["carrier"] = ""
	}
	
	if 	c.String() != "" {
		resultJson["clientip"] = c.String()
	}
	
	if (*c.session.SessionId) != "" {
		resultJson["sessionid"] = (*c.session.SessionId)
	}
	
	if (*c.session.SdkVersion) != "" {
		resultJson["sdkversion"] = (*c.session.SdkVersion)
	}	
	
	data_str, err := ffjson.Marshal(resultJson)
	
	if err != nil {
		logger.Fatal(err)
		return
	}

	if c.Arch != nil {
		c.Arch.PushData(string(data_str))
	}
}

func (c *ClientV1) archiveHttpData(httpData *HttpData) {
	
	defer func() {
  	if r := recover(); r != nil {
		logger.Fatal(r)
  		}
 	}()
	
	if httpData == nil {
		return
	}

	var err error

	data, err := json.Marshal(httpData)
	if err != nil {
		logger.Fatal("json.Marshal(httpData)", err)
		return
	}

	var httpMap map[string]interface{}
    if err := ffjson.Unmarshal(data, &httpMap); err != nil {
		logger.Fatal("json.Unmarshal(httpData)", err)
        return
    }

	c.archiveJsonData(httpMap, CmdType_HTTP_DATA)
}

func (c *ClientV1) archiveWebViewData(webData *WebViewData) {
	
	defer func() {
  	if r := recover(); r != nil {
		logger.Fatal(r)
  		}
 	}()

	if webData == nil {
		return
	}
	var err error

	data, err := json.Marshal(webData)
	if err != nil {
		logger.Fatal("json.Marshal(webData)", err)
		return
	}

	var dataMap map[string]interface{}
    if err := ffjson.Unmarshal(data, &dataMap); err != nil {
		logger.Fatal("ffjson.UnmarshalFast(httpData)", err)
        return
    }
	
	c.archiveJsonData(dataMap, CmdType_WEBVIEW_DATA)
}

func (c *ClientV1) archiveMethodData(methodData *MethodData) {
	
	defer func() {
  	if r := recover(); r != nil {
		logger.Fatal(r)
  		}
 	}()
	
	if methodData == nil {
		return
	}

	var err error

	data, err := json.Marshal(methodData)
	if err != nil {
		logger.Fatal("json.Marshal(methodData)", err)
		return
	}

	var dataMap map[string]interface{}
    if err := ffjson.Unmarshal(data, &dataMap); err != nil {
		logger.Fatal("ffjson.UnmarshalFast(httpData)", err)
        return
    }
	
	c.archiveJsonData(dataMap, CmdType_METHOD_DATA)
}
  


func (c *ClientV1) archiveBatteryData(batteryData *BatteryData) {
	
	defer func() {
  	if r := recover(); r != nil {
		logger.Fatal(r)
  		}
 	}()

	if batteryData == nil {
		return
	}
	
	/*
		丢掉客户端会上报电量均为1的异常数据
	*/
	count := len(batteryData.Batterys)
	index := 1
	
	for _, item := range batteryData.Batterys {
		
		if (*item.CurrentBattery) != 1 {
			break;
		}
		
		if index == count {
//			logger.Info("Error Data And Skip it : ", (*c.session.Appkey))
			return;
		}
		index++;
	}
		
	
	var err error

	data, err := json.Marshal(batteryData)
	if err != nil {
		logger.Fatal("json.Marshal(batteryData)", err)
		return
	}

	var dataMap map[string]interface{}
    if err := ffjson.Unmarshal(data, &dataMap); err != nil {
		logger.Fatal("ffjson.UnmarshalFast(batteryData)", err)
        return
    }
	
	c.archiveJsonData(dataMap, CmdType_BATTERY_DATA)
}

func (c * ClientV1) archiveFpsData(fpsData *FpsData, cmdType CmdType) {
	
	defer func() {
  	if r := recover(); r != nil {
		logger.Fatal(r)
  		}
 	}()

	if fpsData == nil {
		return
	}
	var err error

	data, err := json.Marshal(fpsData)
	if err != nil {
		logger.Fatal("json.Marshal(fpsData)", err)
		return
	}

	var dataMap map[string]interface{}
    if err := ffjson.Unmarshal(data, &dataMap); err != nil {
		logger.Fatal("ffjson.UnmarshalFast(fpsData)", err)
        return
    }
	
	// save snapshot and clear it from map
	// save snapshot jpg file
	var imageByte []byte
	imageByte = fpsData.GetSnapShot()
//	logger.Info("image byte len:", len(imageByte))
//	logger.Info("appkey2:", (*c.session.Appkey))
	
	if len(imageByte) > 0 {
		
		if (*c.session.Appkey) != "" &&
		(* c.session.SessionId) != "" {
				
			path := (*c.session.Appkey) + "/" + (*c.session.SessionId) + "/";		
			savePath := "/letapm/upload/fps/" + path	
				
		 	_, err = os.Stat(savePath)
			exists := err == nil || os.IsExist(err)
			
			if !exists {
				logger.Info("mkdir", savePath)
				os.MkdirAll(savePath, 0777)
			}
					
			now := time.Now().UnixNano()
			fileFullPath := savePath + fmt.Sprintf("%d", now) + ".jpg"
	
 			err = ioutil.WriteFile(fileFullPath, imageByte, 0666)
			if err != nil {
				logger.Fatal("Save Fps Snapshot Appkey:", c.Appkey, "SessionID:", c.session.SessionId, "Error:", err)
			} else {
				dataMap["snapshot_path"] = path + fmt.Sprintf("%d", now) + ".jpg"
			}
		}
	}
	
	delete(dataMap,"snapShot")
	
	c.archiveJsonData(dataMap, cmdType)
}


func (c *ClientV1) archiveConnectData(connectData *SocketConnect) {
	
	defer func() {
  	if r := recover(); r != nil {
		logger.Fatal(r)
  		}
 	}()
	
	if connectData == nil {
		return
	}

//	var err error

//	data, err := json.Marshal(connectData)
//	if err != nil {
//		logger.Fatal(err)
//		return
//	}

//	dataJson, err := simplejson.NewJson(data)
//	if err != nil {
//		logger.Fatal(err)
//		return
//	}

//	dataMap, err := dataJson.Map()
//	if err != nil {
//		logger.Fatal(err)
//		return
//	}

//	resultJson := simplejson.New()

//	appInfo, err := c.appInfo2Json()
//	if err != nil {
//		logger.Fatal(err)
//		return
//	}

//	deviceInfo, err := c.deviceInfo2Json()

//	if err != nil {
//		logger.Fatal(err)
//		return
//	}

//	resultJson.Set("app", appInfo)
//	resultJson.Set("device", deviceInfo)
//	resultJson.Set("createtime", time.Now().UnixNano())
//	resultJson.Set("socket_connect", dataMap)
//	resultJson.Set("datatype", CmdType_SOCKET_CONNECT)
	
////	if c.network.NetworkName != nil &&
////		(*c.network.NetworkName) != "" {
////		resultJson.Set("network", (*c.network.NetworkName))
////	}
	
//	if c.network.NetworkName != nil {
//		resultJson.Set("network", (*c.network.NetworkName))
//		resultJson.Set("carrier", (*c.network.CarrierName))
//	}
	
//	if 	c.String() != "" {
//		resultJson.Set("clientip", c.String())
//	}
	
//	if (*c.session.SessionId) != "" {
//		resultJson.Set("sessionid", (*c.session.SessionId))
//	}
	
//	if (*c.session.SdkVersion) != "" {
//		resultJson.Set("sdkversion", (*c.session.SdkVersion))
//	}	

//	data_str, err := resultJson.MarshalJSON()

//	if err != nil {
//		logger.Fatal(err)
//		return
//	}

//	if c.Arch != nil {
//		c.Arch.PushData(string(data_str))
//	}
}
                                                                                                                                                                                         
func (c *ClientV1) archiveSendrecvData(sendrecvData *SocketSendRecvData) {
	
	defer func() {
  	if r := recover(); r != nil {
		logger.Fatal(r)
  		}
 	}()
	
	if sendrecvData == nil {
		return
	}
//	var err error

//	data, err := json.Marshal(sendrecvData)
//	if err != nil {
//		logger.Fatal("json.Marshal(sendrecvData)", err)
//		return
//	}

//	dataJson, err := simplejson.NewJson(data)
//	if err != nil {
//		logger.Fatal(err)
//		return
//	}

//	dataMap, err := dataJson.Map()
//	if err != nil {
//		logger.Fatal(err)
//		return
//	}

//	resultJson := simplejson.New()

//	appInfo, err := c.appInfo2Json()
//	if err != nil {
//		logger.Fatal(err)
//		return
//	}

//	deviceInfo, err := c.deviceInfo2Json()

//	if err != nil {
//		logger.Fatal(err)
//		return
//	}

//	resultJson.Set("app", appInfo)
//	resultJson.Set("device", deviceInfo)
//	resultJson.Set("createtime", time.Now().UnixNano())
//	resultJson.Set("socket_sendrecv", dataMap)
//	resultJson.Set("datatype", CmdType_SOCKET_SENDRECV_DATA)
	
////	if c.network.NetworkName != nil &&
////		(*c.network.NetworkName) != "" {
////		resultJson.Set("network", (*c.network.NetworkName))
////	}
	
//	if c.network.NetworkName != nil {
//		resultJson.Set("network", (*c.network.NetworkName))
//		resultJson.Set("carrier", (*c.network.CarrierName))
//	}
	
//	if 	c.String() != "" {
//		resultJson.Set("clientip", c.String())
//	}
	
//	if (*c.session.SessionId) != "" {
//		resultJson.Set("sessionid", (*c.session.SessionId))
//	}
	
//	if (*c.session.SdkVersion) != "" {
//		resultJson.Set("sdkversion", (*c.session.SdkVersion))
//	}		

//	data_str, err := resultJson.MarshalJSON()

//	if err != nil {
//		logger.Fatal(err)
//		return
//	}

//	if c.Arch != nil {
//		c.Arch.PushData(string(data_str))
//	}
}

func (c *ClientV1) archiveMetricData(webData *MetricData) {
	
	defer func() {
  	if r := recover(); r != nil {
		logger.Fatal(r)
  		}
 	}()
	
	
	if webData == nil {
		return
	}
//	var err error
	
//	data, err := json.Marshal(webData)
//	if err != nil {
//		logger.Fatal("json.Marshal(MetricData)", err)
//		return
//	}
	
//	dataJson, err := simplejson.NewJson(data)
//	if err != nil {
//		logger.Fatal(err)
//		return
//	}

//	dataMap, err := dataJson.Map()
//	if err != nil {
//		logger.Fatal(err)
//		return
//	}

//	resultJson := simplejson.New()

//	appInfo, err := c.appInfo2Json()
//	if err != nil {
//		logger.Fatal(err)
//		return
//	}

//	deviceInfo, err := c.deviceInfo2Json()

//	if err != nil {
//		logger.Fatal(err)
//		return
//	}

//	resultJson.Set("app", appInfo)
//	resultJson.Set("device", deviceInfo)
//	resultJson.Set("createtime", time.Now().UnixNano())
//	resultJson.Set("metric", dataMap)
	
//	t2 := time.Now().Sub(c.ClientStartTime)	
//	resultJson.Set("timebond", int64(t2.Seconds()))	
//	resultJson.Set("datatype", CmdType_METRIC_DATA)
	
////	if c.network.NetworkName != nil &&
////		(*c.network.NetworkName) != "" {
////		resultJson.Set("network", (*c.network.NetworkName))
////	}
	
//	if c.network.NetworkName != nil {
//		resultJson.Set("network", (*c.network.NetworkName))
//		resultJson.Set("carrier", (*c.network.CarrierName))
//	}
	
//	if 	c.String() != "" {
//		resultJson.Set("clientip", c.String())
//	}
	
//	if (*c.session.SessionId) != "" {
//		resultJson.Set("sessionid", (*c.session.SessionId))
//	}
	
//	if (*c.session.SdkVersion) != "" {
//		resultJson.Set("sdkversion", (*c.session.SdkVersion))
//	}		

//	data_str, err := resultJson.MarshalJSON()

//	if err != nil {
//		logger.Fatal(err)
//		return
//	}

//	if c.Arch != nil {
//		c.Arch.PushData(string(data_str))
//	}
}

