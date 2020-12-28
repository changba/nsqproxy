### 删除消费者和服务器的关联关系

- 接口功能

> 根据主键ID删除一个消费者和服务器的关联关系

```
GET /admin/api/consumeServerMap/delete
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
    result: "ok"
}
```