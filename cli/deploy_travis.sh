NAME="zoo-demo"

PORT="22"
USER="qwezarty"
ADDR="66.42.76.102"

# checking compiled file
[[ ! -f ./$NAME ]] && go build ./

# checking remote dirs...
ssh -p $PORT $USER@$ADDR <<- EOF
	[[ ! -d ~/$NAME ]] && mkdir ~/$NAME
	[[ ! -d ~/$NAME/engine ]] && mkdir ~/$NAME/engine
EOF

# scp executable file to remote
scp -P $PORT ./$NAME $USER@$ADDR:~/$NAME
ssh -q -p $PORT $USER@$ADDR "[[ ! -f ~/$NAME/engine/zoo.db ]]" && scp -P $PORT ./engine/zoo.db $USER@$ADDR:~/$NAME/engine

# kill existed process and run a new instance
ssh -p $PORT $USER@$ADDR <<- EOF
	[[ -n \$(pgrep $NAME) ]] && pkill $NAME
	nohup ./$NAME/$NAME >~/$NAME/$NAME.log 2>&1 &
EOF

