FROM golang:1.20 as builder

COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn make build

FROM debian:stable-slim

RUN mkdir -p /app/configs

COPY --from=builder /src/_output/yoo /app
COPY --from=builder /src/configs/yoo.yaml /app/configs

WORKDIR /app

EXPOSE 8080

VOLUME ["/data/configs"]

CMD ["./yoo", "-c", "/data/configs/yoo.yaml"]