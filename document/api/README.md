# API

## 通用响应值

|返回字段|字段类型|说明|
|:----- |:------|:-----------------------------|
|code   |int    |成功为200，其他为失败|
|msg    |string |成功为ok，失败时为原因|
|result |string |响应正文|

## 服务相关

* 查看状态：[service/status](service/status.md)
* 查看角色：[service/getrole](service/getrole.md)
* pprof：[debug/pprof](debug/pprof.md)

## 后台相关

### 服务器管理

* 新增：[admin/api/workServer/create](admin/api/workserver/create.md)
* 修改：[admin/api/workServer/update](admin/api/workserver/update.md)
* 删除一个：[admin/api/workServer/delete](admin/api/workserver/delete.md)
* 查询一个：[admin/api/workServer/get](admin/api/workserver/get.md)
* 查询所有：[admin/api/workServer/page](admin/api/workserver/page.md)

### 消费者管理

* 新增：[admin/api/consumeConfig/create](admin/api/consumeconfig/create.md)
* 修改：[admin/api/consumeConfig/update](admin/api/consumeconfig/update.md)
* 删除一个：[admin/api/consumeConfig/delete](admin/api/consumeconfig/delete.md)
* 查询一个：[admin/api/consumeConfig/get](admin/api/consumeconfig/get.md)
* 查询所有：[admin/api/consumeConfig/page](admin/api/consumeconfig/page.md)
* 查询一个消费者关联服务器：[admin/api/consumeConfig/workList](admin/api/consumeConfig/worklist.md)

### 消费者和服务器关联关系管理

* 新增：[admin/api/consumeServerMap/create](admin/api/consumeservermap/create.md)
* 修改：[admin/api/consumeServerMap/update](admin/api/consumeservermap/update.md)
* 删除一个：[admin/api/consumeServerMap/delete](admin/api/consumeservermap/delete.md)
* 查询一个（含关联服务器信息）：[admin/api/consumeServerMap/getWork](admin/api/consumeservermap/getWork.md)