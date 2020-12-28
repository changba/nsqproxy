### 更新消费者和服务器的关联关系

- 接口功能

> 更新一个消费者和服务器的关联关系

```
GET /admin/admin/consumeServerMap/update
```

- 请求参数

|参数|必选|类型|说明|
|:----- |:-------|:-----|----- |
|id |true |int |关联关系ID |
|consumeid |true |string |消费者ID |
|serverid |true |string |服务器ID |
|weight |false |int |权重，默认0 |
|invalid |false |int |是否有效，默认0有效，1无效 |

- 返回结果

```
{
    code: 200,
    msg: "ok",
    result: "ok"
}
```