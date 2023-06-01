# 基础镜像
FROM golang:1.20.1 AS build

# 设置工作目录
WORKDIR /app

# 复制 Go 项目代码到容器内
COPY . .
ENV GOPROXY https://goproxy.cn,direct
# 构建可执行文件
RUN go build -o main ./cmd/myapp

# 运行时镜像
#FROM alpine:latest
FROM busybox:latest

# 设置工作目录
RUN mkdir /config && mkdir /download
WORKDIR /app

# 从 build 阶段复制可执行文件
COPY --from=build /app/main .

# 设置容器启动命令
#CMD ["./myapp"]
