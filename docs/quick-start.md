
#### 开始使用
完整的管理admin api文档: 启动服务后浏览器访问admin 端口的 /docs 路由。如： http://127.0.0.1:8806/docs ，管理端口默认未开启权限控制，如果要开启权限控制，则可以通过新增环境变量
NS_LOGIN_ENABLED=true 来开启。开启后admin端口会需要通过登陆来访问

[登陆文档](./user_login.md)


### 1 快速配置解析规则

```shell


admin=127.0.0.1:8806
#1 创建机房
curl -X POST ${admin}/idcs -H 'Content-type:application/json' -d '{"name":"hd"}'

# 创建服务器组
curl -X POST ${admin}/groups -H 'Content-type:application/json' -d '
{
"name":"g-hd",
"default_servers":"243.22.123.44",
"healthy_check_mode":"tcp",
"lb_mode":"loop",
"port":80,
"idc_id":"hd"
}
'

# 给服务器组添加服务器
curl -X POST ${admin}/groups/g-hd/servers -H 'Content-type:application/json' -d '
{
  "server_ip":"10.1.24.56",
  "weight":100
}
'


# 创建地址池
curl -X POST ${admin}/pools -H 'Content-type:application/json' -d '
{
"name":"pool-all",
"lb_mode":0
}
'
# 给地址池添加group
curl -X POST ${admin}/pools/pool-all/groups -H 'Content-type:application/json' -d '
{
"group_id":"g-hd",
"weight":100
}
'
# 创建services
curl -X POST ${admin}/services -H 'Content-type:application/json' -d '
{
   "name": "srv-test",
   "pool_id": "pool-all",
   "ttl": 6000
}
'

#创建路由
curl -X POST ${admin}/services/srv-test/routes -H 'Content-type:application/json' -d '
{
  
   "domains": "somestupidname.org",
   "name": "stupid_route",
   "priority": 0,
   "rules": ""
}
'
#创建泛解析路由，* 号只能出现在域名开头
curl -X POST ${admin}/services/srv-test/routes -H 'Content-type:application/json' -d '
{
  
   "domains": "*.stupidname.org",
   "name": "stupid_route2",
   "priority": 0,
   "rules": ""
}
'



```


### 2 请求域名解析接口
#### 接口使用方式：

GET 127.0.0.1:4567/resolve?hosts=somestupidname.org

```shell
curl '127.0.0.1:4567/resolve?hosts=somestupidname.org'
```
**请求参数说明**

- hosts: 解析的域名，多个用英文 ',' 分隔。

除了hosts 为必要参数，还可以添加其他的自定义参数。

响应示例：

```
HTTP/1.1 200 OK
Server: ifly-nameserver
Date: Tue, 22 Jun 2021 06:28:59 GMT
Content-Type: application/json
Content-Length: 143

{"dns":[{"host":"somestupidname.org","ips":[{"ip":"243.22.123.44","port":80}],"ttl":6000}],"client_ip":"127.0.0.1"}

```

