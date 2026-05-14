package main

import (
	"os"

	"gonode/pkg/engine"
	"gonode/pkg/utils"
)

func main() {

	if len(os.Args) < 2 {
		utils.PrintHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "start":
		engine.HandleStartCommand()
	case "list":
		engine.SendCommand("list")
	case "stop":
		engine.SendCommand("stop")
	case "help":
		utils.PrintHelp()
	case "daemon-logic":
		engine.RunDaemonLogic()
	default:
		utils.PrintHelp()
	}
}