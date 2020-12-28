package model

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

const (
	InvalidAvailable   = 0 //可用
	InvalidUnavailable = 1 //不可用
)

type PageResult struct {
	Total  int64       `json:"total"`
	Page   int         `json:"page"`
	Result interface{} `json:"result"`
}

var db *gorm.DB

func NewDB(host, port, username, password, dbname string) {
	var err error
	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8&parseTime=true&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("open mysql error: " + err.Error())
	}
	err = ConsumeConfig{}.CreateTable()
	if err != nil {
		panic("create table error: " + err.Error())
	}
	err = ConsumeServerMap{}.CreateTable()
	if err != nil {
		panic("create table error: " + err.Error())
	}
	err = WorkServer{}.CreateTable()
	if err != nil {
		panic("create table error: " + err.Error())
	}
}

func IsErrRecordNotFound(err error) bool {
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}
	return false
}

//获取可用的消费者列表
func GetAvailableConsumeList() ([]ConsumeConfig, error) {
	//执行sql
	consumeConfigList := make([]ConsumeConfig, 0)
	consumeConfigListDB := make([]ConsumeConfig, 0)
	if db == nil {
		return consumeConfigList, errors.New("db is nil")
	}
	result := db.Where("invalid = ?", InvalidAvailable).Find(&consumeConfigListDB)
	if result.Error != nil {
		if IsErrRecordNotFound(result.Error) {
			return consumeConfigList, nil
		}
		//数据库出错，直接跳过本次查库
		return consumeConfigList, result.Error
	}
	if result.RowsAffected <= 0 {
		return consumeConfigList, nil
	}
	for _, consumeConfig := range consumeConfigListDB {
		consumeConfig.Consumer = nil
		consumeConfig.TimeoutWrite = consumeConfig.TimeoutWrite * time.Second
		consumeConfig.TimeoutRead = consumeConfig.TimeoutRead * time.Second
		consumeConfig.TimeoutDial = consumeConfig.TimeoutDial * time.Second
		consumeConfig.SetStatusWait()

		//获取消费者配置和worker机关联关系
		consumeServerMapList := make([]ConsumeServerMap, 0)
		consumeServerMapListDB := make([]ConsumeServerMap, 0)
		result = db.Where("consumeid = ? AND weight > 0 AND invalid = ?", consumeConfig.Id, InvalidAvailable).Find(&consumeServerMapListDB)
		if result.Error != nil || result.RowsAffected <= 0 || len(consumeServerMapListDB) <= 0 {
			continue
		}

		//获取每个消费者对应的work机器列表
		for _, consumeServerMap := range consumeServerMapListDB {
			workServer := WorkServer{}
			result = db.Where("id = ? AND invalid = ?", consumeServerMap.Serverid, InvalidAvailable).First(&workServer)
			if result.Error != nil || result.RowsAffected <= 0 || workServer.Id <= 0 {
				continue
			}
			workServer.Protocol = strings.ToLower(workServer.Protocol)
			workServer.SetStatusAvailable()
			consumeServerMap.WorkServer = workServer
			consumeServerMapList = append(consumeServerMapList, consumeServerMap)
		}
		if len(consumeServerMapList) <= 0 {
			continue
		}
		consumeConfig.ServerList = consumeServerMapList
		consumeConfigList = append(consumeConfigList, consumeConfig)
	}
	return consumeConfigList, nil
}
