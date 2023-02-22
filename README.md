# dousheng_service
抖声 service 模块（第五届字节跳动青训营后端专场）

## 快速开始

1. 配置 Go 1.19 环境
2. 下载 [ffmpeg](https://github.com/BtbN/FFmpeg-Builds/releases)，把 bin 目录下的可执行文件放到 GOPATH 的 bin 目录下
3. 配置 nacos、postgresql、redis、minio，修改 ./config 目录下的文件导入到 nacos 中，使用 ./sql 目录下的文件在 postgresql 中建表
4. 修改每个服务监听的 nacos 地址
5. 运行 ./build.sh
6. 运行 ./run.sh

## 服务架构

![img](D:\GolandProjects\dousheng_service\img\img.png)
