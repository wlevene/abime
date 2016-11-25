package gateway



var MaxReportMetricCliecntCountDay int = 3

type AppConfig struct {
	Appkey    		string
	VailedApp 		bool
	VailedMetric 	bool
	ReportMetricCliecntCountToday int
}

func NewAppConfig(appkey string) *AppConfig {

	p := &AppConfig{}

	p.Appkey = appkey
	p.VailedApp = true
	p.VailedMetric = true
	p.ReportMetricCliecntCountToday = 0;
		
	return p
}

func (p *AppConfig) vailedMetricAndCount() bool {
	
	
	if p.VailedApp == false {
		return false
	}
	
	if p.VailedMetric == false {
		return false
	}
	
	if p.ReportMetricCliecntCountToday >= MaxReportMetricCliecntCountDay {
		return false
	}
	
	if p.ReportMetricCliecntCountToday < 0 {
		return false
	}
	
	p.ReportMetricCliecntCountToday++
	
	return true
	
}