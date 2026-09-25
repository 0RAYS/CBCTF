package main

import (
	"CBCTF/internal/generatorworker"
	"log"
	"os"
)

func main() {
	var err error
	if len(os.Args) == 3 && os.Args[1] == "--install" {
		err = generatorworker.Install(os.Args[2])
	} else {
		err = generatorworker.Serve(os.Getenv("CBCTF_WORKER_TOKEN"))
	}
	if err != nil {
		log.Fatal(err)
	}
}
