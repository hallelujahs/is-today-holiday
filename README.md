# 今天是休息日吗？

在使用 Apple Shortcut（捷径）时，经常会有查询一下今天是不是休息日，以继续执行与工作日不同的操作。

找了下现成的一些公开的接口，经常会出现超时或者收费的情况，所以简单写了一个查询当前是否为工作日的接口。

## 如何使用

在使用之前，请保证你的服务器上已经安装 docker 和 docker-compose 等工具。

使用该服务，最好搭建一个统一的 OpenAPI 网关，其他的应用接口作为 Upstream 挂在该网关后面，用 OpenAPI 网关提供鉴权、限流等。

``` shell
# 这里使用 app.net 网络，主要是为了让 OpenAPI 网关访问
docker network create app.net

mkdir -p /data/log/app/is-today-holiday
mkdir -p /data/app/is-today-holiday

cp -r configs /data/app/is-today-holiday
cp deployments/docker-compose.yaml /data/app/is-today-holiday

cd /data/app/is-today-holiday
docker-compose up -d
```
