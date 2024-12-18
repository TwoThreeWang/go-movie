#!/bin/bash
echo '开始拉取最新代码'
git reset --hard origin/master
git pull origin master
echo '打包镜像'
docker build -t go-movie:latest .
echo '停止并删除旧容器'
docker rm -f go-movie
echo '启动容器'
docker run --name go-movie -d -v ./config.yaml:/config.yaml -p 8787:8787 -v ./templates:/templates -v ./public:/public:rw -v ./data:/data:rw --log-opt max-size=10m --restart=always go-movie:latest
echo '清理不再使用的镜像、容器和数据卷'
docker system prune --all --force --volumes