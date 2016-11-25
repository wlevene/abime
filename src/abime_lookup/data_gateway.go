package abime_lookup

import (
	"time"
	"sync"
	
	"github.com/donnie4w/go-logger/logger"
)

const (
	
	MAX_BATTERY_DEVICE  =  200	
	// 每天凌晨 x 点清空 BatteryGateway数据
	DEL_BATTERY_GATEWAY_EVERYDAY_HOUR = 0
	
	
	
	MAX_FPS_DEVICE  =  200
	// 每天凌晨 x 点清空 FpsGateway数据
	DEL_FPS_GATEWAY_EVERYDAY_HOUR = 0
)

type DataGateway struct {
	
	battery_device_map     	map[string][]string
	batteryGatewayMutex 	sync.Mutex
	
	fps_device_map			map[string][]string
	fpsGatewayMutex			sync.Mutex
	
}

func NewDataGateway() *DataGateway {
	p := &DataGateway{}
	p.battery_device_map = make(map[string][]string)
	
	go func() {
		p.resetBatteryGateway()		
		p.resetFpsGateway()		
	}()

	return p
}

func (gateway * DataGateway)resetBatteryGateway() {
	startTime := time.Now()
	hour := startTime.Hour()
	min := startTime.Minute()
	sec := startTime.Second()
	
	t := time.Tick((24 - time.Duration(hour) + DEL_BATTERY_GATEWAY_EVERYDAY_HOUR) * time.Hour + 
		(60 - time.Duration(min)) * time.Minute + 
		(60 - time.Duration(sec)) * time.Second)
	
	<- t
	gateway.batteryGatewayMutex.Lock()
	logger.Info("Clear Battery Gateway :)")
	gateway.battery_device_map = nil
	gateway.battery_device_map = make(map[string][]string)

	gateway.batteryGatewayMutex.Unlock()
	
	for _ = range time.Tick(24 * time.Hour) {
		gateway.batteryGatewayMutex.Lock()
		logger.Info("Clear Battery Gateway :)")
		gateway.battery_device_map = nil
		gateway.battery_device_map = make(map[string][]string)
		gateway.batteryGatewayMutex.Unlock()
	}
}


func (gateway * DataGateway)resetFpsGateway() {
	
	startTime := time.Now()
	hour := startTime.Hour()
	min := startTime.Minute()
	sec := startTime.Second()
	
	t := time.Tick((24 - time.Duration(hour) + DEL_FPS_GATEWAY_EVERYDAY_HOUR) * time.Hour + 
		(60 - time.Duration(min)) * time.Minute + 
		(60 - time.Duration(sec)) * time.Second)
	
	<- t
	gateway.fpsGatewayMutex.Lock()
	logger.Info("Clear Fps Gateway :)")
	gateway.fps_device_map = nil
	gateway.fps_device_map = make(map[string][]string)

	gateway.fpsGatewayMutex.Unlock()
	
	for _ = range time.Tick(24 * time.Hour) {
		gateway.fpsGatewayMutex.Lock()
		logger.Info("Clear Fps Gateway :)")
		gateway.fps_device_map = nil
		gateway.fps_device_map = make(map[string][]string)
		gateway.fpsGatewayMutex.Unlock()
	}
}


func (gateway * DataGateway)IsPassBatteryGateway(appkey string) bool {
	
	if appkey == "" {
		return false
	}
	
	if gateway.battery_device_map == nil {
		return false
	}
	
	devices := gateway.battery_device_map[appkey]
	
	if len(devices) < MAX_BATTERY_DEVICE {
		logger.Info("Battery APPkey:", appkey, "devices count:", len(devices))
		return true
	}
	
	return false
}

func (gateway * DataGateway)IsPassFpsGateway(appkey string) bool {
	
	if appkey == "" {
		return false
	}
	
	if gateway.fps_device_map == nil {
		return false
	}
	
	devices := gateway.fps_device_map[appkey]
	
	if len(devices) < MAX_FPS_DEVICE {
		logger.Info("FPS APPkey:", appkey, "devices count:", len(devices))
		return true
	}
	
	return false
}


func (gateway * DataGateway)ReportDeviceToBatteryGateway(appkey, device string) {
	
	if appkey == "" {
		return
	}
	
	if device == "" {
		return
	}
	
	gateway.batteryGatewayMutex.Lock()
	
	if gateway.battery_device_map == nil {
		gateway.battery_device_map = make(map[string][]string)
	}
	
	devices := gateway.battery_device_map[appkey]
	
	if len(devices) >= MAX_BATTERY_DEVICE {
		gateway.batteryGatewayMutex.Unlock()		
//		logger.Info("Battey Gateway Is Full appkey:", appkey)
		return
	}
	
	
	for _, save_device := range devices {
		if save_device == device {
			gateway.batteryGatewayMutex.Unlock()
			return
		}
	}
	
	devices = append(devices, device)
	gateway.battery_device_map[appkey] = devices
	
	gateway.batteryGatewayMutex.Unlock()
}



func (gateway * DataGateway)ReportDeviceToFpsGateway(appkey, device string) {
	
	if appkey == "" {
		return
	}
	
	if device == "" {
		return
	}
	
	gateway.fpsGatewayMutex.Lock()
	
	if gateway.fps_device_map == nil {
		gateway.fps_device_map = make(map[string][]string)
	}
	
	devices := gateway.fps_device_map[appkey]
	
	if len(devices) >= MAX_FPS_DEVICE {
		gateway.fpsGatewayMutex.Unlock()		
//		logger.Info("Battey Gateway Is Full appkey:", appkey)
		return
	}
	
	for _, save_device := range devices {
		if save_device == device {
			gateway.fpsGatewayMutex.Unlock()
			return
		}
	}
	
	devices = append(devices, device)
	gateway.fps_device_map[appkey] = devices
	
	gateway.fpsGatewayMutex.Unlock()
}
