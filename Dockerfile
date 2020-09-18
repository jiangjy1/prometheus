# 镜像来源
FROM registry.cn-hangzhou.aliyuncs.com/billion_test/alpine-jre8:latest


# 拷贝当前目录的应用到镜像
COPY main.go /application/
COPY linux_build.sh /application/

# 声明工作目录,不然找不到依赖包，如果有的话
WORKDIR /application

# 声明动态容器卷
#VOLUME /application/logs


# 指定容器需要映射到宿主机器的端口
# 服务端口,后期可以用同一个,映射出去不同端口
EXPOSE 8080



# 启动命令
CMD ["./linux_build.sh"]
