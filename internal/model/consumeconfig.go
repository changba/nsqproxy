package model

import (
	"errors"
	"github.com/nsqio/go-nsq"
	"strings"
	"sync/atomic"
	"time"
)

const (
	consumeConfigStatusWait    = int32(0)
	consumeConfigStatusSuccess = int32(1)
	consumeConfigStatusFailed  = int32(2)
	consumeConfigStatusClosed  = int32(3)
)

// 表示一个队列，即唯一的Topic+Channel
type ConsumeConfig struct {
	//主键ID
	Id int `json:"id" gorm:"primaryKey"`
	//队列名
	Topic string `json:"topic"`
	//通道名
	Channel string `json:"channel"`
	//描述
	Description string `json:"description"`
	//责任人
	Owner string `json:"owner"`
	//积压报警阈值
	MonitorThreshold int `json:"monitorThreshold"`
	//该队列的并发量
	HandleNum int `json:"handleNum"`
	//NSQD最多同时推送多少个消息
	MaxInFlight int `json:"maxInFlight"`
	//失败，超时等情况是否重新入队
	IsRequeue bool `json:"isRequeue"`
	//超时时间
	TimeoutDial time.Duration `json:"timeoutDial"`
	//读超时时间
	TimeoutRead time.Duration `json:"timeoutRead"`
	//写超时时间
	TimeoutWrite time.Duration `json:"timeoutWrite"`
	//是否有效
	Invalid int `json:"invalid"`
	//创建时间
	CreatedAt time.Time `json:"createdAt"`
	//更新时间
	UpdatedAt time.Time `json:"updatedAt"`

	//那些work机器可以消费该队列
	ServerList []ConsumeServerMap `json:"serverList" gorm:"-"`
	//nsq客户端的消费者
	Consumer *nsq.Consumer `json:"-" gorm:"-"`
	//nsq客户端的消费者在本系统中的唯一ID
	consumerUniqueId [16]byte `gorm:"-"`
	//运行状态
	status int32 `gorm:"-"`
}

func (ConsumeConfig) TableName() string {
	return "nsqproxy_consume_config"
}

//定义结构体，然后让ORM来帮你建表，这个结构体的标签我觉得写起来贼麻烦，还不如直接来建表语句。
func (ConsumeConfig) CreateTable() error {
	sql := "CREATE TABLE IF NOT EXISTS `nsqproxy_consume_config` (" +
		"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`topic` varchar(100) NOT NULL DEFAULT '' COMMENT 'topic名'," +
		"`channel` varchar(100) NOT NULL DEFAULT '' COMMENT 'channel名'," +
		"`description` varchar(1000) NOT NULL DEFAULT '' COMMENT '描述'," +
		"`owner` varchar(12) NOT NULL DEFAULT '' COMMENT '责任人'," +
		"`monitor_threshold` int(11) NOT NULL DEFAULT '50000' COMMENT '报警监控的阈值, 0是白名单'," +
		"`handle_num` int(11) NOT NULL DEFAULT '2' COMMENT '消费者的并发量'," +
		"`max_in_flight` int(11) NOT NULL DEFAULT '2' COMMENT '未返回时nsqd最大的可推送量'," +
		"`is_requeue` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否需要重新入队，0是不需要，1是需要'," +
		"`timeout_dial` int(11) NOT NULL DEFAULT '3590' COMMENT '超时时间'," +
		"`timeout_read` int(11) NOT NULL DEFAULT '3590' COMMENT '超时时间-读'," +
		"`timeout_write` int(11) NOT NULL DEFAULT '3590' COMMENT '超时时间-写'," +
		"`invalid` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否有效，0是有效'," +
		"`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间'," +
		"`updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间'," +
		"PRIMARY KEY (`id`)," +
		"UNIQUE KEY `uniq_topic_channel` (`topic`,`channel`)" +
		") ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='消费者的配置';"
	return db.Exec(sql).Error
}

//两份配置是否相等
func (c ConsumeConfig) IsEqual(newConsume ConsumeConfig) bool {
	if c.Id != newConsume.Id || c.Topic != newConsume.Topic || c.Channel != newConsume.Channel || c.Description != newConsume.Description || c.Owner != newConsume.Owner {
		return false
	}
	if c.MonitorThreshold != newConsume.MonitorThreshold || c.HandleNum != newConsume.HandleNum || c.MaxInFlight != newConsume.MaxInFlight || c.IsRequeue != newConsume.IsRequeue {
		return false
	}
	if c.TimeoutDial != newConsume.TimeoutDial || c.TimeoutRead != newConsume.TimeoutRead || c.TimeoutWrite != newConsume.TimeoutWrite || c.Invalid != newConsume.Invalid {
		return false
	}
	if len(c.ServerList) != len(newConsume.ServerList) {
		return false
	}
	for index, ServerMap := range c.ServerList {
		if !ServerMap.IsEqual(newConsume.ServerList[index]) {
			return false
		}
	}
	return true
}

func (c ConsumeConfig) GetStatus() int32 {
	return atomic.LoadInt32(&c.status)
}

func (c *ConsumeConfig) SetStatusWait() {
	atomic.StoreInt32(&c.status, consumeConfigStatusWait)
}

func (c ConsumeConfig) StatusIsSuccess() bool {
	return c.GetStatus() == consumeConfigStatusSuccess
}

func (c *ConsumeConfig) SetStatusSuccess() {
	atomic.StoreInt32(&c.status, consumeConfigStatusSuccess)
}

func (c *ConsumeConfig) SetStatusFailed() {
	atomic.StoreInt32(&c.status, consumeConfigStatusFailed)
}

func (c ConsumeConfig) StatusIsClose() bool {
	return c.GetStatus() == consumeConfigStatusClosed
}

func (c *ConsumeConfig) SetStatusClosed() {
	atomic.StoreInt32(&c.status, consumeConfigStatusClosed)
}

func (c *ConsumeConfig) SetConsumerUniqueId(consumerUniqueId [16]byte) {
	c.consumerUniqueId = consumerUniqueId
}

func (c *ConsumeConfig) GetConsumerUniqueId() [16]byte {
	return c.consumerUniqueId
}

//把数组形式的配置列表转换成map，key是id
func FormatConsumeConfigListForMap(consumeConfigList []ConsumeConfig) map[int]ConsumeConfig {
	consumeConfigMap := make(map[int]ConsumeConfig)
	for _, consumeConfig := range consumeConfigList {
		consumeConfigMap[consumeConfig.Id] = consumeConfig
	}
	return consumeConfigMap
}

func (c *ConsumeConfig) Create() (int, error) {
	result := db.Create(c)
	if result.Error != nil {
		return 0, result.Error
	} else if result.RowsAffected <= 0 {
		return 0, errors.New("RowsAffected is zero")
	} else if c.Id <= 0 {
		return 0, errors.New("primaryKey is zero")
	}
	return c.Id, nil
}

func (c *ConsumeConfig) Delete() (int64, error) {
	if c.Id <= 0 {
		return 0, errors.New("primaryKey is zero")
	}
	result := db.Delete(c, c.Id)
	return result.RowsAffected, result.Error
}

func (c *ConsumeConfig) Update() (int64, error) {
	if c.Id <= 0 {
		return 0, errors.New("primaryKey is zero")
	}
	result := db.Select("Id", "Topic", "Channel", "Description", "Owner", "MonitorThreshold", "HandleNum", "MaxInFlight", "IsRequeue", "TimeoutDial", "TimeoutRead", "TimeoutWrite", "Invalid", "UpdatedAt").Updates(c)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (c *ConsumeConfig) Get() (int64, error) {
	if c.Id <= 0 {
		return 0, errors.New("primaryKey is zero")
	}
	result := db.First(c)
	return result.RowsAffected, result.Error
}

func (c *ConsumeConfig) Page(topic string, page int) (PageResult, error) {
	var cList []ConsumeConfig
	//where部分
	whereList := make([]string, 0)
	if len(topic) > 0 {
		whereList = append(whereList, " topic LIKE '%"+topic+"%'")
	}
	d := db.Table(c.TableName()).Where(strings.Join(whereList, " AND "))
	//count部分
	var total int64
	result := d.Count(&total)
	if result.Error != nil || result.RowsAffected != 1 {
		total = 0
	}
	//page部分
	if page <= 0 {
		page = 1
	}
	result = d.Offset((page - 1) * 20).Limit(20).Find(&cList)
	pageRet := PageResult{
		Total:  total,
		Page:   page,
		Result: cList,
	}
	return pageRet, result.Error
}

func (c *ConsumeConfig) WorkList() error {
	//获取消费者配置
	n, err := c.Get()
	if err != nil && !IsErrRecordNotFound(err) {
		return err
	}
	if n == 0 || c.Id <= 0 || IsErrRecordNotFound(err) {
		return nil
	}
	//获取map列表
	csMapList, err := (&ConsumeServerMap{}).AllByConsumeid(c.Id)
	if err != nil {
		return err
	}
	if len(csMapList) <= 0 {
		return nil
	}
	for k, csMap := range csMapList {
		if csMap.Serverid < 0 {
			continue
		}
		work := WorkServer{}
		work.Id = csMap.Serverid
		n, err := work.Get()
		if err != nil || n == 0 || work.Id <= 0 {
			continue
		}
		csMapList[k].WorkServer = work
	}
	c.ServerList = csMapList
	return nil
}
