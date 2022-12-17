package main

import (
	"log"

	"github.com/yskszk63/devcontainer-shell/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
