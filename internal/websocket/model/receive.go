package model

type Receive struct {
	Type string `json:"type"`
}

type HeartbeatRecv struct {
	Msg string `json:"msg"`
	Receive
}
