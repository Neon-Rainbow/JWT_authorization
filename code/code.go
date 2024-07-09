package code

type ResponseCode int

const (
	Success ResponseCode = iota + 1000

	RequestTimeout
	ServerBusy
	LoginParamsError
	LoginGetUserInformationError
	LoginPasswordError
	LoginGenerateTokenError
	LoginUserIsFrozen
	LoginUserNotFound

	RegisterParamsError
	RegisterCheckUserExistsError
	RegisterUsernameExists
	RegisterTelephoneExists
	RegisterCreateUserError

	RequestUnauthorized

	RefreshTokenError
)

var codeMsgMap = map[ResponseCode]string{
	Success:                      "success",
	RequestTimeout:               "request timeout",
	ServerBusy:                   "server busy",
	LoginParamsError:             "login params error",
	LoginGetUserInformationError: "Get User Information error",
	LoginPasswordError:           "login password error",
	LoginGenerateTokenError:      "login generate token error",
	LoginUserIsFrozen:            "login user is Frozen",
	LoginUserNotFound:            "login user not found",
	RegisterParamsError:          "Register params error",
	RegisterCheckUserExistsError: "Register check user exists error",
	RegisterUsernameExists:       "Register username exists",
	RegisterTelephoneExists:      "Register telephone exists",
	RegisterCreateUserError:      "Register create user error",

	RequestUnauthorized: "request unauthorized",

	RefreshTokenError: "refresh token error",
}

func (code ResponseCode) Message() string {
	return codeMsgMap[code]
}
