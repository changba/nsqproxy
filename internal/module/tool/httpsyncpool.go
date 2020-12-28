package tool

import (
	"errors"
	"net"
	"net/http"
	"sync"
	"time"
)

//用sync.pool管理http client
//用sync.Pool和不用，QPS提升10%，CPU下降50%
type HttpClientPool struct{
	pool sync.Pool
}

func NewHttpClientPool()*HttpClientPool {
	return &HttpClientPool{
		pool: sync.Pool{
			New: func() interface{} {
				return NewHttpClient()
			},
		},
	}
}

func (h *HttpClientPool)GetClient()*http.Client{
	return h.pool.Get().(*http.Client)
}

func (h *HttpClientPool)PutClient(client *http.Client){
	h.pool.Put(client)
}

func(h *HttpClientPool)Dial(req *http.Request)(*http.Response, error){
	client := h.GetClient()
	if client == nil{
		return nil, errors.New("HttpClientPool.GetClient is nil")
	}
	defer h.PutClient(client)
	return client.Do(req)
}

//创建一个http client
func NewHttpClient()*http.Client{
	//参数是复制的DefaultClient，只是改了MaxIdleConns和MaxIdleConnsPerHost
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          500,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   500,
		},
	}
}