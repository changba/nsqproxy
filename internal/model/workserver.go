package model

import (
	"errors"
	"sync/atomic"
	"time"
)

//worker当前可用
const workserverStatusAvailable = 1

//worker当前不可用
const workserverStatusUnavailable = 0

type WorkServer struct {
	Id    int   `json:"id" gorm:"primaryKey"`
	//地址，IP:PORT
	Addr        string `json:"addr"`
	//协议，如HTTP、FastCGI、CBNSQ
	Protocol    string `json:"protocol"`
	//扩展字段
	Extra		string `json:"extra"`
	//描述
	Description string `json:"description"`
	//责任人
	Owner       string `json:"owner"`
	//是否有效
	Invalid     int `json:"invalid"`
	//创建时间
	CreatedAt time.Time `json:"createdAt"`
	//更新时间
	UpdatedAt time.Time `json:"updatedAt"`

	status      int32 `gorm:"-"`
}

func (WorkServer) TableName() string {
	return "nsqproxy_work_server"
}

func (WorkServer) CreateTable()error{
	sql := "CREATE TABLE IF NOT EXISTS `nsqproxy_work_server` (" +
		"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`addr` varchar(255) NOT NULL DEFAULT '' COMMENT '地址'," +
		"`protocol` varchar(11) DEFAULT 'CBNSQ' COMMENT '使用的协议，支持HTTP、FastCGI、CBNSQ'," +
		"`extra` varchar(1000) DEFAULT NULL COMMENT '扩展字段，比如协议是FastCGI时，需要传入PHP-FPM的执行的PHP文件的路径'," +
		"`description` varchar(1000) DEFAULT '' COMMENT '描述'," +
		"`owner` varchar(12) DEFAULT NULL COMMENT '责任人'," +
		"`invalid` tinyint(4) DEFAULT '0' COMMENT '是否有效，0是有效'," +
		"`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间'," +
		"`updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间'," +
		"PRIMARY KEY (`id`)," +
		"UNIQUE KEY `index_uq_addr` (`addr`)" +
		") ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='所有的可以消费的服务器列表';"
	return db.Exec(sql).Error
}

//两份配置是否相等
func (w WorkServer) IsEqual(newWork WorkServer) bool {
	if w.Id != newWork.Id || w.Addr != newWork.Addr  || w.Protocol != newWork.Protocol || w.Extra != newWork.Extra || w.Description != newWork.Description {
		return false
	}
	if w.Owner != newWork.Owner || w.Invalid != newWork.Invalid {
		return false
	}
	return true
}

func (w WorkServer) GetStatus(){
	atomic.LoadInt32(&w.status)
}

func (w *WorkServer) SetStatusAvailable(){
	atomic.StoreInt32(&w.status, workserverStatusAvailable)
}

func (w *WorkServer) SetStatusUnAvailable(){
	atomic.StoreInt32(&w.status, workserverStatusUnavailable)
}


func (w *WorkServer) Create()(int, error){
	result := db.Create(w)
	if result.Error != nil{
		return 0, result.Error
	}else if result.RowsAffected <= 0{
		return 0, errors.New("RowsAffected is zero")
	}else if w.Id <= 0{
		return 0, errors.New("primaryKey is zero")
	}
	return w.Id, nil
}

func (w *WorkServer) Delete()(int64, error){
	if w.Id <= 0{
		return 0, errors.New("primaryKey is zero")
	}
	result := db.Delete(w, w.Id)
	return result.RowsAffected, result.Error
}

func (w *WorkServer) Update()(int64, error){
	if w.Id <= 0{
		return 0, errors.New("primaryKey is zero")
	}
	result := db.Select("Id", "Addr", "Protocol", "Extra", "Description", "Owner", "Invalid", "UpdatedAt").Updates(w)
	if result.Error != nil{
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (w *WorkServer) Get()(int64, error){
	if w.Id <= 0{
		return 0, errors.New("primaryKey is zero")
	}
	result := db.First(w)
	return result.RowsAffected, result.Error
}

func (w *WorkServer) Page(page int)(PageResult, error){
	var wList []WorkServer
	d := db.Table(w.TableName())
	//count部分
	var total int64
	result := d.Count(&total)
	if result.Error != nil || result.RowsAffected != 1{
		total = 0
	}
	//page部分
	if page <= 0{
		page = 1
	}
	result = d.Offset((page-1)*20).Limit(20).Find(&wList)
	pageRet := PageResult{
		Total:  total,
		Page:   page,
		Result: wList,
	}
	return pageRet, result.Error
}

func (w *WorkServer) All()([]WorkServer, error){
	var wList []WorkServer
	result := db.Find(&wList)
	return wList, result.Error
}