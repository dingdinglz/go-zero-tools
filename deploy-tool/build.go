package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/spf13/cobra"
)

func BuildCMD() *cobra.Command {
	buildCMD := &cobra.Command{
		Use:   "build",
		Short: "build all targets",
		Long:  "build all targets",
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup
			for _, item := range config.Build.Path {
				fmt.Println("开始build", item)
				wg.Add(1)
				go func() {
					buildPath, _ := filepath.Abs(item)
					buildExec := exec.Command("go", "build")
					buildExec.Dir = buildPath
					res, e := buildExec.CombinedOutput()
					if e != nil {
						fmt.Println(item, "!!!", string(res))
						panic(e)
					}
					fmt.Println("build成功", item)
					wg.Done()
				}()
			}
			wg.Wait()
		},
	}
	return buildCMD
}
