# make命令 一览

> make命令可以快速执行一些封装好的命令，如build、kill、run等。

* `make build` 编译为golang程序，编译后的可执行文件在bin/目录
* `make build-linux` 编译为可在Linux上执行的golang程序，编译后的可执行文件在bin/目录
* `make build-all` 编译为可在Linux、Windows、OSX上执行的golang程序，编译后的可执行文件在bin/目录
* `make clean` 删除所有编译后的可执行文件，即清空bin/目录
* `make kill` 关闭正在运行的nsqproxy进程
* `make test` 执行go test
* `make run` 运行 nohup ./bin/nsqproxy &
* `make statik` 将静态资源文件编译成go文件。即statik -src=web/public/ -dest=internal -f
* `make vue-install` 安装VUE，即npm install
* `make vue-install-taobao` 同make vue-install，使用淘宝的源进行安装，防止官方源被墙
* `make vue-build` 将VUE文件编译打包并复制到web/public/目录下
* `make vue-dev` 将VUE文件编译打包并复制到web/public/目录下