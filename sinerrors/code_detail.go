package sinerrors

type errorDetailInfo struct {
	code    Code
	message string
	toast   string
	err     error
}

// Error return Code in string form
func (e *errorDetailInfo) Error() string {
	if e.err == nil {
		//return strconv.FormatInt(int64(e.code), 10)+":"+e.toast
		return e.code.Error()
	}
	return e.err.Error()
}

// Code get error code.
func (e *errorDetailInfo) Code() int32 {
	return e.code.Code()
}

// Message get code message.
func (e *errorDetailInfo) Message() string {
	return e.message
}

func (e *errorDetailInfo) Toast() string {
	return e.code.Toast()
}

// Equal for compatible.
func (e *errorDetailInfo) Equal(err error) bool {
	if e2, ok := err.(*errorDetailInfo); ok {
		return *e2 == *e
	}
	return false
}

func NewTmpError(errCode int32, errMsg, errToast string) *errorDetailInfo {
	return &errorDetailInfo{
		code:    eInt(errCode),
		message: errMsg,
		toast:   errToast,
	}
}
