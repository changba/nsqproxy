package proxy

import (
	"github.com/changba-server/nsqproxy/internal/model"
	"testing"
)

func init() {
	model.NewDB("0.0.0.0", "3306", "root", "", "nsqproxy")
}

func TestNewHandler(t *testing.T) {
	p := NewProxy()
	consumeConfig := p.consumeConfigList[0]
	NewHandler(consumeConfig)
}
