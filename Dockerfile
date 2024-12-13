FROM golang:alpine AS builder

ENV CGO_ENABLE 0 # 默认禁用了 cgo
ENV GOPROXY https://goproxy.cn,direct

# 安装时区相关的软件包，拷贝需要的时区到最终镜像
RUN apk update --no-cache && apk add --no-cache tzdata

# 安装 ca-certificates，拷贝证书到最终镜像
RUN apk update --no-cache && apk add --no-cache ca-certificates

# 设置工作目录为 /build
WORKDIR /build

# 将整个工程目录拷贝到工作目录下
COPY . .

# 执行编译
RUN go build -ldflags="-s -w" -o /app/apaas_ob_agent .

FROM alpine:latest

# 挂载匿名卷，防止用户将日志写到容器存储层
VOLUME ["/var/log"]

# 拷贝时区文件并设置为中国时区
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

# 拷贝证书文件
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /app

ARG RUN_NAME=apaas_ob_agent

# 拷贝 agent 的配置文件
COPY --from=builder /build/conf/config.yaml ./conf/config.yaml

# 拷贝编译后的二进制文件
COPY --from=builder /app/${RUN_NAME} ./${RUN_NAME}

CMD ["./apaas_ob_agent"]
