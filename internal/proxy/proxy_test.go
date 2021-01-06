package proxy

import (
	"github.com/changba/nsqproxy/config"
	"github.com/changba/nsqproxy/internal/model"
	"testing"
	"time"
)

func init() {
	testing.Init()
	config.NewSystemConfig()
	model.NewDB("0.0.0.0", "3306", "root", "", "nsqproxy")
}

func TestNewProxy(t *testing.T) {
	p := NewProxy()
	if p.consumeConfigList[0].Id <= 0 {
		t.Fatal("初始化proxy-consumeConfig失败")
	}
	if p.IsStop() {
		t.Fatal("初始化proxy-exitFlag失败")
	}
}

func TestRun(t *testing.T) {
	p := NewProxy()
	p.Run()
	time.Sleep(2 * time.Second)

	for _, consumeConfig := range p.consumeConfigList {
		if !consumeConfig.StatusIsSuccess() {
			t.Fatal(consumeConfig.Topic + " consumeConfig.status is not success")
		}
		if consumeConfig.Consumer == nil {
			t.Fatal(consumeConfig.Topic + " consumeConfig.Consumer is nil")
		}
	}
	p.Stop()
	for _, consumeConfig := range p.consumeConfigList {
		if !consumeConfig.StatusIsClose() {
			t.Fatal(consumeConfig.Topic + " consume status is not closed")
		}
	}
	time.Sleep(2 * time.Second)
}
