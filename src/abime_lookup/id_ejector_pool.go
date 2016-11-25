package abime_lookup

import ()

type IdEjectorPool struct {
	pool map[string]*RedisIdEjector
}

func NewIdEjectorPool() *IdEjectorPool {
	p := &IdEjectorPool{}

	return p
}

func (p *IdEjectorPool) CreateIdEjector(appkey string, startID string) bool {
	if p.pool == nil {
		p.pool = make(map[string]*RedisIdEjector)
	}

	if appkey == "" {
		return false
	}

	ret := p.pool[appkey]
	if ret != nil {
		return true
	}

	ret = NewRedisIDEjector(appkey, startID)

	p.pool[appkey] = ret

	return true
}

func (p *IdEjectorPool) NewID(appkey string) string {

	if appkey == "" {
		return ""
	}

	if p.pool == nil {
		return ""
	}

	idEjector := p.pool[appkey]
	if idEjector == nil {
		return ""
	}

	result := idEjector.NewID()

	return result
}
