package model

import (
	"errors"
	"time"
)

type ConsumeServerMap struct {
	Id        int `json:"id" gorm:"primaryKey"`
	Consumeid int `json:"consumeid"`
	Serverid  int `json:"serverid"`
	Weight    int `json:"weight"`
	Invalid   int `json:"invalid"`
	//创建时间
	CreatedAt time.Time `json:"createdAt"`
	//更新时间
	UpdatedAt time.Time `json:"updatedAt"`

	WorkServer WorkServer `json:"workServer" gorm:"-"`
}

func (ConsumeServerMap) TableName() string {
	return "nsqproxy_consume_server_map"
}

func (ConsumeServerMap) CreateTable() error {
	sql := "CREATE TABLE IF NOT EXISTS `nsqproxy_consume_server_map` (" +
		"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`consumeid` int(11) NOT NULL COMMENT '消费者id'," +
		"`serverid` int(11) NOT NULL COMMENT '服务器id'," +
		"`weight` int(11) DEFAULT '0' COMMENT '权重'," +
		"`invalid` tinyint(4) DEFAULT '0' COMMENT '是否有效，0是有效'," +
		"`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间'," +
		"`updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间'," +
		"PRIMARY KEY (`id`)," +
		"UNIQUE KEY `index_uq_cid_sid` (`consumeid`,`serverid`)" +
		") ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='消费者和可消费的服务器之间的关联关系';"
	return db.Exec(sql).Error
}

//两份配置是否相等
func (m ConsumeServerMap) IsEqual(newMap ConsumeServerMap) bool {
	if m.Id != newMap.Id || m.Consumeid != newMap.Consumeid || m.Serverid != newMap.Serverid || m.Weight != newMap.Weight || m.Invalid != newMap.Invalid {
		return false
	}
	if !m.WorkServer.IsEqual(newMap.WorkServer) {
		return false
	}
	return true
}

func (m *ConsumeServerMap) Create() (int, error) {
	result := db.Create(m)
	if result.Error != nil {
		return 0, result.Error
	} else if result.RowsAffected <= 0 {
		return 0, errors.New("RowsAffected is zero")
	} else if m.Id <= 0 {
		return 0, errors.New("primaryKey is zero")
	}
	return m.Id, nil
}

func (m *ConsumeServerMap) Delete() (int64, error) {
	if m.Id <= 0 {
		return 0, errors.New("primaryKey is zero")
	}
	result := db.Delete(m, m.Id)
	return result.RowsAffected, result.Error
}

func (m *ConsumeServerMap) Update() (int64, error) {
	if m.Id <= 0 {
		return 0, errors.New("primaryKey is zero")
	}
	result := db.Select("Id", "Consumeid", "Serverid", "Weight", "Invalid", "UpdatedAt").Updates(m)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (m *ConsumeServerMap) Get() (int64, error) {
	if m.Id <= 0 {
		return 0, errors.New("primaryKey is zero")
	}
	result := db.First(m)
	return result.RowsAffected, result.Error
}

func (m *ConsumeServerMap) AllByConsumeid(consumeid int) ([]ConsumeServerMap, error) {
	mList := make([]ConsumeServerMap, 0)
	result := db.Where("consumeid = ?", consumeid).Find(&mList)
	return mList, result.Error
}
