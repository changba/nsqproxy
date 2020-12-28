### 更新服务器

- 接口功能

> 更新一个服务器配置

```
GET /admin/api/workServer/update
```

- 请求参数

|参数|必选|类型|说明|
|:----- |:-------|:-----|----- |
|id |true |int |服务器ID |
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
    result: "ok"
}
```