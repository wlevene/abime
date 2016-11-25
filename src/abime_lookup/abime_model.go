package abime_lookup

type Abime struct {
	Addr        	string `json:"addr"`
	ClientCount 	int64  `json:"clientcount"`
	AbimeName   	string `json:"name"`
	MessageCount   	uint64 `json:"msgcount"`
	RunTime 		string `json:"runtime"`

	ReportTime int64 `json:"reporttime"`
}
