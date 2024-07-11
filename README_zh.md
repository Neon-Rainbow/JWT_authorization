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

[Postman json](./docs/JWT.postman_collection.json)

## 状态码

| Error Code                | Description                        |
|---------------------------|------------------------------------|
| 1000                      | Success                            |
| 1001                      | RequestTimeout                     |
| 1002                      | ServerBusy                         |
| 1003                      | LoginParamsError                   |
| 1004                      | LoginGetUserInformationError       |
| 1005                      | LoginPasswordError                 |
| 1006                      | LoginGenerateTokenError            |
| 1007                      | LoginUserIsFrozen                  |
| 1008                      | LoginUserNotFound                  |
| 1009                      | RegisterParamsError                |
| 1010                      | RegisterCheckUserExistsError       |
| 1011                      | RegisterUsernameExists             |
| 1012                      | RegisterTelephoneExists            |
| 1013                      | RegisterCreateUserError            |
| 1014                      | RequestUnauthorized                |
| 1015                      | RefreshTokenError                  |
| 1016                      | FrozenUserIDRequired               |
| 1017                      | FrozenUserError                    |
| 1018                      | ThawUserIDRequired                 |
| 1019                      | ThawUserError                      |
| 1020                      | DeleteUserError                    |