#! /bin/bash -x
workDir=$(cd $(dirname $0); pwd)
echo 当前工作目录: $workDir
mod=ftp
target=$workDir/npacket
sshPassword=ics@beyondinfo
remoteAddr=192.168.3.11
remotePath=/kds/bin

# Go程序构建目录
goBuildDir=$workDir
echo 编译目录:$goBuildDir

# 创建构建目录
cd $workDir
go build -o npacket main.go

echo 编译完成，开始拷贝ui程序到$remoteAddr:$remotePath
# 使用sshpass直接拷贝文件到指定机器
sshpass -p $sshPassword scp $target root@$remoteAddr:$remotePath
echo 拷贝完成


