FROM builds/zoo-demo AS build

RUN set -ex; \
	cd /root/go/src/github.com/qwezarty/zoo-demo; \
	git pull origin master; \
	go build ./; \
	mkdir -p /zoo-demo/engine; \
	cp ./zoo-demo /zoo-demo/ && cp ./engine/engine.db /zoo-demo/engine/;

FROM debian:stretch-slim AS dist
EXPOSE 30096
COPY --from=build /zoo-demo /zoo-demo
WORKDIR /zoo-demo
CMD ["./zoo-demo"]
