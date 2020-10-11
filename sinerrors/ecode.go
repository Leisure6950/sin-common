package sinerrors

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type errInfo struct {
	code  int32
	msg   string
	toast string
}

var (
	toasts   sync.Map               // map[int]string
	errInfos sync.Map               // map[int32]errInfo
	codes    = map[int32]struct{}{} // register codes.
	mux      sync.Mutex
)

// NOTE: ecode must unique in global, the New will check repeat and then panic.

// Error returns a  ecode.Codes and register associated ecode message
// NOTE: Error codes and messages should be kept together.
// ecode must unique in global, the Error will check repeat and then panic.
func AddError(e int32, msg string) Code {
	errInfos.Store(e, errInfo{
		code:  e,
		msg:   msg,
		toast: "",
	})
	return eInt(e)
}
func AddErrorWithToast(e int32, msg string, toast string) Code {
	errInfos.Store(e, errInfo{
		code:  e,
		msg:   msg,
		toast: toast,
	})
	return eInt(e)
}
func add(e int32) Code {
	if _, ok := codes[e]; ok {
		fmt.Printf("ecode: %d already exist \n", e)
	}
	codes[e] = struct{}{}
	return eInt(e)
}

// Codes ecode error interface which has a code & message.
type eCodes interface {
	// Error return Code in string form
	Error() string
	// Code get error code.
	Code() int32
	// Message get code message.
	Message() string
	// Toast get code toast
	Toast() string
	// Equal for compatible.
	Equal(error) bool
}

// A Code is an int error code spec.
type Code int32

func (e Code) Error() string {
	return strconv.FormatInt(int64(e), 10)
}

// Code return error code
func (e Code) Code() int32 { return int32(e) }

// Message return error message
func (e Code) Message() string {
	v, ok := errInfos.Load(e.Code())
	if !ok {
		return e.Error()
	}
	return v.(errInfo).msg
}

func (e Code) Toast() string {
	v, ok := errInfos.Load(e.Code())
	if !ok {
		return e.Error()
	}
	return v.(errInfo).toast
}

// Equal for compatible.
func (e Code) Equal(err error) bool { return EqualError(e, err) }

// 错误码拼接详细错误信息
func (c Code) DetailF(f string, args ...interface{}) error {
	return &errorDetailInfo{
		code:    c,
		message: c.Message() + "(" + fmt.Sprintf(f, args...) + ")",
		err:     nil,
	}
}

// 错误码吗拼接详细错误信息
func (c Code) DetailW(params ...string) error {
	return &errorDetailInfo{
		code:    c,
		message: c.Message() + "(" + strings.Join(params, ",") + ")",
		err:     nil,
	}
}

// eInt parse code int to error.
func eInt(i int32) Code { return Code(i) }

// eString parse code string to error.
func eString(e string) Code {
	if e == "" {
		return eInt(0)
	}
	// try error string
	i, err := strconv.Atoi(e)
	if err != nil {
		return eInt(500)
	}
	return Code(i)
}

// Cause cause from error to ecode.
func Cause(e error) eCodes {
	if e == nil {
		return eInt(0)
	}
	ec, ok := errors.Cause(e).(eCodes)
	if ok {
		return ec
	}
	return eString(e.Error())
}

// Equal equal a and b by code int.
func Equal(a, b eCodes) bool {
	if a == nil {
		a = eInt(0)
	}
	if b == nil {
		b = eInt(0)
	}
	return a.Code() == b.Code()
}

// EqualError equal error
func EqualError(code eCodes, err error) bool {
	return Cause(err).Code() == code.Code()
}
