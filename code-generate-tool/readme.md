## code-generate-tool使用说明

code-generate-tool，下面简称cgt

cgt使用须在code-generate-tool目录下运行

注意，cgt仅供生成go代码使用

### prepare

使用前，你需要完成[goctl](https://go-zero.dev/en/docs/tasks/installation/goctl)和[protoc](https://go-zero.dev/en/docs/tasks/installation/protoc)的安装

你需要调整api.go中api文件的文件名

当然你也可以根据你的业务需要调整路径等等，目前的结构设计是

- project
    - api
        - *.api（api文件）
    - app
        - rpcName
            - *.rpc（rpc文件）
    - code-generate-tool
        - cgt
        - （在这个目录下运行cgt）

### 编译

``` bash
go build -o cgt
```

### rpc部分

#### 新建一个rpc定义文件

``` bash
cgt rpc new [rpcName]
```

例如：`cgt rpc new hello`

然后rpc的定义文件会被保存在`app/{rpcName}.proto`

#### 生成rpc代码

``` bash
cgt rpc generate [rpcName]
```

例如：`cgt rpc generate hello`

app/{rpcName}.proto后运行，会在`app/{rpcName}/`下生成业务代码，完成`internal/logic/`下的代码即可

### api部分

#### 生成api

``` bash
cgt api generate
```

修改`api/*.api`后运行，会在`api/`下生成业务代码，完成`internal/logic/`下的代码即可