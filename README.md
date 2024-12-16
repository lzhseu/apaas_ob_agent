# 运行

## make file 方式
运行 ```make run``` 即可

## 裸脚本方式
运行 ```./build.sh && ./run.sh```

## docker 运行方式
1. 构造镜像：```docker build -t agent:v1.0.0 .``` 
   - 注：agent:v1.0.0 可以被替换，agent 表示镜像仓库，v1.0.0 表示镜像标签
   - 一般只需要一次构建即可，除非有改动
   - 也可以不在本地构建，从远程仓库拉（如果有远程仓库的话）
2. 运行容器：```docker run -dit --rm agent:v1.0.0``` 
3. 进入容器：
   - 查看容器ID：```docker container ls```
   - 进入容器：```docker exec -it 容器ID /bin/sh```