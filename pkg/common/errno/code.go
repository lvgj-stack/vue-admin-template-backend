package errno

var (
	// OK represents a successful request.
	OK                   = &Errno{Code: 0, Message: "OK"}
	ErrBind              = &Errno{Code: 10002, Message: "Error occurred when binding request."}
	ErrUserNotFound      = &Errno{Code: 10003, Message: "User not found."}
	ErrPasswordIncorrect = &Errno{Code: 10004, Message: "Password incorrect."}
	ErrToken             = &Errno{Code: 10005, Message: "Error Token."}
	ErrValidation        = &Errno{Code: 10006, Message: "Error validate user."}
	ErrEncrypt           = &Errno{Code: 10007, Message: "Error Encrypt the password."}
	ErrUserAlreadyExist  = &Errno{Code: 10008, Message: "Username already exist."}

	// InternalServerError represents all unknown server-side errors.
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error"}
)
