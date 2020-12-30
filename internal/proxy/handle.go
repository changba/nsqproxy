package proxy

import (
	"errors"
	"github.com/changba-server/nsqproxy/config"
	"github.com/changba-server/nsqproxy/internal/model"
	"github.com/changba-server/nsqproxy/internal/module/logger"
	"github.com/changba-server/nsqproxy/internal/worker"
	"github.com/nsqio/go-nsq"
	"strconv"
)

//响应成功
const WorkResponseSuccess = "200"

type Handler struct {
	consumeConfig model.ConsumeConfig
	workerList    *loadBalance
}

func NewHandler(consumeConfig model.ConsumeConfig) Handler {
	return Handler{
		consumeConfig: consumeConfig,
		workerList:    newLoadBalance(LoadBalanceMethodLoop, consumeConfig.ServerList),
	}
}

//处理收到的订阅的消息
//返回error就会重新入队，返回nil则不会
func (h *Handler) HandleMessage(message *nsq.Message) error {
	//根据权重选一个work
	workServer, err := h.workerList.pickWorker()
	if err != nil {
		h.recordSubLog(logger.LOG_ERROR, message.ID, message.Body, workServer.WorkServer.Addr, "pick work error. "+err.Error())
		return err
	}
	h.recordSubLog(logger.LOG_DEBUG, message.ID, message.Body, workServer.WorkServer.Addr, "")
	//初始化Worker
	w, err := worker.NewWorker(workServer.WorkServer.Addr, workServer.WorkServer.Protocol, workServer.WorkServer.Extra, h.consumeConfig.TimeoutDial, h.consumeConfig.TimeoutWrite, h.consumeConfig.TimeoutRead)
	if err != nil {
		h.recordSubLog(logger.LOG_ERROR, message.ID, message.Body, workServer.WorkServer.Addr, "new worker error. "+err.Error())
		return err
	}
	//向Worker发送数据
	response, err := w.Send(message)
	if err != nil {
		h.recordSubLog(logger.LOG_ERROR, message.ID, message.Body, workServer.WorkServer.Addr, "send worker error. "+err.Error())
		if h.isRequeueByError(err) {
			return err
		}
	}
	//返回值长度判断
	if len(response) < 3 {
		//约定返回值前3位为状态码，如：200 success
		errMsg := "response length less 3. response[" + strconv.Itoa(len(response)) + "]: " + string(response)
		h.recordSubLog(logger.LOG_WARNING, message.ID, message.Body, workServer.WorkServer.Addr, errMsg)
		//是否需要重新入队
		if h.consumeConfig.IsRequeue {
			return errors.New(errMsg)
		} else {
			return nil
		}
	}
	codeString := string(response[0:3])
	//返回值判断
	if codeString == WorkResponseSuccess {
		return nil
	}
	errMsg := "response code error. code: " + codeString + ", content: " + string(response)
	h.recordSubLog(logger.LOG_WARNING, message.ID, message.Body, workServer.WorkServer.Addr, errMsg)
	//是否需要重新入队
	if !h.consumeConfig.IsRequeue {
		return nil
	}
	return errors.New(errMsg)
}

//是否重新入队，根据error类型和配置决定
//连不上worker直接重新入队，write失败直接重新入队，read失败根据配置决定是否重新入队
func (h *Handler) isRequeueByError(err error) bool {
	if worker.IsErrorConnect(err) {
		return true
	} else if worker.IsErrorWrite(err) {
		return true
	} else if worker.IsErrorRead(err) && h.consumeConfig.IsRequeue {
		return true
	}
	return false
}

// 处理失败时会调用此方法
// 连续失败五次次，第六次时在go-nsq客户端giving up时会执行
func (h *Handler) LogFailedMessage(message nsq.Message) {
	h.recordSubLog(logger.LOG_FATAL, message.ID, message.Body, "", "LogFailedMessage giving up")
}

//封装写log方法，上面代码太啰嗦
func (h *Handler) recordSubLog(logLevel logger.LogLevel, messageId nsq.MessageID, messageBody []byte, workAddr, errMsg string) {
	if len(workAddr) <= 0 {
		workAddr = "null"
	}
	if len(errMsg) <= 0 {
		errMsg = "nil"
	}
	config.SystemConfig.SubLogger.WithLevelf(logLevel, "[%s/%s] nsqproxy %s %s messageid:%s messagebody:%s",
		h.consumeConfig.Topic, h.consumeConfig.Channel,
		workAddr, errMsg, string(messageId[:]), string(messageBody),
	)
}
