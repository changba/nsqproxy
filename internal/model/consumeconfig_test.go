package model

import (
	"testing"
)

func TestConsumeConfig_CreateTable(t *testing.T) {
	NewDB("0.0.0.0", "3306", "root", "", "nsqproxy")
	c := ConsumeConfig{}
	err := c.CreateTable()
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

func TestConsumeConfig_CURD(t *testing.T) {
	NewDB("0.0.0.0", "3306", "root", "", "nsqproxy")
	c := ConsumeConfig{
		Topic:       "golang_test_topic",
		Channel:     "golang_test_channel",
		Description: "go test",
		Owner:       "go",
	}
	//新增
	result := db.Table("nsqproxy_consume_config").Create(&c)
	id := c.Id
	if result.Error != nil {
		t.Fatalf("[insert]error: %s", result.Error.Error())
	}
	if result.RowsAffected != 1 {
		t.Fatalf("[insert]RowsAffected is not 1. result: %d", result.RowsAffected)
	}
	if id <= 0 {
		t.Fatalf("[insert]consumeid is 0")
	}

	//查询
	find := ConsumeConfig{}
	result = db.Table("nsqproxy_consume_config").First(&find, id)
	if result.Error != nil {
		t.Fatalf("[select]error: %s", result.Error.Error())
	}
	if result.RowsAffected != 1 {
		t.Fatalf("[select]RowsAffected is not 1. result: %d", result.RowsAffected)
	}
	if !c.IsEqual(find) {
		t.Fatalf("[select]Query result failed")
	}

	//修改
	c.MonitorThreshold = 100
	c.HandleNum = 10
	c.MaxInFlight = 20
	result = db.Save(&c)
	if result.Error != nil {
		t.Fatalf("[update]error: %s", result.Error.Error())
	}
	if result.RowsAffected != 1 {
		t.Fatalf("[update]RowsAffected is not 1. result: %d", result.RowsAffected)
	}
	if c.Id != id {
		t.Fatalf("[update]id has changed")
	}

	//查询
	find = ConsumeConfig{}
	result = db.First(&find, id)
	if result.Error != nil {
		t.Fatalf("[select2]error: %s", result.Error.Error())
	}
	if find.MonitorThreshold != 100 || find.HandleNum != 10 || find.MaxInFlight != 20 {
		t.Fatalf("[select2]update failed")
	}
	if result.RowsAffected != 1 {
		t.Fatalf("[select2]RowsAffected is not 1. result: %d", result.RowsAffected)
	}
	if !c.IsEqual(find) {
		t.Fatalf("[select2]Query result failed")
	}

	//删除
	result = db.Delete(&c)
	if result.Error != nil {
		t.Fatalf("[delete]error: %s", result.Error.Error())
	}
	if result.RowsAffected != 1 {
		t.Fatalf("[delete]RowsAffected is not 1. result: %d", result.RowsAffected)
	}

	//查询
	find = ConsumeConfig{}
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
