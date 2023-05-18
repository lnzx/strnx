# strnx

主要用来监控钱包，查看每日收益，节点运行状态, 适合很多钱包的用户。

默认端口8080,无法更改。

## env (环境变量)

`DATABASE_URL` 

e.g.

`postgresql://USER:PASSWORD@127.0.0.1:PORT/defaultdb?sslmode=[disable|verify-full]`

`KEY` 密码,默认123456

可选参数，发送短信相关配置:

`SMS_API_KEY` 发送短信平台的api key,另一个项目。

`MOBILE` 接受短信手机号，如: 13988888888

可选参数，节点相关配置:

`SSH_USER` SSH默认用户名

`SSH_PASS` SSH默认密码

`HTTP_API_TOKEN`

# 项目结构
### web目录 (前端页面 vite+vue3+bulma css)

# 功能
### 节点管理未实现

# demo
strns.fly.dev 

密码: 123456 用户名随便

# 详细安装教程:

https://lnzx.notion.site/lnzx/strnx-86737b747d5e4187ad3195624bf80414

# 计划

1. 加入权限
2. 增加速率限制
3. 优化页面大小，按需引用
4. 完成节点管理：增加节点，一键升级。