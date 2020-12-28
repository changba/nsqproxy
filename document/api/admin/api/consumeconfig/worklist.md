### 查询一个消费者关联服务器

- 接口功能

> 根据消费者ID，查询消费者所关联的所有服务器列表

```
GET /admin/api/consumeConfig/workList
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
        topic: "test_topic",
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
        serverList: [
            {
                id: 1,
                consumeid: 1,
                serverid: 1,
                weight: 1,
                invalid: 0,
                createdAt: "2018-10-16T14:36:33+08:00",
                updatedAt: "0001-01-01T00:00:00Z",
                workServer: {
                    id: 1,
                    addr: "1.1.1.1:80,
                    protocol: "HTTP",
                    extra: "",
                    description: "通用机器1",
                    owner: "",
                    invalid: 0,
                    createdAt: "2018-10-16T14:30:58+08:00",
                    updatedAt: "0001-01-01T00:00:00Z"
                }
            },
            {
                id: 2,
                consumeid: 1,
                serverid: 2,
                weight: 1,
                invalid: 0,
                createdAt: "2018-10-16T14:36:33+08:00",
                updatedAt: "0001-01-01T00:00:00Z",
                workServer: {
                    id: 2,
                    addr: "1.1.1.2:80,
                    protocol: "HTTP",
                    extra: "",
                    description: "通用机器2",
                    owner: "",
                    invalid: 0,
                    createdAt: "2018-10-16T14:30:58+08:00",
                    updatedAt: "0001-01-01T00:00:00Z"
                }
            },
            {
                id: 3,
                consumeid: 1,
                serverid: 3,
                weight: 1,
                invalid: 0,
                createdAt: "2018-10-16T14:36:33+08:00",
                updatedAt: "0001-01-01T00:00:00Z",
                workServer: {
                    id: 3,
                    addr: "1.1.1.3:80,
                    protocol: "HTTP",
                    extra: "",
                    description: "通用机器3",
                    owner: "",
                    invalid: 0,
                    createdAt: "2018-10-16T14:30:58+08:00",
                    updatedAt: "0001-01-01T00:00:00Z"
                }
            }
        ]
    }
}
```