# 测试游戏服务

## 网关

登陆完成后, 需要进入游戏。游戏内的业务和账号相关的业务是无关俩的, 因此需要新开一个服务去处理游戏逻辑

但是对于客户端而言, 只有一个对外地址, 所以需要引入网关

### 服务网关

> 服务网关 === 路由转发 + 过滤器

+ 路由转发: 接收一切外界请求, 转发到后端的服务
+ 过滤器: 在服务网关中可以完成一系列的横切功能, 例如权限校验、限流以及监控等

### 流程变更

1. 客户端请求地址: ws://127.0.0.0.1:8004
2. 发起登录请求时, 转发请求到登录服务器(8003)处理
3. 发起进入游戏请求时候, 转发请求到游戏服务器(8001)处理


