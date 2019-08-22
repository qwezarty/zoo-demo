# starting go build latest
FROM builds/zoo-demo AS build
# copy latest code to build
COPY . /go/src/github.com/qwezarty/zoo-demo
# building...
RUN set -ex; \
	export GO111MODULE=on; \
	cd $GOPATH/src/github.com/qwezarty/zoo-demo; \
	go build .; \
	mkdir -p /zoo-demo/engine; \
	cp ./zoo-demo /zoo-demo/ && cp ./engine/engine.db /zoo-demo/engine/;

# release os, this version is required and should be same as os of ffmpeg
FROM debian:buster-slim AS dist
EXPOSE 30096
COPY --from=build /zoo-demo /zoo-demo
WORKDIR /zoo-demo
CMD ["./zoo-demo"]
