# 运行

## make file 方式
运行 ```make run``` 即可

## 裸脚本方式
运行 ```./build.sh && ./run.sh```

## docker 运行方式
1. 构造镜像：```docker build -t agent:v1.0.0 .``` 
   - 注：agent:v1.0.0 可以被替换，agent 表示镜像仓库，v1.0.0 表示镜像标签
   - 一般只需要一次构建即可，除非有改动
   - 也可以不在本地构建，从远程仓库拉（如果有远程仓库的话）
2. 运行容器：```docker run -dit --rm agent:v1.0.0``` 
3. 进入容器：
   - 查看容器ID：```docker container ls```
   - 进入容器：```docker exec -it 容器ID /bin/sh```

# 配置文件说明
## 配置文件作用
1. config.yaml：用户侧配置，用户根据实际情况修改，例如 feishu_app, feishu_secret 等
2. schema 目录下的配置文件：配置 Prometheus 指标，相对固定，用户可不关注

## 配置加载优先级
「环境变量加载」>「配置文件加载」 

环境变量枚举:
   - AGENT_SERVER_PORT：Agent 服务端口，如 8888
   - AGENT_LOG_LEVEL：Agent 的日志级别
   - AGENT_LOG_FILE_ENABLE：是否开启日志写到文件，开启：日志会存储到本机文件中
   - AGENT_LOG_FILE_FILENAME：日志文件完整路径，如 /var/log/apaas_ob_agent/run.log，前提：开启日志写到文件
   - AGENT_LOG_LOKI_ENABLE：是否开启日志写到 Loki
   - AGENT_LOG_LOKI_ROOT_URL：Loki 的根 URL，如 http://127.0.0.1:3100
   - FEISHU_APP_ID：监听的飞书应用 APP ID
   - FEISHU_APP_SECRET：监听的飞书应用 APP Secret