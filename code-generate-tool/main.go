package main

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cgt",
		Short: "backend code generator",
		Long:  "backend code generator , designed to generate rpc and so on",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.AddCommand(RpcCMD())
	rootCmd.AddCommand(ApiCMD())

	e := rootCmd.Execute()
	if e != nil {
		panic(e.Error())
	}
}
