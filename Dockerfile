# ---------- 构建阶段 ----------
FROM golang:1.25-alpine AS builder

# 安装 git
RUN apk add --no-cache git

# 使用国内代理
ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 下载依赖
COPY go/go.mod go/go.sum ./go/
WORKDIR /app/go
RUN go mod download

# 复制源码
COPY go ./

# 进入 main.go 所在目录编译
WORKDIR /app/go/cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main .

# ---------- 运行阶段 ----------
FROM alpine:latest

# 安装证书
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 拷贝编译好的二进制
COPY --from=builder /app/main .

# 创建上传目录
RUN mkdir -p ./uploads/avatars

# 暴露端口
EXPOSE 8888

# 启动命令
CMD ["./main"]
