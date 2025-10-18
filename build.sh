#!/bin/bash
echo '开始拉取最新代码'
git reset --hard origin/master
git pull origin master
echo '打包镜像'
DOCKER_BUILDKIT=1 docker build -t go-movie:latest .
echo '启动容器'
docker-compose down
docker-compose up --force-recreate -d --remove-orphans