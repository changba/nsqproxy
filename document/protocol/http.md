### HTTP 协议

- 协议描述

> 使用HTTP POST协议 推送消息给Worker机。Content-Type为application/x-www-form-urlencoded

```
消息ID：MESSAGE_ID，在HTTP Header中。
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

##### Golang使用示例

```golang
//MESSAGE_ID
messageId := r.Header.Get("MESSAGE_ID")

//方法一 - []byte类型的消息体
data, _ := ioutil.ReadAll(r.Body)

//方法二 - KV格式
value := r.PostFormValue("key")
//或
value = r.FormValue("key")

//方法三 - KV格式
r.ParseForm()
value = r.Form.Get("key")
//或
value = r.Form["key"][0]
```