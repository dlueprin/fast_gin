FROM golang:alpine AS builder

ENV CGO_ENABLE=0
ENV GOPROXY https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

WORKDIR /build
ADD . .
RUN go build -o main

FROM alpine

WORKDIR /app

#把一阶段构建的镜像的main复制到二阶段镜像的app目录下
COPY --from=builder /build/main /app
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
#安装时区数据库包。使得能识别转化不同时区时间
RUN apk add tzdata

CMD ["./main"]