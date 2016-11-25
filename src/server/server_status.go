package server

type SvrStatus struct {
	CurrentClientCount int64 `json:"current_client_count"`
	RunTime string `json:"runtime"`
	Name string `json:"name"`
	MsgCount uint64 `json:"msg_count"`
}
