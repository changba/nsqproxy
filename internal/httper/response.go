package httper

import (
	"encoding/json"
	"github.com/ChangbaServer/nsqproxy/internal/module/logger"
	"net/http"
)

const HttpCodeOK = http.StatusOK
const HttpCodeBadRequest = http.StatusBadRequest
const HttpCodeHttpCodeForbidden = http.StatusUnauthorized
const HttpCodeForbidden = 403
const HttpCodeNotFound = 404
const HttpCodeInternalServerError = 500
const HttpCodeNotImplemented = 501
const HttpCodeBadGateway = 502
const HttpCodeServiceUnavailable = 503

type resp struct {
	w http.ResponseWriter
	Code int `json:"code"`
	Message string `json:"msg"`
	Result interface{} `json:"result"`
}

func Success(w http.ResponseWriter, result interface{}){
	response(w, HttpCodeOK, "ok", result)
}

func Failed(w http.ResponseWriter, code int, message string){
	response(w, code, message, nil)
}

func response(w http.ResponseWriter, code int, message string, result interface{}){
	r := resp{
		Code: code,
		Message: message,
		Result: result,
	}
	j, err := json.Marshal(r)
	if err != nil {
		logger.Errorf("response json error: %s", err.Error())
	}
	_, _ = w.Write(j)
}