cd $(dirname $0)

nohup ./gateway_build > ./gateway.log 2>&1 &

nohup ./message_build > ./message.log 2>&1 &

nohup ./user_build > ./user.log 2>&1 &

nohup ./video_build > ./video.log 2>&1 &