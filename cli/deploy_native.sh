NAME="zoo-demo"

PORT="22"
USER="qwezarty"
ADDR="66.42.76.102"

GOPATH="/home/qwezarty/Code/go"

# checking remote dirs...
ssh -p $PORT $USER@$ADDR <<- EOF
	[[ ! -d ~/$NAME ]] && mkdir ~/$NAME
	[[ ! -d ~/$NAME/engine ]] && mkdir ~/$NAME/engine
EOF

# clone or pull code
ssh -p $PORT $USER@$ADDR <<- EOF
	go get -u github.com/qwezarty/zoo-demo
EOF

ssh -p $PORT $USER@$ADDR <<- EOF
	cd $GOPATH/src/github.com/qwezarty/zoo-demo && go build ./ && cp ./$NAME ~/$NAME && cp ./engine/zoo.db ~/$NAME/engine
EOF

ssh -p $PORT $USER@$ADDR <<- EOF
	cd ~/$NAME
	[[ -n \$(pgrep $NAME) ]] && pkill $NAME
	nohup ./$NAME >./$NAME.log 2>&1 &
EOF

