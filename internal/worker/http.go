package worker

import (
	"errors"
	"github.com/ChangbaServer/nsqproxy/internal/module/tool"
	"github.com/nsqio/go-nsq"
	"io/ioutil"
	"net/http"
	"strings"
)

type HTTPWorker struct {
	workerConfig workerConfig
	clientPool *tool.HttpClientPool
}

func (w *HTTPWorker) new(wc workerConfig){
	w.workerConfig = wc
	w.clientPool = tool.NewHttpClientPool()
}

//给HTTP发消息
func (w *HTTPWorker) Send(message *nsq.Message) ([]byte, error) {
	//构造HTTP请求
	//values := url.Values{}
	//values.Set("param", string(message.Body))
	req, err := http.NewRequest("POST", "http://" + w.workerConfig.addr + "/" + w.workerConfig.extra, strings.NewReader(string(message.Body)))
	if err != nil {
		return nil, err
	}
	//含下划线会被nginx抛弃，横线会被转为下划线。
	req.Header.Set("MESSAGE_ID", string(message.ID[:]))
	req.Header.Set("MESSAGE-ID", string(message.ID[:]))
	req.Header.Set("CONTENT-TYPE", "application/x-www-form-urlencoded")
	//获取http.Client
	client := w.clientPool.GetClient()
	if client == nil{
		return nil, errors.New("HttpClientPool.GetClient is nil")
	}
	defer w.clientPool.PutClient(client)
	client.Timeout = w.workerConfig.timeoutDial
	//发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, newWorkerErrorWrite(err)
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, newWorkerErrorRead(err)
	}
	return content, nil
}