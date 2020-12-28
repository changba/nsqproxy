package model

import (
	"testing"
)

func TestWorkServer_CreateTable(t *testing.T) {
	NewDB("0.0.0.0", "3306", "root", "", "nsqproxy")
	w := WorkServer{}
	err := w.CreateTable()
	if err != nil {
		t.Fatalf("create table error: %s", err.Error())
	}
	db, err := db.DB()
	if err != nil {
		t.Fatalf("get db error: %s", err.Error())
	}
	err = db.Close()
	if err != nil {
		t.Fatalf("close db error: %s", err.Error())
	}
}

func TestWorkServer_CURD(t *testing.T) {
	NewDB("0.0.0.0", "3306", "root", "", "nsqproxy")
	w := WorkServer{
		Addr:     "0.0.0.0:80",
		Protocol: "HTTP",
	}
	//新增
	result := db.Table("nsqproxy_work_server").Create(&w)
	id := w.Id
	if result.Error != nil {
		t.Fatalf("[insert]error: %s", result.Error.Error())
	}
	if result.RowsAffected != 1 {
		t.Fatalf("[insert]RowsAffected is not 1. result: %d", result.RowsAffected)
	}
	if id <= 0 {
		t.Fatalf("[insert]serverid is 0")
	}

	//查询
	find := WorkServer{}
	result = db.Table("nsqproxy_work_server").First(&find, id)
	if result.Error != nil {
		t.Fatalf("[select]error: %s", result.Error.Error())
	}
	if result.RowsAffected != 1 {
		t.Fatalf("[select]RowsAffected is not 1. result: %d", result.RowsAffected)
	}
	if !w.IsEqual(find) {
		t.Fatalf("[select]Query result failed")
	}

	//修改
	w.Addr = "0.0.0.0:81"
	w.Protocol = "CBNSQ"
	result = db.Save(&w)
	if result.Error != nil {
		t.Fatalf("[update]error: %s", result.Error.Error())
	}
	if result.RowsAffected != 1 {
		t.Fatalf("[update]RowsAffected is not 1. result: %d", result.RowsAffected)
	}
	if w.Id != id {
		t.Fatalf("[update]id has changed")
	}

	//查询
	find = WorkServer{}
	result = db.First(&find, id)
	if result.Error != nil {
		t.Fatalf("[select2]error: %s", result.Error.Error())
	}
	if find.Addr != "0.0.0.0:81" || find.Protocol != "CBNSQ" {
		t.Fatalf("[select2]update failed")
	}
	if result.RowsAffected != 1 {
		t.Fatalf("[select2]RowsAffected is not 1. result: %d", result.RowsAffected)
	}
	if !w.IsEqual(find) {
		t.Fatalf("[select2]Query result failed")
	}
	//删除
	result = db.Delete(&w)
	if result.Error != nil {
		t.Fatalf("[delete]error: %s", result.Error.Error())
	}
	if result.RowsAffected != 1 {
		t.Fatalf("[delete]RowsAffected is not 1. result: %d", result.RowsAffected)
	}

	//查询
	find = WorkServer{}
	result = db.First(&find, id)
	if find.Id != 0 {
		t.Fatalf("[select3]id is not 0")
	}
	if result.Error == nil {
		t.Fatalf("[select3]no error")
	}
	if result.RowsAffected != 0 {
		t.Fatalf("[select3]RowsAffected is not 1. result: %d", result.RowsAffected)
	}

	//关闭
	db, err := db.DB()
	if err != nil {
		t.Fatalf("get db error: %s", err.Error())
	}
	err = db.Close()
	if err != nil {
		t.Fatalf("close db error: %s", err.Error())
	}
}
