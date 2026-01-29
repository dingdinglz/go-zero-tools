package main

import (
	"net"
	"os"
	"strconv"
)

func FileExist(_path string) bool {
	_, e := os.Stat(_path)
	return e == nil
}

// CheckPortInUse 检查指定端口是否被占用 (TCP)
func CheckPortInUse(port int) bool {
	// 尝试监听该端口
	// "tcp" 表示检查 TCP 端口，也可以改为 "udp"
	// fmt.Sprintf(":%d", port) 相当于监听 0.0.0.0:port
	address := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", address)

	if err != nil {
		// 如果监听失败，大概率是因为端口被占用
		// (也有可能是权限不足，比如监听 < 1024 的端口但没有 root 权限)
		return true
	}

	// 如果监听成功，必须立刻关闭，以免真正占用该端口
	listener.Close()
	return false
}

// CheckAllPortsInUse 检查所有端口是否都被占用
func CheckAllPortsInUse(ports []int) bool {
	for _, port := range ports {
		if !CheckPortInUse(port) {
			return false
		}
	}
	return true
}
