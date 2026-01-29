## deploy-tool

快速跑起来所有服务

### 配置

参考[config.example.json](./config.example.json)和[config.go](./config.go)

你需要创建一个`config.json`文件

build项：配置path数组，每一项就是编译的内容

dev项：配置键值对，每一项的名字对应服务名，path是运行路径，command是启动命令，port是使用的端口号列表，deploy-tool会通过这个是否被使用判断是否启动成功，wait是一个服务名列表，表示依赖那个服务，等那个服务启动了才会开始尝试启动这个服务，env是环境变量列表

stop项，要停止的可执行应用名

### 编译每一项

```
./deploy-tool build
```

### 启动每一项

```
./deploy-tool dev
```

### 停止每一项

```
./deploy-tool stop
```

### 清理所有日志

```
./deploy-tool clean
```