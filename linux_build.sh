#!/usr/bin/env bash

echo '*** Start build ...'
#set -x

# Set app version
Version="0.0.1"

# 指定要编译的平台
MacOS='darwin'
LinuxOS='linux'
WindowsOS='windows'

# Platform detection
PLATFORM='unknown'
DETECTED=$(uname | tr '[:upper:]' '[:lower:]')
if [[ "$DETECTED" == 'linux' ]]; then
   PLATFORM=${LinuxOS}
elif [[ "$DETECTED" == 'darwin' ]]; then
   PLATFORM=${MacOS}
fi

PLATFORM=${LinuxOS}
OS=${PLATFORM}

# 获取当前时间
BuildTime=`date +'%Y.%m.%d.%H%M%S'`

# 获取 Go 的版本
BuildGoVersion=`go version`

# 检查源码在最近一次 git commit 基础上，是否有本地修改，且未提交的文件
GitStatus=`git status -s`

# 获取源码最近一次 git commit log，包含 commit sha 值，以及 commit message
GitCommitLog=`git log --pretty=oneline -n 1`
# 将 log 原始字符串中的单引号替换成双引号
GitCommitLog=${GitCommitLog//\'/\"}

# 将以上变量序列化至 LDFlags 变量中
LDFlags=" \
    -X 'gitlab.bb.local/golang/bininfo.Version=${Version}' \
    -X 'gitlab.bb.local/golang/bininfo.BuildTime=${BuildTime}' \
    -X 'gitlab.bb.local/golang/bininfo.BuildGoVersion=${BuildGoVersion}' \
    -X 'gitlab.bb.local/golang/bininfo.GitStatus=${GitStatus}' \
    -X 'gitlab.bb.local/golang/bininfo.GitCommitLog=${GitCommitLog}' \
    -w \
"

ROOT_DIR=`pwd`

source ~/.bash_profile

cd ${ROOT_DIR}/ && CGO_ENABLED=0 GOOS="$OS" go build -ldflags "$LDFlags" -o ${ROOT_DIR}/http_exporter

#echo 'build done.'
echo '*** Build Done.'
