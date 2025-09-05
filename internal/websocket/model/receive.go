package model

type Recv struct {
	Type string `json:"type"`
}

type HeartbeatRecv struct {
	Msg string `json:"msg"`
	Recv
}
