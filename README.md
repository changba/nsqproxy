# NSQ Proxy
NSQ Proxy是Golang开发的NSQ和Worker之间的中间件，根据数据库配置，负责消息转发。NSQProxy启动后，接受NSQD队列内容，然后通过HTTP/FastCGI/CBNSQ等协议转发给Worker机执行。

[![go report card](https://goreportcard.com/badge/github.com/ChangbaServer/nsqproxy "go report card")](https://goreportcard.com/report/github.com/ChangbaServer/nsqproxy)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://github.com/ChangbaServer/nsqproxy/blob/master/LICENSE)
[![Downloads](https://img.shields.io/github/downloads/ChangbaServer/nsqproxy/total.svg)](https://github.com/ChangbaServer/nsqproxy/releases)
[![Release](https://img.shields.io/github/release/ChangbaServer/nsqproxy.svg?label=Release)](https://github.com/ChangbaServer/nsqproxy/releases)

## 解决的问题

* 各Topic执行机器可配
* 各Topic消费速度可配
* 各Worker机协议可配
    * HTTP：将消息发送给配好的URL。
    * FastCGI：将消费发送给配置的服务端，如PHP-FPM。
    * CBNSQ：自定义的基于TCP的文本协议。
* 可视化界面管理
* 队列积压超出阈值报警
* 散乱在各处的消费者集中化管理
* 通过网络分发，无需安装.so等扩展库，因此无需修改线上环境

## 有图有真相

<img src="https://raw.githubusercontent.com/ChangbaServer/nsqproxy/main/assets/images/nsqproxy_flow_chart.png" alt="流程图">

<img src="https://raw.githubusercontent.com/ChangbaServer/nsqproxy/main/assets/images/admin_consume_config.png" alt="消费者管理">

<img src="https://raw.githubusercontent.com/ChangbaServer/nsqproxy/main/assets/images/admin_work_server.png" alt="worker机管理">

## 使用

### 安装

* 要求Go1.11及以上
* 下载本项目
* `go get github.com/ChangbaServer/nsqproxy`
* `cd nsqproxy`
* `go run cmd/nsqproxy.go -dbHost=127.0.0.1 -dbPort=3306 -dbUsername=root -dbPassword=rootpsd -dbName=nsqproxy -logLevel=debug -nsqlookupdHTTP=127.0.0.1:4161`
* 命令行 `curl http://0.0.0.0:19421/status` 输出ok
* 浏览器打开 http://0.0.0.0:19421/admin

### 启动参数
启动时命令行传入参数。所有的参数都有默认值。

`-httpAddr string` 监听的HTTP端口 (default "0.0.0.0:19421")

`-masterAddr string` 主库IP端口，为空则本机为主机

`-logLevel string` 日志等级，可选有debug、info、warning、error、fatal (default "info")

`-logPath string` 系统日志路径 (default "logs/proxy.log")

`-subLogPath string` 消费log，消费详情由于量大成功消费log仅在日志等级为debug时启用 (default "logs/sub.log")

`-nsqlookupdHTTP string` nsqLookupd的HTTP地址，多个用逗号分割如"127.0.0.1:4161,127.0.0.1:4163" (default "127.0.0.1:4161")

`-updateConfigInterval int` 定时向Mysql更新消费者配置的间隔时间，单位秒 (default 60)

`-dbHost string` Mysql的IP (default "127.0.0.1")

`-dbPort string` Mysql的端口 (default "3306")

`-dbPassword string` Mysql的密码 (default "")

`-dbUsername string` Mysql的账号 (default "root")

`-dbName string` Mysql的库名 (default "nsqproxy")

### make命令

`make build` 编译为golang程序，编译后的可执行文件在bin/目录

`make build-linux` 编译为可在Linux上执行的golang程序，编译后的可执行文件在bin/目录

`make build-all` 编译为可在Linux、Windows、OSX上执行的golang程序，编译后的可执行文件在bin/目录

`make clean` 删除所有编译后的可执行文件，即清空bin/目录

`make kill` 关闭正在运行的nsqproxy进程

`make test` 执行go test

`make run` 运行 nohup ./bin/nsqproxy &

`make statik` 将静态资源文件编译成go文件。即statik -src=web/public/ -dest=internal -f

`make vue-build` 将VUE文件编译打包并复制到web/public/目录下

`make vue-install` 安装VUE，即npm install

`make vue-install-taobao` 同make vue-install，使用淘宝的源进行安装，防止官方源被墙

## 二次开发

### 前端
使用VUE开发，所有源码均在/web/vue-admin目录中，开发完成后需要编译，编译后的文件存放在/web/public/目录中。使用开源项目statik将静态文件/web/public/变成一个go文件internal/statik/statik.go，这样前端的静态文件也会被我们编译到同一个二进制文件中了。

* make vue-install（如果国内被墙可以使用淘宝的源进行安装：make vue-install-taobao）
* make vue-build
* make statik

### 接口文档
* 通过接口对数据库增删改查：[查看接口文档](document/api/index.md)
* 支持的下发给Worker机的各协议说明：[查看协议文档](document/protocol/index.md)

## TODO LIST

* 协议增加protobuf
* 后台增加用户权限管理
* 报警HOOK
* 日志按天分割

## License

© [Changba.com](https://changba.com), 2020~time.Now

Released under the [MIT License](https://github.com/ChangbaServer/nsqproxy/blob/main/LICENSE)