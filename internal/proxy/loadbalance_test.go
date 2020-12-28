package proxy

import (
	"github.com/ChangbaServer/nsqproxy/internal/model"
	"testing"
)

func init(){
	model.NewDB("0.0.0.0", "3306", "root", "", "nsqproxy")
}

func TestNewLoadBalance(t *testing.T) {
	p := NewProxy()
	if len(p.consumeConfigList) <= 0 || len(p.consumeConfigList[0].ServerList) <= 0 {
		t.Fatal("consume workerlist为空")
	}
	loadbalance := newLoadBalance(LoadBalanceMethodLoop, p.consumeConfigList[0].ServerList)
	if loadbalance.method != LoadBalanceMethodLoop || loadbalance.dispatcher == nil {
		t.Fatal("newLoadBalance failed")
	}
}

func TestLoadBalance_pickWorker(t *testing.T) {
	p := NewProxy()
	loadbalance := newLoadBalance(LoadBalanceMethodLoop, p.consumeConfigList[0].ServerList)
	work, err := loadbalance.pickWorker()
	if err != nil {
		t.Fatalf("pickWorker error: %s", err.Error())
	}
	if work.Id <= 0 || work.Weight <= 0 || len(work.WorkServer.Addr) <= 0 || len(work.WorkServer.Protocol) <= 0 {
		t.Fatalf("pickWorker error: no available work")
	}
}

func TestLoadBalanceLoop_new(t *testing.T) {
	p := NewProxy()
	loop := &loadBalanceLoop{}
	loop.new(p.consumeConfigList[0].ServerList)
	if len(loop.list) <= 0{
		t.Fatalf("loop.list is empty")
	}
	//检验capacity是否相等
	capacity := 0
	//检验生成的list中，每台机器的占比是否相等
	list := make(map[int]int, 0)
	for _, v := range p.consumeConfigList[0].ServerList {
		capacity += v.Weight * 100
		list[v.Id] = v.Weight * 100
	}
	if loop.capacity != capacity{
		t.Fatalf("loop.capacity is not match. action:%d expect:%d", loop.capacity, capacity)
	}
	count := make(map[int]int, 0)
	for _, v := range loop.list {
		count[v.Id] = v.Weight * 100
	}
	if len(count) != len(list){
		t.Fatalf("loop.list length does not match. action:%d expect:%d", len(count), len(list))
	}
	for k, v := range list {
		if _, ok := count[k]; !ok {
			t.Fatalf("the id %d of loop.list not exists", k)
		}
		if v != count[k] {
			t.Fatalf("the id %d total of loop.list does not match. action:%d expect:%d", k, count[k], v)
		}
	}
}
