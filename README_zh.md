# JWT_authorization

[English version](./README.md)

一个使用了JWT + MySQL + Redis技术栈的用户认证,权限控制的Demo项目

## 项目结构

```
.
├── README.md
├── README_zh.md
├── code
│   └── code.go
├── config
│   └── config.go
├── config.json
├── go.mod
├── go.sum
├── internal
│   ├── controller
│   │   ├── contextKey.go
│   │   ├── deleteUser.go
│   │   ├── frozen.go
│   │   ├── getAdminInfo.go
│   │   ├── getUserID.go
│   │   ├── getUserInfo.go
│   │   ├── login.go
│   │   ├── refreshToken.go
│   │   ├── register.go
│   │   └── response.go
│   ├── dao
│   │   ├── token.go
│   │   └── user.go
│   └── service
│       ├── EncryptPassword.go
│       ├── deleteUser.go
│       ├── frozen.go
│       ├── login.go
│       ├── refreshToken.go
│       └── register.go
├── main.go
├── middleware
│   ├── adminMiddleware.go
│   └── jwtMiddleware.go
├── model
│   ├── apiError.go
│   ├── token.go
│   └── user.go
├── route
│   └── route.go
└── util
    ├── MySQL
    │   └── MySQL.go
    ├── Redis
    │   └── Redis.go
    ├── initSQL.go
    └── jwt
        └── jwt.go
```

## 项目配置
需要编写config.json并且放在项目根目录下
config.json格式:
```json
{
  "address":"127.0.0.1",
  "port": 8080,
  "MySQL": {
    "host": "127.0.0.1",
    "port": 3306,
    "username": "admin",
    "password": "admin",
    "database": "database"
  },
  "Redis": {
    "host": "127.0.0.1",
    "port": 6379,
    "username": "admin",
    "password": "admin",
    "database": 0
  },
  "JWT": {
    "secret": "secret"
  },
  "passwordSecret": "passwordSecret"
}
```

## 项目启动
终端执行下列指令:
```shell
go run main.go
```
## 项目接口:

### POST /api/auth/login
注册接口
请求参数:
```json
{
  "username": "admin",
  "password": "admin"
}
```


### POST /api/auth/register
登录接口
请求参数:
```json
{
  "username": "admin",
  "password": "admin"
}
```

### POST /api/auth/register
注册接口
请求参数:
```json
{
  "username": "admin",
  "password": "admin"
}
```

### POST /api/auth/refresh
刷新token接口

请求参数:
+ Query Params: ?refresh_token={refresh_token}

### POST /api/user/frozen
冻结用户接口

请求参数:
+ Header: Authorization: Bearer {access_token}

### POST /api/user/delete_account
删除用户接口

请求参数:
+ Header: Authorization: Bearer {access_token}

### POST /api/admin/frozen
冻结用户接口

请求参数:
+ Query Params: ?user_id={user_id}
+ Header: Authorization: Bearer {access_token}

### POST /api/admin/thaw
解冻用户接口

请求参数:
+ Query Params: ?user_id={user_id}
+ Header: Authorization: Bearer {access_token}

### POST /api/admin/delete_account
删除用户接口

请求参数:
+ Query Params: ?user_id={user_id}
+ Header: Authorization: Bearer {access_token}