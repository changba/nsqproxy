### 删除一个服务器

- 接口功能

> 根据主键ID删除一个服务器

```
GET /admin/api/workServer/delete
```

- 请求参数

|参数|必选|类型|说明|
|:----- |:-------|:-----|----- |
|id |true |int |服务器ID |
 
- 返回结果

```
{
    code: 200,
    msg: "ok",
    result: "ok"
}
```