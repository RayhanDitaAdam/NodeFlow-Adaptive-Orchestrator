package main

import (
	"fmt"
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
	case "check":
		if len(os.Args) > 3 && os.Args[2] == "propagation" {
			utils.CheckPropagation(os.Args[3], os.Args[4])
		} else {
			fmt.Println("Usage: gonode check propagation <domain> <expected_ip>")
		}
	case "list":
		engine.SendCommand("list")
	case "stop":
		engine.SendCommand("stop")
	case "help":
		if len(os.Args) > 2 && os.Args[2] == "nginx" {
			utils.PrintNginxHelp()
		} else {
			utils.PrintHelp()
		}
	case "daemon-logic":
		engine.RunDaemonLogic()
	default:
		utils.PrintHelp()
	}
}