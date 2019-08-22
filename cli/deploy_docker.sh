name="zoo-demo"
project="$GOPATH/src/github.com/qwezarty/zoo-demo"

port="22"
user="qwezarty"
addr="66.42.76.102"

echo -e "\033[1m==> Checking image...\033[0m"
if [ -z $(docker images -q $name) ]; then
	echo "  --> IMAGE NOT FOUND!"
	echo "  --> Forget to run cli/build_docker.sh first?"
	echo "  --> Exiting with error..."
	exit 1
fi

echo -e "\033[1m==> Resolving remote dependencies...\033[0m"
if test ! -f $HOME/.ssh/id_rsa.pub; then
	echo "no ssh key found locally, I'll generate for you"
	ssh-keygen -q -t rsa -b 2048
fi
ssh-copy-id -p $port $user@$addr &>/dev/null
ssh -qt -p $port $user@$addr <<- EOF 1>/dev/null
	# create dirs and files
	[[ ! -d ~/$name ]] && mkdir ~/$name
	[[ ! -d ~/$name/engine ]] && mkdir ~/$name/engine
	# checking if docker exists
	if ! hash docker 2>/dev/null; then
		. /etc/os-release
		if [ \$NAME != "Ubuntu" ]; then
			echo "docker is not installed at remote server(and not a ubuntu server), you should follow docker official docs to install!" >&2
			exit 1
		fi
		sudo apt-get update 
		sudo apt-get install -y docker.io
	fi
	# add current user to docker group
	if groups \$USER | grep &>/dev/null '\bdocker\b'; then
		sudo usermod -aG docker \$USER
	fi
EOF
[[ $? != "0" ]] && { echo "  --> Exiting with dependencies solving error..."; exit 1; }

echo -e "\033[1m==> Compressing...\033[0m"
cd $project
docker save -o ./$name.tar $name

echo -e "\033[1m==> Sending image to remote...\033[0m"
scp -P $port ./zoo-demo.tar $user@$addr:~/$name
[[ $? != "0" ]] && { echo "  --> Exiting with scp error..."; exit 1; }
scp -P $port ./engine/engine.db $user@$addr:~/$name/engine/
[[ $? != "0" ]] && { echo "  --> Exiting with scp error..."; exit 1; }

echo -e "\033[1m==> Loading image...\033[0m"
ssh -qt -p $port $user@$addr <<- EOF 1>/dev/nul
	if [ ! -z \$(docker ps -aq --filter ancestor="$name") ]; then
		docker rm -f \$(docker ps -aq --filter ancestor="$name")
	fi
	docker load < ~/$name/$name.tar;
	docker image prune -f
EOF
[[ $? != "0" ]] && { echo "  --> Exiting with loading error..."; exit 1; }

echo -e "\033[1m==> Starting new container...\033[0m"
ssh -qt -p $port $user@$addr <<- EOF 1>/dev/null
	docker run --restart always -d -p 30097:30096 \
	   -v ~/$name/engine/engine.db:/$name/engine/engine.db \
		$name
EOF

echo -e "\033[1m==> Cleaning caches...\033[0m"
rm -f ./$name.tar

