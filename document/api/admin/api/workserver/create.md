### 创建服务器

- 接口功能

> 创建一个新的服务器

```
GET /admin/api/workServer/create
```

- 请求参数

|参数|必选|类型|说明|
|:----- |:-------|:-----|----- |
|addr |true |string |地址，IP:PORT |
|protocol |true |string |协议，如HTTP、FastCGI、CBNSQ |
|extra |false |string |扩展字段，默认空 |
|description |false |string |描述，默认空 |
|owner |false |string |责任人，默认空 |
|invalid |false |int |是否有效，默认0有效，1无效 |

- extra参数 详解
  - 当protocol为HTTP时，extra表示URL的路径，即http://addr/extra
  - 当protocol为FastCGI时，extra表示要执行的文件名路径，即/home/wwwroot/app/index.php

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