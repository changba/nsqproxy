package worker

import (
	"github.com/nsqio/go-nsq"
	"testing"
	"time"
)

func TestCBNSQWorker_Send(t *testing.T) {
	wc := newWorkerConfig("127.0.0.1:19910", "CbNsQ", "", 1*time.Second, 1*time.Second, 1*time.Second)
	handler := &CBNSQWorker{}
	handler.new(wc)

	messageId := nsq.MessageID([16]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p'})
	body := []byte("Hello world")
	message := nsq.NewMessage(messageId, body)
	res, err := handler.Send(message)
	if err != nil{
		t.Fatalf("send error: %s", err.Error())
	}
	if string(res) != "200 ok"{
		t.Fatalf("response body is not match")
	}
}