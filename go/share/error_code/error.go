package error_code

import "fmt"

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	ServerErrorCode = 10000 + iota
	DatabaseErrorCode
)
const (
	UserExistsCode = 20000 + iota
	PasswordIsEasyCode
	UserNotExistsCode
	UserNumberExistsCode
	UserEmailExistsCode
	PasswordFailCode
	OldPasswordFailCode
	PassWordSameCode
	CheckPhoneFailCode
	NotLoginCode
	InvalidTokenCode
	ShouldBindErrorCode
	TokenOutErrorCode
)
const (
	DeviceNotFoundCode = 30000 + iota
	DeviceIsBindCode
	DeviceNotBindCode
	NotDeviceOwnerCode
)

var (
	ServerError      = &APIError{Code: ServerErrorCode, Message: "服务器内部错误。"}
	DatabaseError    = &APIError{Code: DatabaseErrorCode, Message: "数据库操作失败。"}
	UserExists       = &APIError{Code: UserExistsCode, Message: "用户已注册。"}
	UserNotExists    = &APIError{Code: UserNotExistsCode, Message: "用户不存在。"}
	PasswordIsEasy   = &APIError{Code: UserNotExistsCode, Message: "密码过于简单。"}
	UserNumberExists = &APIError{Code: UserNumberExistsCode, Message: "手机号已注册。"}
	UserEmailExists  = &APIError{Code: UserEmailExistsCode, Message: "邮箱已注册。"}
	PasswordFail     = &APIError{Code: PasswordFailCode, Message: "账号或密码错误。"}
	OldPasswordFail  = &APIError{Code: OldPasswordFailCode, Message: "旧密码错误。"}
	PassWordSame     = &APIError{Code: PassWordSameCode, Message: "新密码与旧密码相同。"}
	DeviceNotFound   = &APIError{Code: DeviceNotFoundCode, Message: "设备未注册。"}
	DeviceNotBind    = &APIError{Code: DeviceNotBindCode, Message: "设备未被绑定。"}
	DeviceIsBind     = &APIError{Code: DeviceIsBindCode, Message: "设备已被其他用户绑定。"}
	NotDeviceOwner   = &APIError{Code: NotDeviceOwnerCode, Message: "您不是设备拥有者。"}
	CheckPhoneFail   = &APIError{Code: CheckPhoneFailCode, Message: "手机号验证失败。"}
	NotLogin         = &APIError{Code: NotLoginCode, Message: "用户未登录。"}
	InvalidToken     = &APIError{Code: InvalidTokenCode, Message: "无效的Token格式。"}
	TokenOutError    = &APIError{Code: NotDeviceOwnerCode, Message: "Token已被注销。"}
	ShouldBindError  = &APIError{Code: ShouldBindErrorCode, Message: "绑定参数错误。"}
)

func NewAPIError(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *APIError) SetError() *APIError {
	return &APIError{
		Code:    e.Code,
		Message: e.Message,
	}

}
