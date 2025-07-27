package model

type Send struct {
	Type  string `json:"type"`
	Msg   string `json:"msg"`
	Title string `json:"title"`
}

type HeartbeatSend Send
