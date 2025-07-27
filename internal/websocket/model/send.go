package model

type Send struct {
	Level string `json:"level"`
	Type  string `json:"type"`
	Msg   string `json:"msg"`
	Title string `json:"title"`
}

type HeartbeatSend Send
