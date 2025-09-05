FROM registry.cn-zhangjiakou.aliyuncs.com/matrix-net/ubuntu:24.04
LABEL maintainer "matrix@126.com"

ENV TZ='Asia/Shanghai'

WORKDIR /matrix
COPY vista ./
COPY config.yaml ./config.yaml
RUN chmod +x ./vista
CMD ["/matrix/vista"]