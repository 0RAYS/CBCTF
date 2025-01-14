package main

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/log"
)

func init() {
	config.Init()
	log.Init()
	db.Init()
}

func main() {
}
