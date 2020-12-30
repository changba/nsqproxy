## CBNSQ 协议

- 协议描述

> 自定义的基于TCP的文本协议。

- CBNSQ协议是自定义的非常简单的一个基于文本的TCP的协议。
- 就是消息长度 + 消息ID + 消息内容。 msg的长度(8字节，0填充） + messageId（16字节） + msg正文
- 如：有个用户注册的消息是：{"topic":"test","classname":"UserService","methodname":"addUser","param":["userid", "username", "password"],"addtime":"2020-11-27 14:30:34"}
- 第一部分：消息长度。消息主体是个json，长度140，补齐8位，那么这部分为00000140。8位的极限是99999999/1024/1024≈95M，足够了吧。
- 第二部分：消息ID，这个是NSQ的消息唯一ID，16位。如qwertyuiopasdfgh。注意：这16位是不计入第一部分的消息长度的。
- 第三部分：消息主体。即刚才提到的JSON串。
- 完整的消息为：00000140qwertyuiopasdfgh{"topic":"test","classname":"UserService","methodname":"addUser","param":["userid", "username", "password"],"addtime":"2020-11-27 14:30:34"}

```
消息ID：消息体（从0开始）的第8位到第23位。
消息正文：消息体（从0开始）的第24位到[消息长度]，[消息长度]为消息体的第0位到第7位。
```

### PHP使用示例

```php
<?php
$msg = '00000140qwertyuiopasdfgh{"topic":"test","classname":"UserService","methodname":"addUser","param":["userid", "username", "password"],"addtime":"2020-11-27 14:30:34"}';
//message length
$length = intval(substr($msg, 0, 8));
if(strlen($msg) < 8 + 16 + $length){
    exit('Incomplete message');
}

//MESSAGE_ID
$messageId = substr($msg, 8, 16);

//message body
$body = substr($msg, 8+16, $length);
```

### MeepoPS使用示例

[MeepoPS](https://github.com/lixuancn/MeepoPS) 是PHP的服务端程序，监听端口后，与客户端进行通信。支持CBNSQ协议。

```php
<?php
//引入MeepoPS
require_once 'MeepoPS/index.php';

//使用文本协议传输的Api类
$cbNsq = new \MeepoPS\Api\Cbnsq('0.0.0.0', '19910');

//启动的子进程数量. 通常为CPU核心数
$cbNsq->childProcessCount = 10;

//设置MeepoPS实例名称
$cbNsq->instanceName = 'MeepoPS-CBSNQ';

//设置回调函数 - 这是所有应用的业务代码入口
$cbNsq->callbackNewData = 'callbackNewData';

//启动MeepoPS
\MeepoPS\runMeepoPS();


//以下为回调函数, 业务相关.
//回调 - 收到新消息
function callbackNewData($connect, $data)
{
    file_put_contents('../cbnsq.log', date('Y-m-d H:i:s') . ' ' . $_SERVER['MESSAGE_ID'] . ' ' . $data  ."\n", FILE_APPEND);
    $connect->send("200 ok");
}
```