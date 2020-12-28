### 创建消费者和服务器的关联关系

- 接口功能

> 创建一个消费者和服务器的关联关系

```
GET /admin/api/consumeServerMap/create
```

- 请求参数

|参数|必选|类型|说明|
|:----- |:-------|:-----|----- |
|consumeid |true |string |消费者ID |
|serverid |true |string |服务器ID |
|weight |false |int |权重，默认0 |
|invalid |false |int |是否有效，默认0有效，1无效 |

- 返回结果

```
{
    code: 200,
    msg: "ok",
    result: {
        id: 1,
        consumeid: 1,
        serverid: 1,
        weight: 1,
        invalid: 0,
        createdAt: "2020-11-30T10:45:29+08:00",
        updatedAt: "0001-01-01T00:00:00Z",
    }
}
```