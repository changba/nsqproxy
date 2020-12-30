package httper

import (
	"github.com/changba-server/nsqproxy/internal/model"
	"net/http"
	"strconv"
)

type WorkServer struct {
}

func (WorkServer) Create(w http.ResponseWriter, r *http.Request) {
	work := &model.WorkServer{}
	work.Addr = r.FormValue("addr")
	work.Protocol = r.FormValue("protocol")
	work.Extra = r.FormValue("extra")
	work.Description = r.FormValue("description")
	work.Owner = r.FormValue("owner")
	work.Invalid = model.InvalidAvailable
	invalid, err := strconv.Atoi(r.FormValue("invalid"))
	if err != nil || invalid <= 0 || (invalid != model.InvalidAvailable && invalid != model.InvalidUnavailable) {
		invalid = model.InvalidAvailable
	}
	work.Invalid = invalid
	id, err := work.Create()
	if err != nil {
		Failed(w, HttpCodeBadRequest, "create failed. err: "+err.Error())
		return
	}
	if id <= 0 {
		Failed(w, HttpCodeBadRequest, "id is zero")
		return
	}
	Success(w, work)
}

func (WorkServer) Delete(w http.ResponseWriter, r *http.Request) {
	work := &model.WorkServer{}
	var err error
	work.Id, err = strconv.Atoi(r.FormValue("id"))
	if err != nil || work.Id <= 0 {
		Failed(w, HttpCodeBadRequest, "param error")
		return
	}
	_, err = work.Delete()
	if err != nil {
		Failed(w, HttpCodeBadRequest, "delete failed. err: "+err.Error())
		return
	}
	Success(w, "ok")
}

func (WorkServer) Update(w http.ResponseWriter, r *http.Request) {
	work := &model.WorkServer{}
	var err error
	work.Id, err = strconv.Atoi(r.FormValue("id"))
	if err != nil || work.Id <= 0 {
		Failed(w, HttpCodeBadRequest, "param error")
		return
	}
	work.Addr = r.FormValue("addr")
	work.Protocol = r.FormValue("protocol")
	work.Extra = r.FormValue("extra")
	work.Description = r.FormValue("description")
	work.Owner = r.FormValue("owner")
	invalid, err := strconv.Atoi(r.FormValue("invalid"))
	if err != nil || invalid <= 0 || (invalid != model.InvalidAvailable && invalid != model.InvalidUnavailable) {
		invalid = model.InvalidAvailable
	}
	work.Invalid = invalid

	_, err = work.Update()
	if err != nil {
		Failed(w, HttpCodeBadRequest, "upload failed. err: "+err.Error())
		return
	}
	Success(w, "ok")
}

func (WorkServer) Get(w http.ResponseWriter, r *http.Request) {
	work := &model.WorkServer{}
	var err error
	work.Id, err = strconv.Atoi(r.FormValue("id"))
	if err != nil || work.Id <= 0 {
		Failed(w, HttpCodeBadRequest, "param error")
		return
	}
	n, err := work.Get()
	if err != nil || n == 0 || work.Id <= 0 {
		Failed(w, HttpCodeNotFound, "not found")
		return
	}
	Success(w, work)
}

func (WorkServer) Page(w http.ResponseWriter, r *http.Request) {
	work := &model.WorkServer{}
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	pageResult, err := work.Page(page)
	if err != nil {
		Failed(w, HttpCodeInternalServerError, "please try again. err: "+err.Error())
		return
	}
	Success(w, pageResult)
}

func (WorkServer) All(w http.ResponseWriter, r *http.Request) {
	work := &model.WorkServer{}
	wList, err := work.All()
	if err != nil {
		Failed(w, HttpCodeInternalServerError, "please try again. err: "+err.Error())
		return
	}
	Success(w, wList)
}
