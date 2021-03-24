![](https://img.shields.io/static/v1?label=build&message=v1.0.0&color=)
![](https://img.shields.io/static/v1?label=LICENSE&message=MIT&color=)

## 在线选课DEMO

选课系统为学校提供了基于网络，学生可以远程在线选课的系统。

**并发未测试，欢迎联系交流**

联系邮箱: yb@qxun.io

具体文档请登录演示系统的【帮助】页面查询。

**演示地址**:
- [tk.qxun.io](http://tk.qxun.io)
- 账号: domain
- 密码: tl123456

### 框架及其依赖
1. golang语言
2. [iris框架](https://github.com/kataras/iris)
3. [xorm orm框架](https://github.com/go-xorm/xorm)
4. [go-redis](https://github.com/go-redis/redis)
5. [jwt-go](https://github.com/dgrijalva/jwt-go)
6. redis，mysql数据库
7. [前端(Vue 2 + Element UI)](https://github.com/qxunio/web-take-lessons)

### 项目结构
<pre><code>
├── cmd: 程序入口 main.go
├── configs: 配置
├── controller: controller层
├── db: 数据库配置
├── domain: 领域层模型，转换
├── model: model层
├── service: service层
├── sql: 数据库初始化sql语句
├── third_party: 三方lib
├── tools: 工具包
</code></pre>

### 中间件

- mysql

- redis