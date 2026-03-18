FROM golang:1.25 AS builder

ENV TZ=Asia/Shanghai
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

# 复制所有文件到容器中
COPY . .

RUN go mod tidy
RUN go build -o main ./cmd/main.go

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
ENV service=innerg

WORKDIR /app
COPY --from=builder /app/main /app/main
COPY --from=builder /app/config /app/config

CMD ["./main"]