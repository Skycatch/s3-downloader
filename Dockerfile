FROM golang:alpine

RUN apk add --no-cache git

ENV AWS_REGION=
ENV AWS_SECRET_ACCESS_KEY=
ENV AWS_ACCESS_KEY_ID=

COPY . /go/src/downloader
RUN cd /go/src/downloader && go get && go build -o tasque-downloader *.go
# RUN sysctl net.ipv4.tcp_tw_recycle=1 && sysctl net.ipv4.tcp_tw_reuse=1

CMD ["/go/src/downloader/tasque-downloader", "get"]
