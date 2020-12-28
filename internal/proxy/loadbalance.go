package proxy

import (
	"errors"
	"github.com/ChangbaServer/nsqproxy/internal/model"
)

const LoadBalanceMethodLoop = "loop"

//负载均衡分发器
//目前只支持轮询
type loadBalanceDispatcher interface {
	//初始化方法
	new([]model.ConsumeServerMap)
	//选择
	pick() (model.ConsumeServerMap, error)
}

type loadBalance struct {
	method     string                //负载均衡的方法
	dispatcher loadBalanceDispatcher //负载均衡分发器，具体的算法执行 者
}

func newLoadBalance(method string, consumeServerMapList []model.ConsumeServerMap) *loadBalance {
	//过滤一下，选出可用的机器
	list := make([]model.ConsumeServerMap, 0)
	for _, v := range consumeServerMapList {
		if v.Id <= 0 || v.Weight <= 0 || len(v.WorkServer.Addr) <= 0 || len(v.WorkServer.Protocol) <= 0 {
			continue
		}
		list = append(list, v)
	}

	//构造分发器，目前只支持轮询
	//轮询
	var dispatcher loadBalanceDispatcher
	if method == LoadBalanceMethodLoop {
		dispatcher = &loadBalanceLoop{}
	} else {
		method = LoadBalanceMethodLoop
		dispatcher = &loadBalanceLoop{}
	}
	dispatcher.new(list)

	return &loadBalance{
		method:     method,
		dispatcher: dispatcher,
	}
}

//选择一个work机器 - 轮询的方式根据权重选择
func (l *loadBalance) pickWorker() (model.ConsumeServerMap, error) {
	work, err := l.dispatcher.pick()
	//失败了就再取一次
	if err != nil || work.Id <= 0 || work.Weight <= 0 || len(work.WorkServer.Addr) <= 0 || len(work.WorkServer.Protocol) <= 0 {
		work, err = l.dispatcher.pick()
		if err != nil || work.Id <= 0 || work.Weight <= 0 || len(work.WorkServer.Addr) <= 0 || len(work.WorkServer.Protocol) <= 0 {
			return work, errors.New("no available work")
		}
	}
	return work, nil
}

type loadBalanceLoop struct {
	list     []model.ConsumeServerMap //机器列表
	capacity int                      //权重总分
	seek     int
}

func (l *loadBalanceLoop) new(consumeServerMapList []model.ConsumeServerMap) {
	//所有机器的权重和
	capacity := 0
	for _, v := range consumeServerMapList {
		capacity += v.Weight
	}
	//生成Worker列表，按顺序使用，比rand节省CPU
	multiple := 100
	capacity *= multiple
	wl := make([]model.ConsumeServerMap, capacity)
	//填充
	index := 0
	for i := 0; i < multiple; i++ {
		for _, consumeServerMap := range consumeServerMapList {
			for j := 0; j < consumeServerMap.Weight; j++ {
				wl[index] = consumeServerMap
				index++
			}
		}
	}
	l.list = wl
	l.capacity = capacity
	l.seek = 0
}

//选择一个work机器 - 轮询的方式根据权重选择
func (l *loadBalanceLoop) pick() (model.ConsumeServerMap, error) {
	//这里无需加锁，并发问题最多导致分布不是绝对的均衡。你多执行一个我少执行一个无伤大雅。
	seek := l.seek
	if seek >= l.capacity {
		l.seek = 0
		seek = 0
	}
	work := l.list[seek]
	//失败了就取索引为的0，补取一次
	if work.Id <= 0 || work.Weight <= 0 || len(work.WorkServer.Addr) <= 0 || len(work.WorkServer.Protocol) <= 0 {
		l.seek = 0
		seek = 0
		work = l.list[seek]
		if work.Id <= 0 || work.Weight <= 0 || len(work.WorkServer.Addr) <= 0 || len(work.WorkServer.Protocol) <= 0 {
			return work, errors.New("no available work")
		}
	}
	l.seek++
	return work, nil
}
