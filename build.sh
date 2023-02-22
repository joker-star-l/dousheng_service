cd $(dirname $0)

cd ./gateway
go build -o ../gateway_build
cd ../

cd ./message
go build -o ../message_build
cd ../

cd ./user
go build -o ../user_build
cd ../

cd ./video
go build -o ../video_build
cd ../