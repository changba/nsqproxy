# 快速开始

## 简介
NSQProxy是 队列NSQ 和 消费者 之间的桥梁。负责将 队列中的任务 根据管理后台的配置 通过指定的协议 分发给指定的Worker机。

* 队列中的任务：生产者入队到NSQ的数据。
* 管理后台的配置：每个Topic/Channel都可以配置在哪台机器上消费、同时消费的并发量是几个、是否重新入队、积压多少个开始报警等。
* 通过指定的协议：支持HTTP、FastCGI、CBNSQ等协议将数据发送给Worker机。
* Worker机：执行消费者程序的服务器。
* 数据流转：生产者 -> NSQ -> NSQProxy -> 消费者

> 本文档以HTTP作为通信协议，编写一个go代码作为消费者。

## 启动依赖
NSQProxy是一个中间转发器，因此需要上下游依赖，尽管没有依赖时也可以正常启动。

#### 依赖NSQ
NSQ是真正的队列服务，因此NSQProxy的上游是NSQ。NSQ会将任务下发给NSQProxy，站在NSQ的视角中，NSQProxy是一个真正的消费者。

* 启动NSQLookupd `nsqlookupd -broadcast-address="0.0.0.0" -http-address="0.0.0.0:4161" -tcp-address="0.0.0.0:4160"`
* 启动NSQD `nsqd -broadcast-address="0.0.0.0" -lookupd-tcp-address="0.0.0.0:4160" -tcp-address="0.0.0.0:4150" -http-address="0.0.0.0:4151"`

#### 依赖MySQL
MySQL中存储着各Topic/Channel的配置信息，因此NSQProxy依赖于MySQL。

启动MySQL的方式多种多样，如`mysqld` 和 `service mysql start`等等。

## 下载安装
下载并启动NSQProxy。

* 下载最新版本的压缩包 https://github.com/ChangbaServer/nsqproxy/releases
* 解压
* 启动（注意替换为自己的MySQL信息） `./nsqproxy -dbHost=127.0.0.1 -dbPort=3306 -dbUsername=root -dbPassword=rootpsd -dbName=nsqproxy -logLevel=debug -nsqlookupdHTTP=127.0.0.1:4161`
* 命令行 `curl http://0.0.0.0:19421/status` 输出ok
* 浏览器打开 http://0.0.0.0:19421/admin

## 部署消费者
本文编写一个go代码作为消费者，这个go代码使用HTTP协议监听8888端口。
```golang
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// 启动HTTP
func main() {
	http.HandleFunc("/nsqTask", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("MessageID：" + r.Header.Get("MESSAGE_ID"))
		data, _ := ioutil.ReadAll(r.Body)
		fmt.Println("MessageBody：" + string(data))
		_, _ = w.Write([]byte("200 ok"))
	})
	if err := http.ListenAndServe("0.0.0.0:8888", nil); err != nil {
		panic("ListenAndServe error: " + err.Error())
	}
}
```

运行 `go run test.go`

## 后台操作
把Topic和消费者关联起来。浏览器打开 http://0.0.0.0:19421/admin

1、添加Worker机

<img src="https://raw.githubusercontent.com/ChangbaServer/nsqproxy/main/assets/images/quick_start_add_work_server.png" alt="添加Worker机">

2、添加新消费者配置

<img src="https://raw.githubusercontent.com/ChangbaServer/nsqproxy/main/assets/images/quick_start_add_consume_config.png" alt="添加新消费者配置">

3、把消费者和Worker机关联起来

<img src="https://raw.githubusercontent.com/ChangbaServer/nsqproxy/main/assets/images/quick_start_add_consume_server_map.png" alt="把消费者和Worker机关联起来">

## 具体示例

此时，我们给NSQ入队，NSQ就会把消息推给NSQProxy，NSQProxy根据刚才的配置，就会把消息推送给0.0.0.0:8888的Golang程序。

1、入队给NSQ `curl -d 'name=xiaoming&sex=male&age=18' 'http://0.0.0.0:4151/pub?topic=test_topic'`

2、查看刚才编写的Golang程序的输出，可以拿到消息的唯一ID和消息的内容

<img src="https://raw.githubusercontent.com/ChangbaServer/nsqproxy/main/assets/images/quick_start_demo.png" alt="具体示例">