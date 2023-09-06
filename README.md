## intro
A custom protocol over TCP

这是一个简单的demo，用来解释如何在TCP的基础上实现简单的应用层协议.

协议格式：
```text
|  len  |  kind | data | 
| 2Byte | 2Byte |      |  
```
支持的消息类型
- text 
- json
- file 


## 运行
```shell
go run cmd/server/main.go   # 先运行服务器
go run cmd/client/main.go   # 使用客户端发送测试文件
```

`cmd/client/main.go` 是一个死循环，等待用户输入命令   
输入格式为：`CommandType CommandData`  
一组测试例子：
```text
text hello_1
json message_1
file /path/to/directory/test.mp4
```