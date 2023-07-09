#!/bin/bash
echo ""
#输出当前时间
date --date='0 days ago' "+%Y-%m-%d %H:%M:%S"
echo "Start"
#git分支名称
branch="master"
echo $branch
echo $1
#git项目路径
cd /www/wwwroot/go_movie_git && \
git pull  && \
echo "处理完成"
exit 0