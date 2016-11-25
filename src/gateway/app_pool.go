package gateway

import (
	"fmt"
)

type AppPool struct {
	pool map[string]*AppConfig
}

func NewAppPool() *AppPool {
	p := &AppPool{}

	return p
}


func (p *AppPool) GetAppConfig(appkey string) *AppConfig {
	if p.pool == nil {
		p.pool = make(map[string]*AppConfig)
	}

	if appkey == "" {
		return nil
	}

	ret := p.pool[appkey]
	if ret != nil {
		return ret
	}

	ret = NewAppConfig(appkey)

	p.pool[appkey] = ret

	return ret
}

func (p *AppPool) VailedMetric(appkey string) bool {
	
	if appkey == "" {
		return false
	}

	if p.pool == nil {
		return false
	}
	
	appConfig := p.pool[appkey]
	if appConfig == nil {
		return false
	}

	fmt.Println("------")
	result := appConfig.vailedMetricAndCount()

	return result
}
