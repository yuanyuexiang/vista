FROM registry.cn-zhangjiakou.aliyuncs.com/matrix-net/ubuntu:24.04
LABEL maintainer "matrix@126.com"

ENV TZ='Asia/Shanghai'

# 设置阿里云源并安装依赖，安装必要依赖和 CA 证书
RUN echo "deb http://mirrors.aliyun.com/ubuntu/ noble main restricted universe multiverse\n\
deb http://mirrors.aliyun.com/ubuntu/ noble-updates main restricted universe multiverse\n\
deb http://mirrors.aliyun.com/ubuntu/ noble-backports main restricted universe multiverse\n\
deb http://mirrors.aliyun.com/ubuntu/ noble-security main restricted universe multiverse" > /etc/apt/sources.list \
 && rm -f /etc/apt/sources.list.d/* \
 && apt-get clean \
 && apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates curl \
 && rm -rf /var/lib/apt/lists/* \
 && update-ca-certificates

WORKDIR /matrix
COPY vista ./
COPY config.yaml ./config.yaml
RUN chmod +x ./vista

COPY frontend_demo.html ./
COPY test_wechat.html ./

CMD ["/matrix/vista"]