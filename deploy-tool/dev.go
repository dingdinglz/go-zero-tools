package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// ServiceStatus 服务状态
type ServiceStatus struct {
	Name    string
	Started bool
	Error   error
	Ports   []int
	mu      sync.RWMutex
}

// SetStarted 设置服务启动状态
func (s *ServiceStatus) SetStarted(started bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Started = started
	s.Error = err
}

// IsStarted 获取服务启动状态
func (s *ServiceStatus) IsStarted() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Started
}

// GetError 获取服务错误信息
func (s *ServiceStatus) GetError() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Error
}

func DevCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "dev",
		Short: "dev run all target",
		Long:  "dev run all target",
		Run: func(cmd *cobra.Command, args []string) {
			// 启动前检查端口占用情况
			fmt.Println("==================== 启动前检查端口占用 ====================")
			for name, devConfig := range config.Dev {
				for _, port := range devConfig.Port {
					if CheckPortInUse(port) {
						panic(fmt.Sprintf("❌ [%s] 端口 %d 已被占用，无法启动", name, port))
					}
				}
			}
			fmt.Println("✅ 所有端口检查通过，无占用")

			// 初始化服务状态
			serviceStatus := make(map[string]*ServiceStatus)
			for name, devConfig := range config.Dev {
				serviceStatus[name] = &ServiceStatus{
					Name:    name,
					Started: false,
					Ports:   devConfig.Port,
				}
			}

			// 使用 WaitGroup 等待所有服务启动完成
			var wg sync.WaitGroup
			timeout := time.Duration(config.Timeout) * time.Second
			if timeout <= 0 {
				timeout = 60 * time.Second // 默认超时60秒
			}

			// 启动所有服务
			fmt.Println("==================== 开始启动所有服务 ====================")
			for name, devConfig := range config.Dev {
				wg.Add(1)
				go startService(name, devConfig, serviceStatus, &wg, timeout)
			}

			// 等待所有服务启动完成
			wg.Wait()

			// 最终验证所有服务的端口
			fmt.Println("\n==================== 最终验证所有服务状态 ====================")
			allSuccess := true
			for name, status := range serviceStatus {
				if status.IsStarted() && CheckAllPortsInUse(status.Ports) {
					fmt.Printf("✅ [%s] 启动成功，所有端口正常监听: %v\n", name, status.Ports)
				} else {
					allSuccess = false
					if err := status.GetError(); err != nil {
						fmt.Printf("❌ [%s] 启动失败: %v\n", name, err)
					} else {
						fmt.Printf("❌ [%s] 启动失败: 端口未正常监听 %v\n", name, status.Ports)
					}
				}
			}

			if allSuccess {
				fmt.Println("\n🎉 所有服务启动成功！")
				fmt.Println("所有服务已在后台运行，deploy-tool 即将退出。")
			} else {
				fmt.Println("\n⚠️  部分服务启动失败，请查看日志文件排查问题。")
			}
		},
	}
}

// startService 启动单个服务
func startService(name string, devConfig DevConfig, serviceStatus map[string]*ServiceStatus, wg *sync.WaitGroup, timeout time.Duration) {
	defer wg.Done()

	status := serviceStatus[name]

	// 等待依赖项启动
	if len(devConfig.Wait) > 0 {
		fmt.Printf("⏳ [%s] 等待依赖项启动: %v\n", name, devConfig.Wait)
		if !waitForDependencies(name, devConfig.Wait, serviceStatus, timeout) {
			panic(fmt.Sprintf("❌ [%s] 依赖项启动失败或超时", name))
		}
		fmt.Printf("✅ [%s] 所有依赖项已就绪\n", name)
	}

	// 创建日志文件
	logFileName := fmt.Sprintf("%s.out", name)
	logFile, err := os.Create(logFileName)
	if err != nil {
		panic(fmt.Sprintf("❌ [%s] 创建日志文件失败: %v", name, err))
	}
	defer logFile.Close()

	// 启动服务
	fmt.Printf("🚀 [%s] 正在启动服务...\n", name)
	startCommand := exec.Command(devConfig.Command[0], devConfig.Command[1:]...)
	workPath, _ := filepath.Abs(devConfig.Path)
	startCommand.Dir = workPath
	startCommand.Env = append(os.Environ(), devConfig.Env...)
	startCommand.Stdout = logFile
	startCommand.Stderr = logFile

	err = startCommand.Start()
	if err != nil {
		panic(fmt.Sprintf("❌ [%s] 启动命令执行失败: %v", name, err))
	}

	// 等待端口监听
	fmt.Printf("⏳ [%s] 等待端口监听: %v\n", name, devConfig.Port)
	if waitForPorts(name, devConfig.Port, timeout) {
		status.SetStarted(true, nil)
		fmt.Printf("✅ [%s] 启动成功，端口已监听: %v\n", name, devConfig.Port)
	} else {
		panic(fmt.Sprintf("❌ [%s] 启动超时，端口未监听: %v (请查看 %s 日志)", name, devConfig.Port, logFileName))
	}
}

// waitForDependencies 等待依赖项启动完成
func waitForDependencies(serviceName string, dependencies []string, serviceStatus map[string]*ServiceStatus, timeout time.Duration) bool {
	startTime := time.Now()
	for {
		if time.Since(startTime) > timeout {
			return false
		}

		allReady := true
		for _, dep := range dependencies {
			depStatus, exists := serviceStatus[dep]
			if !exists {
				fmt.Printf("⚠️  [%s] 依赖项 %s 不存在于配置中\n", serviceName, dep)
				return false
			}

			if !depStatus.IsStarted() {
				allReady = false
				break
			}

			// 检查依赖项的端口是否都在监听
			if !CheckAllPortsInUse(depStatus.Ports) {
				allReady = false
				break
			}
		}

		if allReady {
			return true
		}

		time.Sleep(500 * time.Millisecond)
	}
}

// waitForPorts 等待端口开始监听
func waitForPorts(serviceName string, ports []int, timeout time.Duration) bool {
	if len(ports) == 0 {
		return true // 如果没有配置端口，默认认为启动成功
	}

	startTime := time.Now()
	for {
		if time.Since(startTime) > timeout {
			return false
		}

		if CheckAllPortsInUse(ports) {
			return true
		}

		time.Sleep(500 * time.Millisecond)
	}
}
