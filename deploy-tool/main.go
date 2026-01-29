package main

import (
	"github.com/spf13/cobra"
)

func main() {
	LoadConfig()

	rootCommand := &cobra.Command{
		Use:   "deploy-tool",
		Short: "build & deploy the project",
		Long:  "build & deploy the project",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCommand.AddCommand(BuildCMD())
	rootCommand.AddCommand(DevCMD())
	rootCommand.AddCommand(StopCMD())
	rootCommand.AddCommand(CleanCMD())

	rootCommand.Execute()
}
