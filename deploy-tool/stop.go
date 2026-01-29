package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func StopCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "stop old",
		Long:  "stop old",
		Run: func(cmd *cobra.Command, args []string) {
			for _, item := range config.Stop.Name {
				fmt.Println("stoping", item)
				findCommand := exec.Command("pgrep", "-x", item)
				res, _ := findCommand.CombinedOutput()
				killList := strings.Split(string(res), "\n")
				for _, killItem := range killList {
					if killItem == "" {
						continue
					}
					fmt.Println("stopping pid:", killItem)
					stopCommand := exec.Command("kill", killItem)
					_, e := stopCommand.CombinedOutput()
					if e != nil {
						panic(e)
					}
				}
				fmt.Println("killed", item)
			}
		},
	}
}
