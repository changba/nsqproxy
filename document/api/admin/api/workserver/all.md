### 批量查询服务器

- 接口功能

> 批量查询服务器

```
GET /admin/api/workServer/all
```

- 请求参数

> 无
 
- 返回结果

```
{
    code: 200,
    msg: "ok",
    result: [
        {
            id: 1,
            addr: "10.10.1.1:80",
            protocol: "HTTP",
            extra: "index.php",
            description: "通用机器",
            owner: "",
            invalid: 0,
            createdAt: "2018-10-16T14:31:08+08:00",
            updatedAt: "0001-01-01T00:00:00Z"
        },
        {
            id: 2,
            addr: "10.10.1.1:9000",
            protocol: "FASTCGI",
            extra: "/home/wwwroot/test.php",
            description: "通用机器",
            owner: "",
            invalid: 0,
            createdAt: "2018-10-16T14:31:08+08:00",
            updatedAt: "0001-01-01T00:00:00Z"
        }
    ]
}
```