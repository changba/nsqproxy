package main

//go:generate echo "statik -src=../web/public/ -dest=../internal -f"
//go:generate statik -src=../web/public/ -dest=../internal -f

import (
	"github.com/changba/nsqproxy/config"
	"github.com/changba/nsqproxy/internal/backup"
	"github.com/changba/nsqproxy/internal/httper"
	"github.com/changba/nsqproxy/internal/model"
	"github.com/changba/nsqproxy/internal/module/logger"
	"github.com/changba/nsqproxy/internal/module/tool"
	"github.com/changba/nsqproxy/internal/proxy"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 主函数
func main() {
	//初始化系统配置
	config.NewSystemConfig()
	model.NewDB(config.SystemConfig.DbHost, config.SystemConfig.DbPort, config.SystemConfig.DbUsername, config.SystemConfig.DbPassword, config.SystemConfig.DbName)
	//创建一个proxy实例
	p := proxy.NewProxy()
	//异常捕获
	defer func() {
		tool.PanicHandlerForLog()
		logger.Fatalf("nsqproxy will exit")
		os.Exit(2)
	}()
	//开启HTTP
	httper.NewHttper(config.SystemConfig.HttpAddr).Run()
	//灾备
	backup.Backup(config.SystemConfig.MasterAddr)
	//启动一个proxy实例
	logger.Infof("nsqproxy is starting")
	p.Run()
	//监听信号
	listenSignal(p)
	logger.Infof("nsqproxy end success")
}

// 监听信号
func listenSignal(p *proxy.Proxy) {
	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTRAP)
	for {
		sig := <-sigChannel
		logger.Infof("nsqproxy receive signal: %s", sig.String())
		if sig == syscall.SIGTRAP {
			continue
		}
		logger.Infof("nsqproxy is closing consumes...")
		p.SetExitFlag()
		p.Stop()
		//等待10秒
		logger.Infof("nsqproxy will be closed master process ten seconds later.")

		time.Sleep(10)
		//time.Sleep(10 * time.Second)
		break
	}
}
