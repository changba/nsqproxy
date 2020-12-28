package httper

import (
	"github.com/ChangbaServer/nsqproxy/internal/model"
	"net/http"
	"strconv"
	"strings"
)

type ConsumeServerMap struct{

}

func (ConsumeServerMap) Create(w http.ResponseWriter, r *http.Request){
	consumeid, err := strconv.Atoi(r.FormValue("consumeid"))
	if err != nil || consumeid <= 0{
		Failed(w, HttpCodeBadRequest, "param error. consumeid is required")
		return
	}
	weight, err := strconv.Atoi(r.FormValue("weight"))
	if err != nil || weight <= 0{
		Failed(w, HttpCodeBadRequest, "param error. weight is required")
		return
	}
	invalid, err := strconv.Atoi(r.FormValue("invalid"))
	if err != nil || invalid <= 0 || (invalid != model.InvalidAvailable && invalid != model.InvalidUnavailable){
		invalid = model.InvalidAvailable
	}
	//支持多选，逗号分割
	if len(r.FormValue("serverid")) < 0 {
		Failed(w, HttpCodeBadRequest, "param error. serverid is required")
		return
	}
	serveridStrList := strings.Split(r.FormValue("serverid"), ",")
	csMapList := make([]*model.ConsumeServerMap, 0)
	for _, serveridStr := range serveridStrList{
		serverid, err := strconv.Atoi(serveridStr)
		if err != nil || serverid <= 0{
			Failed(w, HttpCodeBadRequest, "param error. serverid must be int")
			return
		}
		csMap := &model.ConsumeServerMap{
			Consumeid:  consumeid,
			Serverid:   serverid,
			Weight:     weight,
			Invalid:    invalid,
		}
		id, err := csMap.Create()
		if err != nil{
			Failed(w, HttpCodeBadRequest, "create failed. err: " + err.Error())
			return
		}
		if id <= 0{
			Failed(w, HttpCodeBadRequest, "id is zero")
			return
		}
		csMapList = append(csMapList, csMap)
	}
	Success(w, csMapList)
}

func (ConsumeServerMap) Delete(w http.ResponseWriter, r *http.Request){
	csMap := &model.ConsumeServerMap{}
	var err error
	csMap.Id, err = strconv.Atoi(r.FormValue("id"))
	if err != nil || csMap.Id <= 0{
		Failed(w, HttpCodeBadRequest, "param error")
		return
	}
	_, err = csMap.Delete()
	if err != nil{
		Failed(w, HttpCodeBadRequest, "delete failed. err: " + err.Error())
		return
	}
	Success(w, "ok")
}

func (ConsumeServerMap) Update(w http.ResponseWriter, r *http.Request){
	csMap := &model.ConsumeServerMap{}
	var err error
	csMap.Id, err = strconv.Atoi(r.FormValue("id"))
	if err != nil || csMap.Id <= 0{
		Failed(w, HttpCodeBadRequest, "param error")
		return
	}
	csMap.Consumeid, err = strconv.Atoi(r.FormValue("consumeid"))
	if err != nil || csMap.Consumeid <= 0{
		Failed(w, HttpCodeBadRequest, "param error. consumeid is required")
		return
	}
	csMap.Weight, err = strconv.Atoi(r.FormValue("weight"))
	if err != nil || csMap.Weight <= 0{
		Failed(w, HttpCodeBadRequest, "param error. weight is required")
		return
	}
	csMap.Serverid, err = strconv.Atoi(r.FormValue("serverid"))
	if err != nil || csMap.Serverid <= 0{
		Failed(w, HttpCodeBadRequest, "param error. serverid is required")
		return
	}
	invalid, err := strconv.Atoi(r.FormValue("invalid"))
	if err != nil || invalid <= 0 || (invalid != model.InvalidAvailable && invalid != model.InvalidUnavailable){
		invalid = model.InvalidAvailable
	}
	csMap.Invalid = invalid

	_, err = csMap.Update()
	if err != nil{
		Failed(w, HttpCodeBadRequest, "upload failed. err: " + err.Error())
		return
	}
	Success(w, "ok")
}

func (ConsumeServerMap) GetWork(w http.ResponseWriter, r *http.Request){
	csMap := &model.ConsumeServerMap{}
	var err error
	csMap.Id, err = strconv.Atoi(r.FormValue("id"))
	if err != nil || csMap.Id <= 0{
		Failed(w, HttpCodeBadRequest, "param error")
		return
	}
	n, err := csMap.Get()
	if n == 0 || csMap.Id <= 0{
		Failed(w, HttpCodeNotFound, "not found")
		return
	}
	if csMap.Serverid <= 0{
		Failed(w, HttpCodeBadRequest, "work.Id is empty")
		return
	}
	work := model.WorkServer{}
	work.Id = csMap.Serverid
	n, err = work.Get()
	if err != nil || n == 0 || work.Id <= 0{
		Failed(w, HttpCodeNotFound, "work not found")
		return
	}
	csMap.WorkServer = work
	Success(w, csMap)
}