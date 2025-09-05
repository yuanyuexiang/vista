FROM registry.cn-zhangjiakou.aliyuncs.com/matrix-net/ubuntu:24.04
LABEL maintainer "matrix@126.com"

ENV TZ='Asia/Shanghai'

# 设置阿里云源并安装依赖，安装必要依赖和 CA 证书
RUN sed -i 's|http://archive.ubuntu.com/ubuntu/|http://mirrors.aliyun.com/ubuntu/|g' /etc/apt/sources.list \
 && apt-get update \
 && apt-get install -y --no-install-recommends \
        ca-certificates curl \
 && rm -rf /var/lib/apt/lists/* \
 && update-ca-certificates

WORKDIR /matrix
COPY vista ./
COPY config.yaml ./config.yaml
RUN chmod +x ./vista
CMD ["/matrix/vista"]