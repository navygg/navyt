FROM registry.matrix.netease.com/dev/gitlabci:1.12

MAINTAINER hjgong@corp.netease.com

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV path /go/src/navyt
WORKDIR ${path}
COPY . ${path}

RUN go build -i -v -o navyt \
    && cp navyt /usr/bin/ \
    && rm -rf /go/pkg/navyt \
    && rm -rf /go/pkg/linux_amd64/navyt

CMD ["--help"]
