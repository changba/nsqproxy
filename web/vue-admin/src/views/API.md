消费者配置：

新增：http://10.10.248.110:19421/api/consumeconfig/create?topic=test_api1
获取一条：http://10.10.248.110:19421/api/consumeconfig/get?id=2
更新：http://10.10.248.110:19421/api/consumeconfig/update?id=86&topic=test_api2
删除：http://10.10.248.110:19421/api/consumeconfig/delete?id=86
获取全部不分页：http://10.10.248.110:19421/api/consumeconfig/all?topic=hjk
获取第一页：http://10.10.248.110:19421/api/consumeconfig/all?page=1

type ConsumeConfig struct {
  //主键ID
  Id int `gorm:"primaryKey"`
  //队列名
  Topic string
  //通道名
  Channel string
  //描述
  Description string
  //责任人
  Owner string
  //积压报警阈值
  MonitorThreshold int
  //该队列的并发量
  HandleNum int
  //NSQD最多同时推送多少个消息
  MaxInFlight int
  //失败，超时等情况是否重新入队
  IsRequeue bool
  //超时时间
  TimeoutDial time.Duration
  //读超时时间
  TimeoutRead time.Duration
  //写超时时间
  TimeoutWrite time.Duration
  //是否暂停
  Invalid int
  //是否暂停
  Pause int
  //创建时间
  CreatedAt time.Time
  //更新时间
  UpdatedAt time.Time
}


分页 total page
list 字段与创建修改不对应
boolean类型默认值