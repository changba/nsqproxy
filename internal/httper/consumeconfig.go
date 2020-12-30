package httper

import (
	"github.com/changba-server/nsqproxy/internal/model"
	"net/http"
	"strconv"
	"time"
)

type ConsumeConfig struct {
}

func (ConsumeConfig) Create(w http.ResponseWriter, r *http.Request) {
	c := &model.ConsumeConfig{}
	c.Topic = r.FormValue("topic")
	c.Channel = r.FormValue("channel")
	//描述
	c.Description = r.FormValue("description")
	//责任人
	c.Owner = r.FormValue("owner")
	//积压报警阈值
	monitorThreshold, err := strconv.Atoi(r.FormValue("monitorThreshold"))
	if err == nil && monitorThreshold >= 0 {
		c.MonitorThreshold = monitorThreshold
	}
	//该队列的并发量
	handleNum, err := strconv.Atoi(r.FormValue("handleNum"))
	if err == nil && handleNum >= 0 {
		c.HandleNum = handleNum
	}
	//NSQD最多同时推送多少个消息
	maxInFlight, err := strconv.Atoi(r.FormValue("maxInFlight"))
	if err == nil && maxInFlight >= 0 {
		c.MaxInFlight = maxInFlight
	}
	//失败，超时等情况是否重新入队
	isRequeue, err := strconv.ParseBool(r.FormValue("isRequeue"))
	if err == nil && isRequeue {
		c.IsRequeue = isRequeue
	}
	//超时时间
	timeoutDial, err := strconv.ParseInt(r.FormValue("timeoutDial"), 10, 64)
	if err == nil && timeoutDial >= 0 {
		c.TimeoutDial = time.Duration(timeoutDial)
	}
	//读超时时间
	timeoutRead, err := strconv.ParseInt(r.FormValue("timeoutRead"), 10, 64)
	if err == nil && timeoutRead >= 0 {
		c.TimeoutRead = time.Duration(timeoutRead)
	}
	//写超时时间
	timeoutWrite, err := strconv.ParseInt(r.FormValue("timeoutWrite"), 10, 64)
	if err == nil && timeoutWrite >= 0 {
		c.TimeoutWrite = time.Duration(timeoutWrite)
	}
	//是否有效
	invalid, err := strconv.Atoi(r.FormValue("invalid"))
	if err != nil || invalid <= 0 || (invalid != model.InvalidAvailable && invalid != model.InvalidUnavailable) {
		invalid = model.InvalidAvailable
	}
	c.Invalid = invalid

	id, err := c.Create()
	if err != nil {
		Failed(w, HttpCodeBadRequest, "create failed. err: "+err.Error())
		return
	}
	if id <= 0 {
		Failed(w, HttpCodeBadRequest, "id is zero")
		return
	}
	Success(w, c)
}

func (ConsumeConfig) Delete(w http.ResponseWriter, r *http.Request) {
	c := &model.ConsumeConfig{}
	var err error
	c.Id, err = strconv.Atoi(r.FormValue("id"))
	if err != nil || c.Id <= 0 {
		Failed(w, HttpCodeBadRequest, "param error")
		return
	}
	m, err := (&model.ConsumeServerMap{}).AllByConsumeid(c.Id)
	if err != nil {
		Failed(w, HttpCodeBadRequest, "get consumeServerMap failed. err: "+err.Error())
		return
	}
	if len(m) > 0 {
		Failed(w, HttpCodeBadRequest, "Cannot delete because there is an associated worker server")
		return
	}
	_, err = c.Delete()
	if err != nil {
		Failed(w, HttpCodeBadRequest, "delete failed. err: "+err.Error())
		return
	}
	Success(w, "ok")
}

func (ConsumeConfig) Update(w http.ResponseWriter, r *http.Request) {
	c := &model.ConsumeConfig{}
	var err error
	c.Id, err = strconv.Atoi(r.FormValue("id"))
	if err != nil || c.Id <= 0 {
		Failed(w, HttpCodeBadRequest, "param error")
		return
	}
	c.Topic = r.FormValue("topic")
	c.Channel = r.FormValue("channel")
	//描述
	c.Description = r.FormValue("description")
	//责任人
	c.Owner = r.FormValue("owner")
	//积压报警阈值
	monitorThreshold, err := strconv.Atoi(r.FormValue("monitorThreshold"))
	if err == nil && monitorThreshold >= 0 {
		c.MonitorThreshold = monitorThreshold
	}
	//该队列的并发量
	handleNum, err := strconv.Atoi(r.FormValue("handleNum"))
	if err == nil && handleNum >= 0 {
		c.HandleNum = handleNum
	}
	//NSQD最多同时推送多少个消息
	maxInFlight, err := strconv.Atoi(r.FormValue("maxInFlight"))
	if err == nil && maxInFlight >= 0 {
		c.MaxInFlight = maxInFlight
	}
	//失败，超时等情况是否重新入队
	isRequeue, err := strconv.ParseBool(r.FormValue("isRequeue"))
	if err == nil && isRequeue {
		c.IsRequeue = isRequeue
	}
	//超时时间
	timeoutDial, err := strconv.ParseInt(r.FormValue("timeoutDial"), 10, 64)
	if err == nil && timeoutDial >= 0 {
		c.TimeoutDial = time.Duration(timeoutDial)
	}
	//读超时时间
	timeoutRead, err := strconv.ParseInt(r.FormValue("timeoutRead"), 10, 64)
	if err == nil && timeoutRead >= 0 {
		c.TimeoutRead = time.Duration(timeoutRead)
	}
	//写超时时间
	timeoutWrite, err := strconv.ParseInt(r.FormValue("timeoutWrite"), 10, 64)
	if err == nil && timeoutWrite >= 0 {
		c.TimeoutWrite = time.Duration(timeoutWrite)
	}
	//是否有效
	invalid, err := strconv.Atoi(r.FormValue("invalid"))
	if err != nil || invalid <= 0 || (invalid != model.InvalidAvailable && invalid != model.InvalidUnavailable) {
		invalid = model.InvalidAvailable
	}
	c.Invalid = invalid

	_, err = c.Update()
	if err != nil {
		Failed(w, HttpCodeBadRequest, "upload failed. err: "+err.Error())
		return
	}
	Success(w, "ok")
}

func (ConsumeConfig) Get(w http.ResponseWriter, r *http.Request) {
	c := &model.ConsumeConfig{}
	var err error
	c.Id, err = strconv.Atoi(r.FormValue("id"))
	if err != nil || c.Id <= 0 {
		Failed(w, HttpCodeBadRequest, "param error")
		return
	}
	n, err := c.Get()
	if n == 0 || c.Id <= 0 {
		Failed(w, HttpCodeNotFound, "not found")
		return
	}
	Success(w, c)
}

func (ConsumeConfig) Page(w http.ResponseWriter, r *http.Request) {
	c := &model.ConsumeConfig{}
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	topic := r.FormValue("topic")
	pageResult, err := c.Page(topic, page)
	if err != nil {
		Failed(w, HttpCodeInternalServerError, "please try again. err: "+err.Error())
		return
	}
	Success(w, pageResult)
}

func (ConsumeConfig) WorkList(w http.ResponseWriter, r *http.Request) {
	c := &model.ConsumeConfig{}
	var err error
	c.Id, err = strconv.Atoi(r.FormValue("id"))
	if err != nil || c.Id <= 0 {
		Failed(w, HttpCodeBadRequest, "param error")
		return
	}
	err = c.WorkList()
	if err != nil {
		Failed(w, HttpCodeInternalServerError, "please try again. err: "+err.Error())
		return
	}
	Success(w, c)
}
