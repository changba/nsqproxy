### FastCGI 协议

- 协议描述

> 使用FastCGI协议 推送消息给Worker机。
> 常见的Worker机为PHP-FPM

```
消息ID：MESSAGE_ID，在FastCGI协议中发送，PHP-FPM会解析到$_SERVER['MESSAGE_ID']
消息正文：在POST请求中。
```

##### PHP使用示例

```php
//MESSAGE_ID
$_SERVER['MESSAGE_ID'];

//message body 下列三种方式均可
$_REQUEST;
$_POST;
file_get_contents("php://input");
```