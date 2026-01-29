package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	// 将首字母转为大写，然后拼接上剩余的字符串
	return strings.ToUpper(string(s[0])) + s[1:]
}

func RpcCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rpc",
		Short: "rpc code generator",
		Long:  "rpc code generator",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(RpcNewCMD())
	cmd.AddCommand(RpcGenerateCMD())

	return cmd
}

func RpcNewCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "new",
		Short: "new rpc code",
		Long:  "new rpc code",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("usage: cgt rpc new [功能英文名]")
				return
			}
			p, _ := filepath.Abs("..")
			dir := filepath.Join(p, "app", args[0])

			e := os.MkdirAll(dir, os.ModePerm)
			if e != nil {
				panic(e.Error())
			}

			p = filepath.Join(p, "app", args[0], args[0]+".proto")
			rpcFileContent := `syntax = "proto3";
option go_package = ".";
package {{a}}.v1;

// gRPC服务定义
service {{b}}Service {

}
`
			rpcFileContent = strings.ReplaceAll(rpcFileContent, "{{a}}", args[0])
			rpcFileContent = strings.ReplaceAll(rpcFileContent, "{{b}}", capitalizeFirst(args[0]))
			e = os.WriteFile(p, []byte(rpcFileContent), os.ModePerm)
			if e != nil {
				panic(e.Error())
			}
			fmt.Println("Done.")
		},
	}
}

func RpcGenerateCMD() *cobra.Command {
	/*
		构造类似：goctl rpc protoc ./proto/agent/v1/agent.proto --go_out=./gen/agent/v1/ --go-grpc_out=./gen/agent/v1/ --zrpc_out=./app/agent/
	*/
	return &cobra.Command{
		Use:   "generate",
		Short: "generate rpc code",
		Long:  "generate rpc code",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("usage: cgt rpc generate [功能英文名]")
				return
			}

			runPath, e := filepath.Abs("..")
			if e != nil {
				panic(e.Error())
			}
			runPath = filepath.Join(runPath, "app", args[0])

			runCMD := exec.Command("goctl", "rpc", "protoc", args[0]+".proto", "--go_out", "./pb", "--go-grpc_out", "./pb", "--zrpc_out", ".")
			runCMD.Dir = runPath
			res, e := runCMD.CombinedOutput()
			fmt.Print(string(res))
			if e != nil {
				panic(e.Error())
			}
		},
	}
}
