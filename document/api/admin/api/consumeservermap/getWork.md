### 查询一个消费者和服务器的关联关系

- 接口功能

> 根据主键ID查询一个消费者和服务器的关联关系

```
GET /admin/admin/consumeServerMap/getWork
```

- 请求参数

|参数|必选|类型|说明|
|:----- |:-------|:-----|----- |
|id |true |int |关联关系ID |
 
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
        workServer: {
            id: 1,
            addr: "0.0.0.0:80",
            protocol: "HTTP",
            extra: "test.php",
            description: "",
            owner: "",
            invalid: 0,
            createdAt: "2020-11-27T11:04:31+08:00",
            updatedAt: "0001-01-01T00:00:00Z"
        }
    }
}
```