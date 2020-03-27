FROM golang:1.13

MAINTAINER navygong@gmail.com

ENV TZ=Asia/Shanghai

ENV path /go/src/navyt
WORKDIR ${path}
COPY . ${path}

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone \
    && go build -i -v -o navyt \
    && cp navyt /usr/bin/ \
    && rm -rf /go/pkg/navyt \
    && rm -rf /go/pkg/linux_amd64/navyt

CMD ["--help"]
