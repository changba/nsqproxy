package httper

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestHttper_Run(t *testing.T) {
	addr := "0.0.0.0:19421"
	httper := NewHttper(addr)
	httper.Run()
	time.Sleep(100 * time.Microsecond)
	url := "http://" + addr + "/status"
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal("请求错误: " + err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ReadAll error: %s", err.Error())
	}
	_ = resp.Body.Close()

	type response struct {
		Code    int         `json:"code"`
		Message string      `json:"msg"`
		Result  interface{} `json:"result"`
	}
	r := response{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		t.Fatalf("json decode failed. err: %s", err.Error())
	}
	if r.Code != 200 || r.Result != "ok" {
		t.Fatalf("response failed. action code:%d expect code:%d action result:%s expect result:%s", r.Code, 200, body, "ok")
	}
}
