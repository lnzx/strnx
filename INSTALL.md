# env (环境变量)

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

`HTTP_API_TOKEN` watchtower HTTP API Token

> 也可以不用docker安装，自己编译可执行文件

# 简单版本

使用docker安装

1. 需要一个postgresql数据，然后使用db.sql初始化数据库。
2. 配置环境变量
- `DATABASE_URL`  数据库连接字符串，类似这样：postgresql://USER:PASSWORD@127.0.0.1:PORT/defaultdb?sslmode=verify-full
- `KEY` 登录用的密码
- 详细参数见:env (环境变量)

# 详细版本

推荐用fly.io平台，免费提供虚拟机，数据库，个人使用，完全够用。

1. 注册fly.io平台帐号
2. 根据[官方教程](https://fly.io/docs/hands-on/install-flyctl/)，安装flyctl，然后登陆.
3. 开始部署,使用docker镜像部署,运行以下命令

```bash
fly launch --image ghcr.io/lnzx/strnx:latest
```

然后根据提示选择即可。

主要选择有部署区域，

是否要数据库，我们选Y(是)，

是否要redis，我们不需要，选N，

是否立即部署，选N，因为配置文件还要手动修改。

完成以后，当前目录会生成fly.toml文件，编辑fly.toml文件，添加环境变量[env]区域，如果使用fly.io平台集成数据库，可以不配置DATABASE_URL 环境环境变量。

```bash
app = "strns"
primary_region = "nrt"

[env]
KEY = "123456"

[build]
  image = "ghcr.io/lnzx/strnx:latest"

[http_service]
  internal_port = 8080
  force_https = true
  auto_stop_machines = false
  auto_start_machines = false
  min_machines_running = 0
```

- 初始化数库,把strns-db换成你自己的数据库app名称,一般是应用名加-db, 把strns换成数据库名，一般和应用名一样。

```bash
fly pg connect -a strns-db -d strns
```

运行完后，复制db.sql里的sql语句执行，完成后会创建3张表 \dt查看

```bash
List of relations
 Schema |  Name  | Type  |  Owner   
--------+--------+-------+----------
 public | daily  | table | postgres
 public | earn   | table | postgres
 public | wallet | table | postgres
```

然后\q退出。

最后一步，在fly.toml文件目录运行以下命令开始部署:

```bash
fly deploy
```

等待完成即可。

## 更新

在fly.tom目录再次运行： `fly deploy` 即可。