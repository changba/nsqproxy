### 更新消费者

- 接口功能

> 更新一个消费者

```
GET /admin/api/consumeConfig/update
```

- 请求参数

|参数|必选|类型|说明|
|:----- |:-------|:-----|----- |
|id |true |int |消费者ID |
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
    result: "ok"
}
```