### 创建消费者

- 接口功能

> 创建一个新的消费者

```
GET /admin/api/consumeConfig/create
```

- 请求参数

|参数|必选|类型|说明|
|:----- |:-------|:-----|----- |
|topic |true |string |Topic名 |
|channel |true |string |Channel名 |
|description |false |string |描述，默认空 |
|owner |false |string |责任人，默认空 |
|monitorThreshold |false |int |积压报警阈值，默认50000 |
|handleNum |false |int |该队列的并发量，默认2 |
|maxInFlight |false |int |NSQD最多同时推送多少个消息，默认2 |
|isRequeue |false |bool |失败，超时等情况是否重新入队，默认false |
|timeoutDial |false |int |超时时间，默认3590秒 |
|timeoutRead |false |int |读超时时间，默认3590秒 |
|timeoutWrite |false |int |写超时时间，默认3590秒 |
|invalid |false |int |是否有效，默认0有效，1无效 |

- 返回结果

```
{
    code: 200,
    msg: "ok",
    result: {
        id: 1,
        topic: "test_topic_1",
        channel: "one",
        description: "描述",
        owner: "责任人",
        monitorThreshold: 0,
        handleNum: 6,
        maxInFlight: 6,
        isRequeue: false,
        timeoutDial: 3590,
        timeoutRead: 3590,
        timeoutWrite: 3590,
        invalid: 0,
        createdAt: "2020-08-20T15:08:39+08:00",
        updatedAt: "0001-01-01T00:00:00Z",
        serverList: null,
    }
}
```