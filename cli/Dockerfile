FROM debian:stretch AS build

RUN set -ex; \
	apt-get update; \
	apt-get install -y wget git build-essential; \
	cd /tmp; \
	wget https://dl.google.com/go/go1.12.6.linux-amd64.tar.gz; \
	tar -C /usr/local -xzf go1.12.6.linux-amd64.tar.gz; \
	export PATH=$PATH:/usr/local/go/bin; \
	export HTTP_PROXY='http://192.168.15.19:8001'; \
	export HTTPS_PROXY='http://192.168.15.19:8001'; \
	go get -u -v github.com/qwezarty/zoo-demo; \
	cd /root/go/src/github.com/qwezarty/zoo-demo; \
	go build ./;

ENV GOPATH /root/go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
