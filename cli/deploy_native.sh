name="zoo-demo"

port="22"
user="qwezarty"
addr="66.42.76.102"

GOPATH="/home/qwezarty/Code/go"
project="$GOPATH/src/github.com/qwezarty/$name"

# checking remote dirs...
ssh -qt -p $port $user@$addr <<- EOF
	export GOPATH=$GOPATH
	[[ ! -d ~/$name ]] && mkdir ~/$name
	[[ ! -d ~/$name/engine ]] && mkdir ~/$name/engine

	if [ -d $project ]; then
		cd $project && git pull
	else
		go get -u github.com/qwezarty/$name
	fi

	cd $project
	go build . 
	[[ $? != "0" ]] && { echo "Exiting with building error"; exit 1; }
	[[ -n \$(pgrep -u \$(whoami) $name) ]] && pkill -u \$(whoami) $name
	cp ./$name ~/$name && cp ./engine/engine.db ~/$name/engine

	cd ~/$name
	nohup ./$name >./$name.log 2>&1 &
EOF
