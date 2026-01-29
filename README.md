## go-zero-tools

一些用于加速使用go-zero框架的项目的开发的小工具，包括但不限于代码生成、编译、运行等等，可以快速集成到现有的项目或者新项目中 

### List

- [code-generate](./code-generate-tool)：用于快速根据api和rpc文件生成代码，是对goctl的二次封装，将该工具的readme喂给ai可以快速让ai学会根据api文件和rpc文件生成代码，而不需要额外传入goctl的使用方法，同时可以规范项目结构
- [deploy-tool](./deploy-tool)：可以快速把所有本地服务和网关跑起来，方便本地测试，也可以作为测试服务器上CD流的一部分