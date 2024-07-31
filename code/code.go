package code

type ResponseCode int

const (
	Success ResponseCode = iota + 1000

	RequestTimeout
	RequestCanceled
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

	FrozenUserIDRequired
	FrozenUserError

	ThawUserIDRequired
	ThawUserError

	DeleteUserTokenError

	PermissionGetError
	PermissionDenied
	PermissionParamsError
	PermissionChangeError
)

var codeMsgMap = map[ResponseCode]string{
	Success:         "success",
	RequestTimeout:  "request timeout",
	RequestCanceled: "request canceled",
	ServerBusy:      "server busy",

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

	FrozenUserIDRequired: "frozen user id required",
	FrozenUserError:      "frozen user error",

	ThawUserIDRequired: "thaw user id required",
	ThawUserError:      "thaw user error",

	DeleteUserTokenError: "delete user error",

	PermissionGetError:    "permission get error",
	PermissionDenied:      "permission denied",
	PermissionParamsError: "permission params error",
	PermissionChangeError: "permission add error",
}

func (code ResponseCode) Message() string {
	return codeMsgMap[code]
}

func (code ResponseCode) Error() string {
	return codeMsgMap[code]
}

func (code ResponseCode) Int() int {
	return int(code)
}

const (
	Permission1 = 1 << iota
	Permission2
	Permission3
	Permission4
	Permission5
	Permission6
	Permission7
	AllPermission = Permission1 | Permission2 | Permission3 | Permission4 | Permission5 | Permission6 | Permission7
)
