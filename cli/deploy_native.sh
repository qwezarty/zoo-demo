NAME="zoo-demo"

PORT="22"
USER="qwezarty"
ADDR="66.42.76.102"

GOPATH="/home/qwezarty/Code/go"
PROJECT="$GOPATH/src/github.com/qwezarty/$NAME"

# checking remote dirs...
ssh -p $PORT $USER@$ADDR <<- EOF
	export GOPATH=$GOPATH
	[[ ! -d ~/$NAME ]] && mkdir ~/$NAME
	[[ ! -d ~/$NAME/engine ]] && mkdir ~/$NAME/engine

	if [ -d $PROJECT ]; then
		cd $PROJECT && git pull
	else
		go get -u github.com/qwezarty/$NAME
	fi

	cd $PROJECT && go build ./ && cp ./$NAME ~/$NAME && cp ./engine/engine.db ~/$NAME/engine

	cd ~/$NAME
	[[ -n \$(pgrep $NAME) ]] && pkill $NAME
	nohup ./$NAME >./$NAME.log 2>&1 &
EOF
