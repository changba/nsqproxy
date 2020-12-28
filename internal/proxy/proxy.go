package proxy

import (
	"errors"
	"fmt"
	"github.com/ChangbaServer/nsqproxy/config"
	"github.com/ChangbaServer/nsqproxy/internal/model"
	"github.com/ChangbaServer/nsqproxy/internal/module/logger"
	"github.com/ChangbaServer/nsqproxy/internal/module/tool"
	"github.com/nsqio/go-nsq"
	"strconv"
	"time"
)

const proxyStatusNormal = 0
const proxyStatusExit = 1

type Proxy struct {
	//消费者配置列表
	consumeConfigList []model.ConsumeConfig
	//是否退出
	exitFlag int
	//角色
	role int
}

// 生成一个新的Proxy实例
// 此方法建议全局只有一次调用
func NewProxy() *Proxy {
	//配置
	consumeConfigList, err := model.GetAvailableConsumeList()
	if err != nil {
		panic("load config file failed. err: " + err.Error())
	}
	if len(consumeConfigList) == 0 {
		logger.Warningf("consume config is empty.")
	}
	return &Proxy{
		consumeConfigList: consumeConfigList,
		exitFlag:          proxyStatusNormal,
	}
}

//开始运行NSQProxy
func (p *Proxy) Run() {
	config.SystemConfig.Role = config.RoleMaster
	//启动消费者
	go p.run()
}

// 立刻启动消费者，随后启动定时器，定时获取最新的系统配置和消费者配置
// 配置未发生变化时，忽略
// 配置发生变化时，停止旧消费者，启动新消费者
func (p *Proxy) run() {
	//立即启动Proxy
	p.start()
	//定时器。指定时间运行一次
	if config.SystemConfig.UpdateConfigInterval <= 0 {
		return
	}
	ticker := time.NewTicker(config.SystemConfig.UpdateConfigInterval)
	for {
		select {
		case <-ticker.C:
			if p.IsStop() {
				logger.Warningf("nsqproxy update consume config quit")
				ticker.Stop()
				break
			}
			// 更新消费者配置
			p.updateConsumeConfigList()
			//更新work机状态统计
			logger.Infof("nsqproxy update consume config success.")
			p.start()
		}
	}
}

//立即运行所有的消费者实例
func (p *Proxy) start() bool {
	//状态初始化
	for k, consumeConfig := range p.consumeConfigList {
		//启动状态监测
		if consumeConfig.StatusIsSuccess() {
			continue
		}
		p.consumeConfigList[k].SetStatusWait()
		_, err := p.startConsume(&(p.consumeConfigList[k]))
		if err != nil {
			p.consumeConfigList[k].SetStatusFailed()
			p.consumeConfigList[k].Consumer = nil
		} else {
			p.consumeConfigList[k].SetStatusSuccess()
		}
	}
	return true
}

//开始订阅一个Topic
func (p *Proxy) startConsume(consumeConfig *model.ConsumeConfig) (*nsq.Consumer, error) {
	//go-nsq配置
	clientConfig := nsq.NewConfig()
	clientConfig.MaxInFlight = consumeConfig.MaxInFlight
	consumer, err := nsq.NewConsumer(consumeConfig.Topic, consumeConfig.Channel, clientConfig)
	if err != nil {
		errMsg := fmt.Sprintf("start consume failed. id: %s, topic: %s, channel：%s, error: %s", strconv.Itoa(consumeConfig.Id), consumeConfig.Topic, consumeConfig.Channel, err.Error())
		logger.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}
	//把本系统的logLevel映射成go-nsq项目的等级
	var nsqLogLevel nsq.LogLevel
	switch config.SystemConfig.SubLogger.Level {
	case logger.LOG_DEBUG:
		nsqLogLevel = nsq.LogLevelDebug
	case logger.LOG_INFO:
		nsqLogLevel = nsq.LogLevelInfo
	case logger.LOG_WARNING:
		nsqLogLevel = nsq.LogLevelWarning
	case logger.LOG_ERROR:
	case logger.LOG_FATAL:
		nsqLogLevel = nsq.LogLevelError
	}
	//consumer.SetLogger(config.SystemConfig.SubLogger, nsqLogLevel)
	consumer.SetLoggerLevel(nsqLogLevel)
	consumeConfig.Consumer = consumer
	consumerUniqueId := tool.GenerateUniqueId(int64(consumeConfig.Id))
	consumeConfig.SetConsumerUniqueId(consumerUniqueId)
	handler := NewHandler(*consumeConfig)
	consumer.AddConcurrentHandlers(&handler, consumeConfig.HandleNum)
	err = consumer.ConnectToNSQLookupds(config.SystemConfig.NsqlookupdHttpAddrList)
	if err != nil {
		errMsg := "connect nsqlookupd failed. consumeruniqueid: " + string(consumerUniqueId[:]) + ", id: " + strconv.Itoa(consumeConfig.Id) + ", topic: " + consumeConfig.Topic + ", channel：" + consumeConfig.Channel + ", error: " + err.Error()
		logger.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}
	logger.Infof("run consume success. consumeruniqueid: %s, id: %d, topic: %s, channel：%s", consumerUniqueId, consumeConfig.Id, consumeConfig.Topic, consumeConfig.Channel)
	return consumer, nil
}

// 更新消费者配置
func (p *Proxy) updateConsumeConfigList() bool {
	//获取最新的配置
	consumeConfigList, err := model.GetAvailableConsumeList()
	if err != nil {
		logger.Errorf("nsqproxy get consume config failed. error: " + err.Error())
		return false
	}
	// 更新消费者配置
	//等待结束的
	stopList := make([]model.ConsumeConfig, 0)
	// 1、循环新配置
	// 新有的，旧也有的，判断配置是否相同，相同的忽略，不同的结束
	// 新有的，旧没有的，等待启动
	oldConsumeConfigMap := model.FormatConsumeConfigListForMap(p.consumeConfigList)
	for newIndex, newC := range consumeConfigList {
		if oldC, ok := oldConsumeConfigMap[newC.Id]; ok {
			//新配置和旧配置一样，则忽略
			if newC.IsEqual(oldC) {
				consumeConfigList[newIndex] = oldC
				continue
			}
			//新配置和旧配置不一样，则停止旧配置产生的消费者
			logger.Infof("nsqproxy update consume config: new!=old, stop the old. consumeruniqueid: %s, id: %d, topic: %s, channel：%s", oldC.GetConsumerUniqueId(), oldC.Id, oldC.Topic, oldC.Channel)
			stopList = append(stopList, oldC)
		}
	}
	// 2、循环旧配置
	// 旧有的，新没有的，结束
	newConsumeConfigMap := model.FormatConsumeConfigListForMap(consumeConfigList)
	for _, oldC := range p.consumeConfigList {
		if _, ok := newConsumeConfigMap[oldC.Id]; !ok {
			//旧配置不在新配置，结束旧配置
			logger.Infof("nsqproxy update consume config: old not exists in new, stop old. consumeruniqueid: %s, id: %d, topic: %s, channel：%s", oldC.GetConsumerUniqueId(), oldC.Id, oldC.Topic, oldC.Channel)
			stopList = append(stopList, oldC)
		}
	}
	// 3、新配置覆盖旧配置
	p.consumeConfigList = consumeConfigList
	// 4、停止配置发生变化的旧消费者
	for k, _ := range stopList {
		p.StopConsume(&(stopList[k]))
	}
	return true
}

func (p *Proxy) Stop() {
	//关闭所有的消费者
	logger.Infof("nsqproxy begin stopping all consume.")
	for k, _ := range p.consumeConfigList {
		p.StopConsume(&(p.consumeConfigList[k]))
	}
	logger.Infof("nsqproxy stop all consume success.")
}

//结束一个订阅
func (p *Proxy) StopConsume(consumeConfig *model.ConsumeConfig) {
	if consumeConfig.Id > 0 && consumeConfig.Consumer != nil {
		consumeConfig.Consumer.ChangeMaxInFlight(0)
		consumeConfig.Consumer.Stop()
		consumeConfig.SetStatusClosed()
		logger.Infof("stop consume success. consumeruniqueid: %s, id: %d, topic: %s, channel：%s", consumeConfig.GetConsumerUniqueId(), consumeConfig.Id, consumeConfig.Topic, consumeConfig.Channel)
	}
}

func (p *Proxy) SetExitFlag() {
	p.exitFlag = proxyStatusExit
}

func (p *Proxy) IsStop() bool {
	return p.exitFlag == proxyStatusExit
}

func (p *Proxy) GetStop() int {
	return p.exitFlag
}
