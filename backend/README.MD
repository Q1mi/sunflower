# GoFrame Template For SingleRepo

Goframe 框架官方文档: 
- https://goframe.org/quick

## 视频教程
[B站视频教程](https://www.bilibili.com/cheese/play/ss167267536)

## 启动后端
** 重要：在执行后续命令前，确保你在 `backend` 目录下。**

### 环境准备

1. 在MySQL中创建数据表

连接你的 MySQL 数据库，执行：`create_tables.sql`

2. 修改 `manifest/config/config.yaml` 中的 MySQL 和 Redis 连接信息

### 安装依赖
```bash
go mod tidy
```

### 启动项目
```bash
go run .
```
