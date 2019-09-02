name="zoo-demo"
project="$GOPATH/src/github.com/qwezarty/$name"
# export CGO_ENABLED="1"

echo -e "\033[1m==> Checking dependencies...\033[0m"
[[ -d $project ]] || { echo "  --> Exiting with project not existed error..."; exit 1; }

echo -e "\033[1m==> Running tests...\033[0m"
cd $project
go test . 1>/dev/null
[[ $? != "0" ]] && echo "  --> Exiting with tests error..." && exit 1

# Making a image for compilation if it's not existed
# Using dockerfile of $project/cli/Dockerfile
if [ -z $(docker images -q builds/$name) ]; then
	echo -e "\033[1m==> Satisfying build enviroment...\033[0m"
	cd $project
	docker build -t builds/$name -f ./cli/Dockerfile . 1>/dev/null
	[[ $? != "0" ]] && { echo "  --> Exiting with docker build error..."; exit 1; }
fi

# Starting building release image
# Using dockerfile of $project/Dockerfile
echo -e "\033[1m==> Building docker image...\033[0m"
cd $project
go build . 1>/dev/null
[[ $? != "0" ]] && { echo "  --> Exiting with go build error..."; exit 1; }
go clean
docker build -t $name . 1>/dev/null
[[ $? != "0" ]] && { echo "  --> Exiting with docker build error..."; exit 1; }

# Removing none-tag images
echo -e "\033[1m==> Cleaning caches...\033[0m"
docker image prune -f | awk '{print "  --> "$0}'

# echo -e "\033[1m==> Pushing image to docker hub...\033[0m"
# docker tag $name qwezarty/$name 1>/dev/null
# docker image prune -f 1>/dev/null
# docker push qwezarty/$name 1>/dev/null
# [[ $? != "0" ]] && { echo "  --> Exiting with image pushing error..."; exit 1; }
