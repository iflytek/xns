
## 从源码安装
```shell
#需要先准备好postgres 数据库，假设已经准备好了postgres
#1 clone 仓库源码，并编译
git clone https://git.iflytek.com/sjliu7/xns.git
cd xns
go build -o xns
#2 设置好环境变量，数据库地址需要更改，IP 地址池也要设置
export NS_SERVER_LISTEN=:4567
export NS_ADMIN_LISTEN=:8806
export NS_PG_HOST=10.1.87.70
export NS_PG_PORT=5432
export NS_PG_USER=u_xns
export NS_PG_DB_NAME=xns
export NS_PG_PASSWORD=123456
export NS_IP_SRC=resource/ip.src  # 设置IP 地址池文件，如果没有可以设置为一个空文件。

#3 启动服务
./xns


```


