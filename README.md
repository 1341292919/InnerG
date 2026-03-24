# InnerG Server
- 本项项目是基于前后端分离架构的实践，前端项目：[InnerG-web](https://github.com/1341292919/InnerG-web)

InnerG 的后端服务，基于 Go + Gin，提供以下核心能力：
- 用户系统（邮箱验证码、注册、登录、资料维护、头像上传）
- AI 咨询会话（会话管理 + 流式对话 SSE）
- 音乐服务（歌单与歌曲查询）
- 多存储协同（MySQL + Redis + MongoDB）
---


## 1. 技术栈

- Go `1.25.0`
- Web 框架：`gin-gonic/gin`
- ORM：`gorm`
- MySQL 驱动：`gorm.io/driver/mysql`
- Redis：`redis/go-redis/v9`
- MongoDB：`mongo-driver/v2`
- 配置：`spf13/viper`
- 鉴权：`golang-jwt/jwt/v4`（EdDSA）

---

## 2. 项目结构

```text
cmd/                # 程序入口
api/v1/             # HTTP Handler
routes/             # 路由注册与中间件
service/            # 业务逻辑层
dao/                # 数据访问层（db/cache/mongo）
config/             # 配置与初始化脚本
pkg/                # 公共组件（jwt/logger/oss/utils/errno/constants）
types/              # 请求/响应结构体
docker/             # 本地开发依赖环境（mysql/redis/mongodb）
```

请求链路：`api -> service -> dao`。

---

## 3. 本地运行

### 3.1 前置条件

请先准备：

- Go `>= 1.25`
- Docker + Docker Compose
- Linux/macOS 环境（Makefile 使用了 `sudo chown`）

### 3.2 启动依赖服务

在项目根目录执行：

```bash
make env-up
```

该命令会启动：

- MySQL（`3306`）
- Redis（`6379`）
- MongoDB（`27017`）

并挂载初始化脚本：

- MySQL：`config/sql/init.sql`
- MongoDB：`config/mongodb/init.js`


### 3.3 启动服务

```bash
go run ./cmd/main.go
```

默认监听地址来自 `config/config.yaml`

### 3.4 停止环境

```bash
make env-down
```

---

## 4. 配置说明

配置文件路径：`config/config.yaml`

### 4.1 关键配置项

- `service.address`：服务监听地址
- `service.private-key`：JWT 私钥（Ed25519 PEM）
- `mysql.*`：MySQL 连接信息
- `redis.*`：Redis 连接信息
- `mongodb.*`：MongoDB 连接信息
- `api.*`：大模型接口配置（URL、Key、Model）
- `smtp.*`：验证码邮件发送配置
- `oss.*`：对象存储配置（头像上传）
- `log.*`：业务日志/Gin 日志路径与前缀

### 4.2 配置热更新

项目通过 `viper.WatchConfig()` 监听配置变更，修改 `config/config.yaml` 后会重新映射到运行时配置。


## 5. 数据与缓存设计

- MySQL：用户、歌曲、歌单、歌单歌曲关联
- MongoDB：会话与消息历史
- Redis：验证码、token 黑名单、音乐详情缓存

---

## 6. 日志

默认日志目录：`./logs`

- 业务日志前缀：`service`
- Gin 日志前缀：`gin`


