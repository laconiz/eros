# eros

封装Golang服务器的各种包，尚在开发中

## database/
* consul  
封装github.com/hashicorp/consul, 更方便地调用consul
* elastic  
封装github.com/olivere/elastic, 更方便地使用ElasticSearch
* redis  
封装github.com/gomodule/redigo/redis, 更方便地使用Redis

## iris/
一套基于oceanus的游戏服务器示例

## logis/
日志包

## network
* cipher  
长连接数据流加密
* encoder  
网络消息序列化
* socket  
tcp协议的服务器与客户端, 统一长连接接口
* steropes  
封装github.com/gin-gonic/gin, 更方便地生成HTTP服务
* websocket  
封装github.com/gorilla/websocket, 统一长连接接口

## oceanus/
actor模型与service mesh杂交的服务器框架

## utils/
* command  
命令行参数处理
* ioc  
依赖注入器
* json  
github.com/json-iterator/go
* mathe  
封装一些数学库
