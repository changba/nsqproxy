# 启动参数 一览

> 启动时命令行可以传入的参数，所有的参数都有默认值。

> `./bin/nsqproxy -h` 也可查看

* NSQProxy相关部分
    * `-httpAddr string` 监听的HTTP端口 (default "0.0.0.0:19421")
    * `-masterAddr string` 主库IP端口，为空则本机为主机，不为空则本机为备机。本机为备机时，会定期给masterAddr发PING，连续5次未收到PONG则认定主机异常，该备机启动。(default "")
    * `-logLevel string` 日志等级，可选有debug、info、warning、error、fatal (default "info")
    * `-logPath string` 系统日志路径 (default "logs/proxy.log")
    * `-subLogPath string` 消费log，消费详情由于量大成功消费log仅在日志等级为debug时启用 (default "logs/sub.log")
    * `-updateConfigInterval int` 定时向Mysql更新消费者配置的间隔时间，单位秒 (default 60)
* NSQ相关部分
    * `-nsqlookupdHTTP string` nsqLookupd的HTTP地址，多个用逗号分割如"127.0.0.1:4161,127.0.0.1:4163" (default "127.0.0.1:4161")
* Mysql相关部分
    * `-dbHost string` Mysql的IP (default "127.0.0.1")
    * `-dbPort string` Mysql的端口 (default "3306")
    * `-dbPassword string` Mysql的密码 (default "")
    * `-dbUsername string` Mysql的账号 (default "root")
    * `-dbName string` Mysql的库名 (default "nsqproxy")