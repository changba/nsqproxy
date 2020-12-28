package backup

import (
	"encoding/json"
	"github.com/ChangbaServer/nsqproxy/internal/module/logger"
	"io/ioutil"
	"net/http"
	"time"
)

//主库返回的响应值
type masterResponse struct {
	Code int `json:"code"`
	Message string `json:"msg"`
	Result interface{} `json:"result"`
}

//向主发送心跳，主观判断主下线后，备才会启动，否则在本函数阻塞
//成功：成功一次就算成功
//失败：连续失败五次算失败。
//即使发生分区、网络阻塞等原因主被错误的判断为失败，此时主备会同时在线，但我们认为，主备同时在线，并不会造成消息的丢失（但是极限下会加快两倍消费速度）。
//考虑引入多节点的哨兵判断，上升复杂度，不如就这样简单做。
func Backup(masterAddr string) {
	//主的空，则忽略
	if masterAddr == "" {
		return
	}
	//访问主机失败次数
	failedNum := 0
	successInterval := 15 * time.Second
	failedInterval := 2 * time.Second
	interval := successInterval
	for {
		resp, err := http.Get("http://" + masterAddr + "/status")
		var body []byte
		var r masterResponse
		if err != nil {
			interval = failedInterval
			failedNum++
			logger.Errorf("the backup connect master error: %s", err.Error())
			logger.Errorf("failed to access the master %d times", failedNum)
			goto SLEEP
		}
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil{
			interval = failedInterval
			failedNum++
			logger.Errorf("the backup read response of master error: %s", err.Error())
			logger.Errorf("failed to access the master %d times", failedNum)
			goto SLEEP
		}
		err = json.Unmarshal(body, &r)
		if err != nil{
			interval = failedInterval
			failedNum++
			logger.Errorf("the backup change response of master to json error: %s", err.Error())
			logger.Errorf("failed to access the master %d times", failedNum)
			goto SLEEP
		}
		_ = resp.Body.Close()
		if r.Code != 200 || r.Result != "ok"{
			interval = failedInterval
			failedNum++
			logger.Errorf("the backup read response of master failed. action:%d %s expect:%d %s", r.Code, r.Result, 200, "ok")
			logger.Errorf("failed to access the master %d times", failedNum)
			goto SLEEP
		}
		logger.Infof("the master is normal. the backup is keeping on block.")
		interval = successInterval
		failedNum = 0

SLEEP:
		if failedNum >= 5{
			logger.Errorf("the backup will run.")
			break
		}
		time.Sleep(interval)
	}
}
