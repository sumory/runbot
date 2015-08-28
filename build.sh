#!/bin/sh

#####################################################################
# usage:
# sh build.sh 构建默认的linux64位程序
# sh build.sh darwin(或linux), 构建指定平台的64为程序
#####################################################################

source /etc/profile



OS="$1"
if [ -n "$OS" ];then
   echo "use defined GOOS: "$OS
else
   echo "use default GOOS: linux"
   OS=linux
fi

echo "start building with GOOS: "$OS

export GOOS=$OS
export GOARCH=amd64


release_dir="release"


mkdir -p ./${release_dir}
rm -rf ./${release_dir}/*


go build -ldflags "-X main.buildstamp `date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash `git rev-parse HEAD`" -x -o ${release_dir}/runbot runbot.go

cp ./config.conf ./${release_dir}/
cp ./start.sh ./${release_dir}/
cp ./readme.md ./${release_dir}/
