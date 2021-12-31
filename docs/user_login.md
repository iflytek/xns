
开启登陆接口权限校验.
配置文件中新增： 

```
login {
    enabled true
    session_timeout 3600 # session 过期时长。session 目前存放在内存。
}
```

### 登陆接口

#### create user
````
POST {{url}}/users/create
Content-Type: application/json
X-HYDRANT-TOKEN: 0d519fa0-c32b-4d7e-9194-fdbf6ecd4e09  

{
  "username": "sjliu9",
  "password": "123456",
  "mode": "admin"  
}

````
 参数解释
- username: 用户名
- password: 密码
- mode: 用户模式： admin or  readonly 分别对应超级权限和只读权限

只有在系统没有任何用户情况下，可以直接调用接口创建 admin 用户。后面只有admin 用户才能创建其他用户

#### login
使用用户名和密码登陆
```
POST {{url}}/users/login
Content-Type: application/json
X-HYDRANT-TOKEN:  {{token}}1

{
  "username": "sjliu8",
  "password": "123456"
}
```

响应
````json
{"token":"4da3639d-ff48-423e-af7b-8bc755398801","header":"X-HYDRANT-TOKEN"}
````
响应参数解释

- token： 登陆token
- header: 请求头，后续客户端带上该请求头，值为token。 请求就可以视为已登陆


#### 获取已经登陆用户的信息

````http request

GET {{url}}/user_info
X-HYDRANT-TOKEN: c2a81695-ccc6-4b9f-b015-de0f1ae0402b

````

