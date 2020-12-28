### 查看角色

- 接口功能

> 查看当前是主服务还是备服务

```
GET /getRole
```

- 请求参数

> 无
 
- 返回结果

> 结果有master或slave

```
{
    code: 200,
    msg: "ok",
    result: "master"
}
 ```