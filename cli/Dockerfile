FROM golang

COPY . /go/src/github.com/qwezarty/zoo-demo

RUN set -ex; \
	go env -w GO111MODULE=on; \
	go env -w GOPROXY="https://goproxy.io,direct"; \
	cd /go/src/github.com/qwezarty/atomsrv; \
	go mod init; \
	go get -u -v ./...; \
	go build .;

