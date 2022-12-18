package main

import (
	"log"

	"github.com/yskszk63/devcontainer-shell/agent-cmd"
)

func main() {
	err := agentcmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
