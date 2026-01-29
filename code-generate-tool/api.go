package main

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func ApiCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "api code generator",
		Long:  "api code generator",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(ApiGenerateCMD())
	return cmd
}

func ApiGenerateCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "generate api",
		Long:  "generate api",
		Run: func(cmd *cobra.Command, args []string) {
			runPath, e := filepath.Abs("../api")
			if e != nil {
				panic(e.Error())
			}
			runCMD := exec.Command("goctl", "api", "go", "-api", "cashflow.api", "-dir", ".")
			runCMD.Dir = runPath
			res, e := runCMD.CombinedOutput()
			fmt.Print(string(res))
			if e != nil {
				panic(e.Error())
			}
		},
	}
}
