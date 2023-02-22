cd $(dirname $0)

nohup sh ./gateway > ./gateway.log 2>&1 &

nohup sh ./message > ./message.log 2>&1 &

nohup sh ./user > ./user.log 2>&1 &

nohup sh ./video > ./video.log 2>&1 &