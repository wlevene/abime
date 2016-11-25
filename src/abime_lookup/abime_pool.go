package abime_lookup

type AbimePool struct {
	pools []Abime
}

func NewAbimePool() *AbimePool {
	p := &AbimePool{}
	p.pools = make([]Abime, 0, 5)

	return p
}

func (p *AbimePool) Abimes() []Abime {
	return p.pools
}

func (p *AbimePool) Count() int {
	return len(p.pools)
}

func (p *AbimePool) RegisterAbime(abime Abime) (bool, bool) {
	exist := false
	result := false
	if abime.Addr == "" {
		return result, exist
	}

	if p.pools == nil {
		p.pools = make([]Abime, 0, 5)
	}

	for _, saved_abime := range p.pools {
		if saved_abime.Addr == abime.Addr {
			result = true
			exist = true
			return result, exist
		}
	}

	result = true

	p.pools = append(p.pools, abime)

	return result, exist
}

func (p *AbimePool) UnRegisterAbime(abime Abime) {

	if abime.Addr == "" {
		return
	}

	if p.pools == nil {
		return
	}

	index := 0
	for _, saved_abime := range p.pools {
		if saved_abime.Addr == abime.Addr {
			p.pools = append(p.pools[:index], p.pools[index:]...)
			return
		}
		index++
	}
}

func (p *AbimePool) UpdateAbime(abime Abime) {
	if abime.Addr == "" {
		return
	}

	index := 0
	for _, saved_abime := range p.pools {
		if saved_abime.Addr == abime.Addr {
			//			saved_abime = abime
			copyPool := p.pools[:]
			p.pools = append(copyPool[:index], abime)
			p.pools = append(p.pools, copyPool[index+1:]...)
			copyPool = copyPool[0:0]
			return
		}
		index++
	}
}

func (p *AbimePool) HeavyestAbime() Abime {

	var ret Abime

	if len(p.pools) <= 0 {
		return ret
	}

	for _, abime := range p.pools {
		if ret.Addr == "" {
			ret = abime
			continue
		}

		if abime.ClientCount > ret.ClientCount {
			ret = abime
		}
	}

	return ret
}

func (p *AbimePool) LessHeavyestAbime() Abime {

	var ret Abime

	if len(p.pools) <= 0 {
		return ret
	}

	for _, abime := range p.pools {
		if ret.Addr == "" {
			ret = abime
			continue
		}

		if abime.ClientCount < ret.ClientCount {
			ret = abime
		}
	}

	return ret
}

func (p *AbimePool) LessHeavyestAbimeAboveAt(abouveCount int64) Abime {

	var ret Abime

	if len(p.pools) <= 0 {
		return ret
	}

	for _, abime := range p.pools {
		
		if abime.ClientCount < abouveCount {
			continue
		}
				
		if ret.Addr == "" {
			ret = abime
			continue
		}

		if abime.ClientCount < ret.ClientCount {
			ret = abime
		}
	}
	
	if ret.Addr == "" {
		ret = p.HeavyestAbime()
		
		if ret.ClientCount > abouveCount {
			ret = p.LessHeavyestAbime()
		}
	}

	return ret
}
