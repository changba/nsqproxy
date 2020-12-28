### 批量查询消费者

- 接口功能

> 批量查询消费者

```
GET /admin/api/consumeConfig/page
```

- 请求参数

|参数|必选|类型|说明|
|:----- |:-------|:-----|----- |
|topic |false |string |Topic名，模糊查询，%topic%的方式查询 |
|page |false |int |分页，默认1，小于等于零时表示第一页，一页20条 |
 
- 返回结果

```
{
    code: 200,
    msg: "ok",
    result: [
        {
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
        },
        {
            id: 2,
            topic: "test_topic_2",
            channel: "one",
            description: "描述",
            owner: "责任人",
            monitorThreshold: 0,
            handleNum: 2,
            maxInFlight: 2,
            isRequeue: false,
            timeoutDial: 3590,
            timeoutRead: 3590,
            timeoutWrite: 3590,
            invalid: 0,
            createdAt: "2020-08-20T15:08:39+08:00",
            updatedAt: "0001-01-01T00:00:00Z",
            serverList: null
        }
    ]
}
```