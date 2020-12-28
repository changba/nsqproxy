### 查询一个消费者

- 接口功能

> 根据主键ID查询一个消费者

```
GET /admin/api/consumeConfig/get
```

- 请求参数

|参数|必选|类型|说明|
|:----- |:-------|:-----|----- |
|id |true |int |消费者ID |
 
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
        serverList: null
    }
}
```