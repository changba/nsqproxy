package worker

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/nsqio/go-nsq"
	"net"
	"strconv"
	"time"
)

// CBNSQ协议是自定义的非常简单的一个基于文本的TCP的协议。
// 就是消息长度 + 消息内容。 msg的长度(8字节，0填充） + messageId（16字节） + msg正文
// 如：有个用户注册的消息是：{"topic":"test","classname":"UserService","methodname":"addUser","param":["userid", "username", "password"],"addtime":"2020-11-27 14:30:34"}
// 第一部分：消息长度。消息主体是个json，长度140，补齐8位，那么这部分为00000140。8位的极限是99999999/1024/1024≈95M，足够了吧。
// 第二部分：消息ID，这个是NSQ的消息唯一ID，16位。如qwertyuiopasdfgh。注意：这16位是不计入第一部分的消息长度的。
// 第三部分：消息主体。即刚才提到的JSON串。
// 完整的消息为：00000140qwertyuiopasdfgh{"topic":"test","classname":"UserService","methodname":"addUser","param":["userid", "username", "password"],"addtime":"2020-11-27 14:30:34"}

type CBNSQWorker struct {
	workerConfig workerConfig
}

func (w *CBNSQWorker) new(wc workerConfig) {
	w.workerConfig = wc
}

func (w *CBNSQWorker) Send(message *nsq.Message) ([]byte, error) {
	//连接到worker
	conn, err := net.DialTimeout("tcp", w.workerConfig.addr, w.workerConfig.timeoutDial)
	if err != nil {
		return nil, newWorkerErrorConnect(err)
	}
	//设置连接的读写超时时间
	_ = conn.SetWriteDeadline(time.Now().Add(w.workerConfig.timeoutWrite))
	_ = conn.SetReadDeadline(time.Now().Add(w.workerConfig.timeoutRead))
	defer conn.Close()
	//给worker发送数据
	data := w.encode(message)
	n, err := conn.Write(data)
	if n == 0 {
		return nil, newWorkerErrorWrite(errors.New("n of conn.Write is 0"))
	}
	if err != nil {
		return nil, newWorkerErrorWrite(errors.New("conn.Write" + err.Error()))
	}
	//从worker读取响应
	buf := make([]byte, 128)
	n, err = conn.Read(buf)
	if err != nil {
		return nil, newWorkerErrorRead(errors.New("conn.Read" + err.Error()))
	}
	if n == 0 {
		return nil, newWorkerErrorRead(errors.New("response length is 0"))
	}
	return buf[:n], nil
}

// CBNSQ协议打包数据
// msg的长度(8字节，0填充） + messageId（16字节） + msg正文
func (w *CBNSQWorker) encode(message *nsq.Message) []byte {
	header := fmt.Sprintf("%08s", strconv.Itoa(len(message.Body)))
	return bytes.Join([]([]byte){[]byte(header), message.ID[:], message.Body}, []byte(""))
}
