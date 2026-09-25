package main

import (
	"log"
	"os"

	"CBCTF/internal/worker"
)

func main() {
	var err error
	if len(os.Args) == 3 && os.Args[1] == "--install" {
		err = worker.Install(os.Args[2])
	} else {
		err = worker.Serve(os.Getenv("CBCTF_WORKER_TOKEN"))
	}
	if err != nil {
		log.Fatal(err)
	}
}
