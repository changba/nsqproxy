package model

import (
	"testing"
)

func TestConsumeServerMap_CreateTable(t *testing.T) {
	NewDB("0.0.0.0", "3306", "root", "", "nsqproxy")
	m := ConsumeServerMap{}
	err := m.CreateTable()
	if err != nil{
		t.Fatalf("create table error: %s", err.Error())
	}
	db, err := db.DB()
	if err != nil{
		t.Fatalf("get db error: %s", err.Error())
	}
	err = db.Close()
	if err != nil{
		t.Fatalf("close db error: %s", err.Error())
	}
}

func TestConsumeServerMap_CURD(t *testing.T) {
	NewDB("0.0.0.0", "3306", "root", "", "nsqproxy")
	m := ConsumeServerMap{
		Consumeid: 1,
		Serverid:  2,
		Weight:    1,
	}
	//新增
	result := db.Table("nsqproxy_consume_server_map").Create(&m)
	id := m.Id
	if result.Error != nil{
		t.Fatalf("[insert]error: %s", result.Error.Error())
	}
	if result.RowsAffected != 1{
		t.Fatalf("[insert]RowsAffected is not 1. result: %d", result.RowsAffected)
	}
	if id <= 0{
		t.Fatalf("[insert]id is 0")
	}

	//查询
	find := ConsumeServerMap{}
	result = db.Table("nsqproxy_consume_server_map").First(&find, id)
	if result.Error != nil{
		t.Fatalf("[select]error: %s", result.Error.Error())
	}
	if result.RowsAffected != 1{
		t.Fatalf("[select]RowsAffected is not 1. result: %d", result.RowsAffected)
	}

	//修改
	m.Invalid = 10
	m.Weight = 11
	result = db.Save(&m)
	if result.Error != nil{
		t.Fatalf("[update]error: %s", result.Error.Error())
	}
	if result.RowsAffected != 1{
		t.Fatalf("[update]RowsAffected is not 1. result: %d", result.RowsAffected)
	}
	if m.Id != id{
		t.Fatalf("[update]id has changed")
	}

	//查询
	find = ConsumeServerMap{}
	result = db.First(&find, id)
	if result.Error != nil{
		t.Fatalf("[select2]error: %s", result.Error.Error())
	}
	if find.Invalid != 10 || find.Weight != 11 {
		t.Fatalf("[select2]update failed")
	}
	if result.RowsAffected != 1{
		t.Fatalf("[select2]RowsAffected is not 1. result: %d", result.RowsAffected)
	}
	if !m.IsEqual(find){
		t.Fatalf("[select]Query result failed")
	}

	//删除
	result = db.Delete(&m)
	if result.Error != nil{
		t.Fatalf("[delete]error: %s", result.Error.Error())
	}
	if result.RowsAffected != 1{
		t.Fatalf("[delete]RowsAffected is not 1. result: %d", result.RowsAffected)
	}

	//查询
	find = ConsumeServerMap{}
	result = db.First(&find, id)
	if find.Serverid != 0{
		t.Fatalf("[select3]id is not 0")
	}
	if result.Error == nil{
		t.Fatalf("[select3]no error")
	}
	if result.RowsAffected != 0{
		t.Fatalf("[select3]RowsAffected is not 1. result: %d", result.RowsAffected)
	}

	//关闭
	db, err := db.DB()
	if err != nil{
		t.Fatalf("get db error: %s", err.Error())
	}
	err = db.Close()
	if err != nil{
		t.Fatalf("close db error: %s", err.Error())
	}
}