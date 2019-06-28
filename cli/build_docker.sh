NAME="zoo-demo"
PROJECT="$GOPATH/src/github.com/qwezarty/zoo-demo"

echo -e "\033[1m==> Running tests...\033[0m"
cd $PROJECT
go test ./... | awk '{print "  --> "$0}'
[[ $? != "0" ]] && echo "  --> Exiting with tests error..." && exit 1

if [ -z $(docker images -q builds/$NAME) ]; then
	echo -e "\033[1m==> Satisfying build enviroment...\033[0m"
	cd $PROJECT/cli
	docker build -t builds/$NAME . | awk '{print "  --> "$0}'
	[[ $? != "0" ]] && echo "  --> Exiting with docker build error..." && exit 1
fi

echo -e "\033[1m==> Building docker image...\033[0m"
docker build -t $NAME . | awk '{print "  --> "$0}'

echo -e "\033[1m==> Cleaning caches...\033[0m"
docker image prune -f | awk '{print "  --> "$0}'

