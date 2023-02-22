cd $(dirname $0)

cd ./gateway
go build -o ../gateway
cd ../

cd ./message
go build -o ../message
cd ../

cd ./user
go build -o ../user
cd ../

cd ./video
go build -o ../video
cd ../