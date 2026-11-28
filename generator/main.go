package main

import (
	"log"
	"os"

	"CBCTF/internal/generator"
)

func main() {
	var err error
	if len(os.Args) == 3 && os.Args[1] == "--install" {
		err = generator.Install(os.Args[2])
	} else {
		err = generator.Serve(os.Getenv("CBCTF_GENERATOR_TOKEN"))
	}
	if err != nil {
		log.Fatal(err)
	}
}
