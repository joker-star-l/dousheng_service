cd $(dirname $0)

nohup sh ./gateway_build > ./gateway.log 2>&1 &

nohup sh ./message_build > ./message.log 2>&1 &

nohup sh ./user_build > ./user.log 2>&1 &

nohup sh ./video_build > ./video.log 2>&1 &