FROM golang:1.23.4-alpine3.20 as builder

WORKDIR /usr/src/app

COPY ./ ./

# 国内镜像源
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 安装运行依赖库
RUN go install && \
    go build -o app .

FROM alpine:latest

WORKDIR /app/

COPY --from=builder /usr/src/app/app .

COPY ./Docker-Service.env ./Service.env

# 运行容器时执行
CMD ["./app"]

EXPOSE 80 27017
