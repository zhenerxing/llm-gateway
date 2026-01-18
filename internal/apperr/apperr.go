package apperr

import(
	"errors"
	"fmt"
)

type Type string

const (
	TypePlatform Type = "platform"
	TypeUpstream Type = "upstream"
)

type Error struct {
	//对外语义
	Code       string
	Message    string
	Type       Type
	RetryAfter int // seconds; 0 means unset
	Details    map[string]any

	//对内技术细节
	cause error
}

// 实现Error() string后，Error就是error接口
func (e *Error) Error() string{
	if e == nil{
		return "<nil>"
	}
	if e.cause != nil{
		return fmt.Sprintf("%s:%s: %v", e.Type, e.Code, e.cause)
	}
	return fmt.Sprintf("%s:%s", e.Type, e.Code)
}

// 接入go的错误链
func (e *Error) Unwrap() error { return e.cause }

func New(code string , typ Type , message string) *Error{
	return &Error{Code: code, Type: typ, Message: message}
}

// 封装错误
func Wrap(code string, typ Type, message string, cause error) *Error {
	return &Error{Code: code, Type: typ, Message: message, cause: cause}
}

// 封装As
func As(err error) (*Error,bool){
	var ae *Error
	if errors.As(err,&ae) && ae != nil{
		return ae , true
	}
	return nil,false
} 

func (e *Error) WithDetails(details map[string]any) *Error {
	if e == nil {
		return nil
	}
	e.Details = details
	return e
}