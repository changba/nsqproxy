# NSQ Proxy
NSQ Proxy是Golang开发的NSQ和Worker之间的中间件，根据数据库配置，负责消息转发。NSQProxy启动后，接受NSQD队列内容，然后通过HTTP/FastCGI/CBNSQ等协议转发给Worker机执行。在唱吧内部使用2年，高效稳定的处理着每日数十亿条消息。

[![go report card](https://goreportcard.com/badge/github.com/ChangbaServer/nsqproxy "go report card")](https://goreportcard.com/report/github.com/ChangbaServer/nsqproxy)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://github.com/ChangbaServer/nsqproxy/blob/master/LICENSE)
[![Downloads](https://img.shields.io/github/downloads/ChangbaServer/nsqproxy/total.svg)](https://github.com/ChangbaServer/nsqproxy/releases)
[![Release](https://img.shields.io/github/release/ChangbaServer/nsqproxy.svg?label=Release)](https://github.com/ChangbaServer/nsqproxy/releases)

## 解决的问题

* 各Topic执行机器可配
* 各Topic消费速度可配
* 各Worker机协议可配
    * HTTP：将消息发送给配好的URL。
    * FastCGI：将消息发送给配置的服务端，如PHP-FPM。
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
请先部署好NSQLookupd、NSQd、MySQL

> 启动NSQLookupd `nsqlookupd -broadcast-address="0.0.0.0" -http-address="0.0.0.0:4161" -tcp-address="0.0.0.0:4160"`

> 启动NSQD `nsqd -broadcast-address="0.0.0.0" -lookupd-tcp-address="0.0.0.0:4160" -tcp-address="0.0.0.0:4150" -http-address="0.0.0.0:4151"`

> 启动MySQL

### 安装

#### 二进制安装

* 下载最新版本的压缩包 https://github.com/ChangbaServer/nsqproxy/releases
* 解压
* 启动（注意替换为自己的MySQL信息） `./nsqproxy -dbHost=127.0.0.1 -dbPort=3306 -dbUsername=root -dbPassword=rootpsd -dbName=nsqproxy -logLevel=debug -nsqlookupdHTTP=127.0.0.1:4161`
* 命令行 `curl http://0.0.0.0:19421/status` 输出ok
* 浏览器打开 http://0.0.0.0:19421/admin

#### 源码安装

* 要求Go1.13及以上
* 下载本项目 `go get github.com/ChangbaServer/nsqproxy`
* `cd nsqproxy`
* `export GO111MODULE=on`
* 编译 `make build`
* 启动（注意替换为自己的MySQL信息） `./bin/nsqproxy -dbHost=127.0.0.1 -dbPort=3306 -dbUsername=root -dbPassword=rootpsd -dbName=nsqproxy -logLevel=debug -nsqlookupdHTTP=127.0.0.1:4161`
* 命令行 `curl http://0.0.0.0:19421/status` 输出ok
* 浏览器打开 http://0.0.0.0:19421/admin

### 快速开始

* [快速体验](document/doc/quick_start.md)
* [启动参数](document/doc/flag.md)
* [make命令](document/doc/make.md)


## 二次开发

### 前端
使用VUE开发，所有源码均在/web/vue-admin目录中，开发完成后需要编译，编译后的文件存放在/web/public/目录中。使用开源项目statik将静态文件/web/public/变成一个go文件internal/statik/statik.go，这样前端的静态文件也会被我们编译到同一个二进制文件中了。

* 启动go服务 `make run`
* 安装VUE `make vue-install`（如果国内被墙可以使用淘宝的源进行安装：make vue-install-taobao）
* 开启VUE开发环境 `make vue-dev`
* 浏览器打开 http://0.0.0.0:9528/admin
* 开发前端相关功能
* 编译VUE `make vue-build`
* 前段文件转换为一个go文件 `make statik`
* 编译go服务 `make build`
* 浏览器打开 http://0.0.0.0:19421/admin

### 接口文档
* 通过接口对数据库增删改查：[查看接口文档](document/api/README.md)

## TODO LIST

* 协议增加protobuf
* 后台增加用户权限管理
* 报警HOOK
* 日志按天分割

## License

© [Changba.com](https://changba.com), 2020~time.Now

Released under the [MIT License](https://github.com/ChangbaServer/nsqproxy/blob/main/LICENSE)