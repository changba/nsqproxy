package config

import (
	"flag"
	"github.com/changba-server/nsqproxy/internal/module/logger"
	"github.com/changba-server/nsqproxy/internal/module/tool"
	"strings"
	"time"
)

const NsqproxyVersion = "1.0.0"

const RoleMaster = "master"
const RoleBackup = "backup"

var SystemConfig = &systemConfig{}

type systemConfig struct {
	//-----本项目相关-----
	//本机监听的端口
	HttpAddr string
	//主库IP端口，本机为主库时请留空。
	MasterAddr string

	//-----NSQ相关-----
	NsqlookupdHttpAddrList []string

	//-----日志相关-----
	//订阅的消息日志logger
	SubLogger *logger.Logger

	//-----消费者相关-----
	UpdateConfigInterval time.Duration

	//-----数据库相关-----
	DbHost     string
	DbPort     string
	DbUsername string
	DbPassword string
	DbName     string

	//-----本机相关-----
	//本机IP
	InternalIP string

	//全局配置，无需自定义
	Role string
}

func NewSystemConfig() {
	//参数
	var httpAddr = flag.String("httpAddr", "0.0.0.0:19421", "监听的HTTP端口")
	var masterAddr = flag.String("masterAddr", "", "主库IP端口，为空则本机为主机")
	//nsq相关
	var nsqlookupdHTTP = flag.String("nsqlookupdHTTP", "127.0.0.1:4161", "nsqLookupd的HTTP地址，多个用逗号分割如'127.0.0.1:4161,127.0.0.1:4163'")
	//log相关
	var logLevel = flag.String("logLevel", "info", "日志等级，可选有debug、info、warning、error、fatal")
	var logPath = flag.String("logPath", "logs/proxy.log", "系统日志路径")
	var subLogPath = flag.String("subLogPath", "logs/sub.log", "消费log，由于量大成功消费log仅在日志等级为debug时启用")
	//MySQL
	var dbHost = flag.String("dbHost", "127.0.0.1", "MySQL的IP")
	var dbPort = flag.String("dbPort", "3306", "MySQL的端口")
	var dbUsername = flag.String("dbUsername", "root", "MySQL的账号")
	var dbPassword = flag.String("dbPassword", "", "MySQL的密码")
	var dbName = flag.String("dbName", "nsqproxy", "MySQL的库名")
	//消费者相关
	var updateConfigInterval = flag.Int64("updateConfigInterval", 60, "更新配置间隔")

	flag.Parse()
	//本机相关
	SystemConfig.HttpAddr = *httpAddr
	if len(SystemConfig.HttpAddr) <= 0 {
		panic("httpAddr参数缺失")
	}
	SystemConfig.MasterAddr = *masterAddr
	//nsqlookupd相关
	SystemConfig.NsqlookupdHttpAddrList = strings.Split(*nsqlookupdHTTP, ",")
	if len(SystemConfig.NsqlookupdHttpAddrList) <= 0 {
		panic("nsqlookupdHTTP 缺失")
	}
	//日志相关
	logLevelLower := strings.ToLower(*logLevel)
	logLevelList := map[string]struct{}{"debug": struct{}{}, "info": struct{}{}, "warning": struct{}{}, "error": struct{}{}, "fatal": struct{}{}}
	if _, ok := logLevelList[logLevelLower]; !ok {
		panic("logLevel可选值为debug、info、warning、error、fatal")
	}
	logger.Init(*logPath, logLevelLower)
	SystemConfig.SubLogger = logger.NewLogger(*subLogPath, "", logLevelLower)
	//数据库
	SystemConfig.DbHost = *dbHost
	SystemConfig.DbPort = *dbPort
	SystemConfig.DbUsername = *dbUsername
	SystemConfig.DbPassword = *dbPassword
	SystemConfig.DbName = *dbName
	//消费者相关
	SystemConfig.UpdateConfigInterval = time.Duration(*updateConfigInterval) * time.Second
	//版本
	SystemConfig.InternalIP = tool.GetInternalIP()

	//全局配置，无需自定义
	SystemConfig.Role = RoleBackup
}

func (s *systemConfig) Close() bool {
	logger.Close()
	s.SubLogger.Close()
	s = nil
	return true
}
