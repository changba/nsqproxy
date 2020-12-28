### 查询一个服务器

- 接口功能

> 根据主键ID查询一个服务器

```
GET /admin/api/workServer/get
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
    result: {
        id: 1,
        addr: "10.10.1.1:80",
        protocol: "HTTP",
        extra: "index.php",
        description: "通用机器",
        owner: "",
        invalid: 0,
        createdAt: "2018-10-16T14:31:08+08:00",
        updatedAt: "0001-01-01T00:00:00Z"
    }
}
```