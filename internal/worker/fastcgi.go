package worker

import (
	"bytes"
	"github.com/changba/nsqproxy/config"
	"github.com/changba/nsqproxy/internal/module/fastcgi"
	"github.com/nsqio/go-nsq"
	"io/ioutil"
)

type FastCGIWorker struct {
	workerConfig workerConfig
}

func (w *FastCGIWorker) new(workerConfig workerConfig) {
	w.workerConfig = workerConfig
}

func (w *FastCGIWorker) Send(message *nsq.Message) ([]byte, error) {
	//连接到work
	fc, err := fastcgi.DialTimeout("tcp", w.workerConfig.addr, w.workerConfig.timeoutDial, w.workerConfig.timeoutWrite, w.workerConfig.timeoutRead)
	if err != nil {
		return nil, newWorkerErrorConnect(err)
	}
	defer fc.Close()
	//fpm参数，php可以用$_SERVER获取
	server := w.getServer(w.workerConfig, message)
	//给work发送数据
	rd := bytes.NewReader(message.Body)
	resp, err := fc.Post(server, "", rd, rd.Len())
	if err != nil {
		return nil, newWorkerErrorWrite(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, newWorkerErrorRead(err)
	}
	return content, nil
}

func (w *FastCGIWorker) getServer(wc workerConfig, message *nsq.Message) map[string]string {
	//fpm参数，php可以用$_SERVER获取
	server := make(map[string]string)
	server["SCRIPT_FILENAME"] = wc.extra
	server["REMOTE_ADDR"] = config.SystemConfig.InternalIP
	server["MESSAGE_ID"] = string(message.ID[:])
	return server
}
