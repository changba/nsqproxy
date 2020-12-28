### 删除一个消费者

- 接口功能

> 根据主键ID删除一个消费者

```
GET /admin/api/consumeConfig/delete
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
    result: "ok"
}
```