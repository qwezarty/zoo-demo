FROM golang

COPY . /go/src/github.com/qwezarty/zoo-demo

RUN set -ex; \
	export GO111MODULE=on; \
	export GOPROXY=https://goproxy.io; \
	cd $GOPATH/src/github.com/qwezarty/zoo-demo; \
	go get .; \
	go build .;

