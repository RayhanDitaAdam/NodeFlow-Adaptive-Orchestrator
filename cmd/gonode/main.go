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
		engine.ListAllServices()
	case "logs":
		if len(os.Args) < 3 {
			fmt.Println("Usage: gonode logs <project-name>")
			return
		}
		utils.TailLogs(os.Args[2])
	case "stop":
		if len(os.Args) < 3 {
			fmt.Println("Usage: gonode stop <project-name>")
			return
		}
		engine.HandleStopCommand(os.Args[2])
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