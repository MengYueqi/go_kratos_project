# Go Kratos Review Service

本项目使用 [Kratos](https://go-kratos.dev/) 微服务框架开发，作为电商评价系统的一部分。基于 Q1mi 课程进行实现。

---

## 🗓️ 迭代进度
- **第一次迭代 DDL：2025-09-15**

---

## 🚀 服务端口

| 服务/模块       | 协议/功能         | 端口号 (宿主机) | 说明                 |
|----------------|------------------|----------------|--------------------|
| **MySQL**      | 数据库            | 3307           | 关系型数据库             |
| **Redis**      | 缓存/消息队列      | 6380           | 使用 5.0.7 版本        |
| **review-server** | HTTP / gRPC   | 8482 / 9492    | 业务中台服务             |
| **review-b**   | HTTP / gRPC      | 8483 / 9493    | 商家侧服务              |
| **review-o**   | HTTP / gRPC      | 8484 / 9494    | 审核侧服务              |
| **review-job** | HTTP / gRPC      | 8485 / 9495    | 异步任务服务（Kafka → ES） |
| **Consul**     | Web UI / DNS     | 8500 / 8600udp | 服务注册与发现            |
| **Canal**      | TCP              | 11111          | MySQL binlog 订阅    |
| **Kafka**      | Broker / 内外部   | 9092 / 29092   | 消息队列服务             |
| **Kafka UI**   | Web UI           | 8080           | Kafka 管理界面         |
| **Elasticsearch** | REST / Transport | 9200 / 9300 | ES API 与集群通信       |
| **Kibana**     | Web UI           | 5601           | 可视化界面              |

---

## ⚙️ 快速启动

### 1. 启动基础依赖
```bash
# 启动 MySQL
docker run --name mysql -p 3307:3306 -e MYSQL_ROOT_PASSWORD=root1234 -d mysql

# 启动 Redis
docker run --name redis507 -p 6380:6379 -d redis:5.0.7

# 启动 Consul
docker run -d --name=consul-dev \
  -p 8500:8500 \
  -p 8600:8600/udp \
  hashicorp/consul agent -dev -client=0.0.0.0
```
其余 Kafka、Elasticsearch、Kibana 等服务可使用提供的 docker-compose.yml 启动。

### 2. 启动服务
进入对应服务目录后执行：
```bash
kratos run
```

## ✅ 验证运行

- 业务中台：http://127.0.0.1:8482
- 商家侧：http://127.0.0.1:8483
- 审核侧：http://127.0.0.1:8484
- Kafka UI：http://127.0.0.1:8080
- Kibana：http://127.0.0.1:5601
- Consul：http://127.0.0.1:8500