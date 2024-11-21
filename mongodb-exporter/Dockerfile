# 第一阶段: 构建阶段
FROM golang:1.23.0 as builder

# 设置工作目录
WORKDIR /build

# 将 Go modules 文件复制到容器
COPY go.mod go.sum ./

# 下载依赖项
RUN go mod download

# 将代码复制到容器中
COPY . .

# 构建二进制文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mongodb-exporter .

# 第二阶段: 运行阶段
FROM alpine:3.17

# 安装必要的包（如 glibc，如果有动态依赖）
RUN apk --no-cache add ca-certificates

# 添加非 root 用户
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 拷贝构建的二进制文件到运行时镜像
COPY --from=builder /build/mongodb-exporter .

# 使用非 root 用户运行
USER appuser

# 暴露运行时端口
EXPOSE 8080

# 设置容器启动时的命令
CMD ["./mongodb-exporter"]
