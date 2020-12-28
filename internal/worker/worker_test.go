package worker

import (
	"errors"
	"github.com/nsqio/go-nsq"
	"testing"
	"time"
)

func TestNewWorker_IsError(t *testing.T) {
	if !IsErrorConnect(newWorkerErrorConnect(errors.New("hello world"))) {
		t.Fatal("IsErrorConnect failed")
	}
	if !IsErrorWrite(newWorkerErrorWrite(errors.New("hello world"))) {
		t.Fatal("IsErrorWrite failed")
	}
	if !IsErrorRead(newWorkerErrorRead(errors.New("hello world"))) {
		t.Fatal("IsErrorRead failed")
	}
}

func TestNewWorker(t *testing.T) {
	//构造消息
	messageId := nsq.MessageID([16]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p'})
	body := []byte("Hello world")
	message := nsq.NewMessage(messageId, body)

	//HTTP
	handler, err := NewWorker("127.0.0.1:80", "HtTp", "index.php", 1*time.Second, 1*time.Second, 1*time.Second)
	if err != nil {
		t.Fatalf("http NewWorker error: %s", err.Error())
	}
	res, err := handler.Send(message)
	if err != nil {
		t.Fatalf("http send error: %s", err.Error())
	}
	if string(res) != "200 ok" {
		t.Fatalf("http response body is not match")
	}

	//FastCGI
	handler, err = NewWorker("127.0.0.1:9000", "FaStCgI", "/var/www/index.php", 1*time.Second, 1*time.Second, 1*time.Second)
	if err != nil {
		t.Fatalf("fastcgi NewWorker error: %s", err.Error())
	}
	res, err = handler.Send(message)
	if err != nil {
		t.Fatalf("fastcgi send error: %s", err.Error())
	}
	if string(res) != "200 ok" {
		t.Fatalf("fastcgi response body is not match")
	}

	//CBNSQ
	handler, err = NewWorker("127.0.0.1:19910", "CbNsQ", "", 1*time.Second, 1*time.Second, 1*time.Second)
	if err != nil {
		t.Fatalf("cbnsq NewWorker error: %s", err.Error())
	}
	res, err = handler.Send(message)
	if err != nil {
		t.Fatalf("cbnsq send error: %s", err.Error())
	}
	if string(res) != "200 ok" {
		t.Fatalf("cbnsq response body is not match")
	}
}
