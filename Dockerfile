# build frontend
# 前端产物与平台无关，固定在构建机平台执行（避免 QEMU 慢速构建；
# 另外 node:latest 已不再提供 linux/arm/v7，用 node:20-alpine 兼容 32 位 ARM）
FROM --platform=$BUILDPLATFORM node:20-alpine AS web_image

# 华为源
# RUN npm config set registry https://repo.huaweicloud.com/repository/npm/

WORKDIR /build

# 使用 package-lock.json（已锁定 naive-ui 2.44.x：2.45+ 的 vapor 重写与 Vue 3.2 不兼容）
COPY ./package.json ./package-lock.json /build

RUN npm ci

COPY . /build

# 前端展示用版本号（与后端 service/assets/version 保持一致）
ARG VITE_APP_VERSION=1.9.1
ENV VITE_APP_VERSION=${VITE_APP_VERSION}

RUN npm run build

# build backend
# 最新alpine3.19导致sqlite3编译失败(https://github.com/mattn/go-sqlite3/issues/1164，
# 临时解决方案:https://github.com/mattn/go-sqlite3/pull/1177)
# sun-panel暂时解决方案使用golang:1.21-alpine3.18（因旧版本使用没问题，短期内较稳定） 
FROM golang:1.21-alpine3.18 as server_image

WORKDIR /build

COPY ./service .

# 中国国内源
# RUN sed -i "s@dl-cdn.alpinelinux.org@mirrors.aliyun.com@g" /etc/apk/repositories \
#     && go env -w GOPROXY=https://goproxy.cn,direct

RUN apk add --no-cache bash curl gcc git musl-dev

RUN go env -w GO111MODULE=on \
    && export PATH=$PATH:/go/bin \
    && go install -a -v github.com/go-bindata/go-bindata/...@latest \
    && go install -a -v github.com/elazarl/go-bindata-assetfs/...@latest \
    && go-bindata-assetfs -o=assets/bindata.go -pkg=assets assets/... \
    && go build -o sun-panel --ldflags="-X sun-panel/global.RUNCODE=release -X sun-panel/global.ISDOCKER=docker" main.go



# run_image
FROM alpine

WORKDIR /app

COPY --from=web_image /build/dist /app/web

COPY --from=server_image /build/sun-panel /app/sun-panel

# 中国国内源
# RUN sed -i "s@dl-cdn.alpinelinux.org@mirrors.aliyun.com@g" /etc/apk/repositories

EXPOSE 3002

RUN apk add --no-cache bash ca-certificates su-exec tzdata \
    && chmod +x ./sun-panel \
    && ./sun-panel -config

CMD ./sun-panel
