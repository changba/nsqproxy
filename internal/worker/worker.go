package worker

import (
	"errors"
	"github.com/nsqio/go-nsq"
	"strings"
	"time"
)

//和Worker的通信协议
//HTTP协议
const protocolHttp = "http"

//FastCGI协议
const protocolFastCGI = "fastcgi"

//自定义CBNSQ协议（唱吧NSQ）
const protocolCBNSQ = "cbnsq"

type workerConfig struct {
	addr         string
	protocol     string //小写
	extra        string
	timeoutDial  time.Duration
	timeoutWrite time.Duration
	timeoutRead  time.Duration
}

func newWorkerConfig(addr, protocol, extra string, timeoutDial, timeoutWrite, timeoutRead time.Duration) workerConfig {
	return workerConfig{
		addr:         addr,
		protocol:     strings.ToLower(protocol),
		extra:        extra,
		timeoutDial:  timeoutDial,
		timeoutWrite: timeoutWrite,
		timeoutRead:  timeoutRead,
	}
}

//worker接口
type Worker interface {
	new(workerConfig)
	Send(*nsq.Message) ([]byte, error)
}

func NewWorker(addr, protocol, extra string, timeoutDial, timeoutWrite, timeoutRead time.Duration) (Worker, error) {
	wc := newWorkerConfig(addr, protocol, extra, timeoutDial, timeoutWrite, timeoutRead)
	var handler Worker
	switch wc.protocol {
	case protocolHttp:
		handler = &HTTPWorker{}
	case protocolFastCGI:
		handler = &FastCGIWorker{}
	case protocolCBNSQ:
		handler = &CBNSQWorker{}
	default:
		return nil, errors.New("worker invalid protocol")
	}
	handler.new(wc)
	return handler, nil
}
