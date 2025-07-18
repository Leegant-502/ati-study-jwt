# 使用多阶段构建
# 构建阶段
FROM golang:alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go mod和sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 运行阶段
FROM alpine:latest

# 安装基本工具
RUN apk --no-cache add tzdata

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .
COPY --from=builder /app/config/config.yaml ./config/

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]
