# syntax=docker/dockerfile:1

FROM golang:1.19-alpine
# 使用阿里云源加速下载
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk update && apk upgrade
RUN apk add --no-cache ffmpeg

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY ./ ./
RUN go mod download

# Build
RUN go build -o /app/simple-tiktok

EXPOSE 8967

# Run
CMD ["/app/simple-tiktok"]