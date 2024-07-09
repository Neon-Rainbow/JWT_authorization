# JWT Authorization

[简体中文](./README_zh.md)

A demo project for user authentication and authorization control using JWT, MySQL, and Redis technology stack.

## Project Structure

```
.
├── README.md
├── README_zh.md
├── code
│   └── code.go
├── config
│   └── config.go
├── config.json
├── go.mod
├── go.sum
├── internal
│   ├── controller
│   │   ├── contextKey.go
│   │   ├── deleteUser.go
│   │   ├── frozen.go
│   │   ├── getAdminInfo.go
│   │   ├── getUserID.go
│   │   ├── getUserInfo.go
│   │   ├── login.go
│   │   ├── refreshToken.go
│   │   ├── register.go
│   │   └── response.go
│   ├── dao
│   │   ├── token.go
│   │   └── user.go
│   └── service
│       ├── EncryptPassword.go
│       ├── deleteUser.go
│       ├── frozen.go
│       ├── login.go
│       ├── refreshToken.go
│       └── register.go
├── main.go
├── middleware
│   ├── adminMiddleware.go
│   └── jwtMiddleware.go
├── model
│   ├── apiError.go
│   ├── token.go
│   └── user.go
├── route
│   └── route.go
└── util
    ├── MySQL
    │   └── MySQL.go
    ├── Redis
    │   └── Redis.go
    ├── initSQL.go
    └── jwt
        └── jwt.go
```

## Project Configuration

You need to create a `config.json` file and place it in the root directory of the project.
The format of `config.json` is:

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

## Project Startup

Run the following command in the terminal:

```shell
go run main.go
```

## Project APIs:

### POST /api/auth/login

Login endpoint

Request parameters:

```json
{
  "username": "admin",
  "password": "admin"
}
```

### POST /api/auth/register

Registration endpoint

Request parameters:

```json
{
  "username": "admin",
  "password": "admin"
}
```

### POST /api/auth/refresh

Refresh token endpoint

Request parameters:
+ Query Params: `?refresh_token={refresh_token}`

### POST /api/user/frozen

Freeze user endpoint

Request parameters:
+ Header: `Authorization: Bearer {access_token}`

### POST /api/user/delete_account

Delete user endpoint

Request parameters:
+ Header: `Authorization: Bearer {access_token}`

### POST /api/admin/frozen

Freeze user endpoint (admin)

Request parameters:
+ Query Params: `?user_id={user_id}`
+ Header: `Authorization: Bearer {access_token}`

### POST /api/admin/thaw

Unfreeze user endpoint (admin)

Request parameters:
+ Query Params: `?user_id={user_id}`
+ Header: `Authorization: Bearer {access_token}`

### POST /api/admin/delete_account

Delete user endpoint (admin)

Request parameters:
+ Query Params: `?user_id={user_id}`
+ Header: `Authorization: Bearer {access_token}`

## Status Code

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