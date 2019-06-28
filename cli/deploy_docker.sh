NAME="zoo-demo"
PROJECT="$GOPATH/src/github.com/qwezarty/zoo-demo"

PORT="22"
USER="qwezarty"
ADDR="66.42.76.102"

echo -e "\033[1m==> Checking image...\033[0m"
if [ -z $(docker images -q $NAME) ]; then
	echo "  --> IMAGE NOT FOUND!"
	echo "  --> Forget to run cli/build_docker.sh first?"
	echo "  --> Exiting with error..."
	exit 1
fi

echo -e "\033[1m==> Compressing...\033[0m"
cd $PROJECT
docker save -o ./$NAME.tar $NAME

echo -e "\033[1m==> Sending image to remote...\033[0m"
ssh -p $PORT $USER@$ADDR -q <<- EOF > /dev/null
	[[ ! -d ~/$NAME ]] && mkdir ~/$NAME
	[[ ! -d ~/$NAME/engine ]] && mkdir ~/$NAME/engine
EOF
scp -P $PORT ./zoo-demo.tar $USER@$ADDR:~/$NAME
[[ $? != "0" ]] && echo "  --> Exiting with scp error..." && exit 1
scp -P $PORT ./engine/engine.db $USER@$ADDR:~/$NAME/engine/
[[ $? != "0" ]] && echo "  --> Exiting with scp error..." && exit 1


echo -e "\033[1m==> Loading image...\033[0m"
ssh -p $PORT $USER@$ADDR -q <<- EOF | awk '{print "  --> "$0}'
	[[ ! -z \$(docker ps -aq --filter ancestor="$NAME") ]] && docker rm -f \$(docker ps -aq --filter ancestor="$NAME"); \
	docker load < ~/$NAME/$NAME.tar;
EOF
[[ $? != "0" ]] && echo "  --> Exiting with loading error..." && exit 1

echo -e "\033[1m==> Starting new container...\033[0m"
ssh -p $PORT $USER@$ADDR -q <<- EOF | awk '{print "  --> "$0}'
	docker run --restart always -d -p 30097:30096 \
	   -v ~/$NAME/engine/engine.db:/$NAME/engine/engine.db \
		$NAME
EOF

echo -e "\033[1m==> Cleaning caches...\033[0m"
ssh -p $PORT $USER@$ADDR -q "docker image prune -f" | awk '{print "  --> "$0}'
rm -f ./$NAME.tar

