## 使用docker 快速部署
xns 
#### 1 部署数据库

````
#部署数据库
docker run -tid --name=postgres_nameserver -v /data/postgres/nameserver:/var/lib/postgresql/data -p "5432:5432" -e "POSTGRES_USER=u_nameserver" -e "POSTGRES_DB=nameserver"  -e "POSTGRES_PASSWORD=123456" postgres:9.5

#初始化数据库
#部署

````

#### 2
