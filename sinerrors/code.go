package sinerrors

import ()

func SetErrorAI(code Code, msg string) Code {
	return AddError(code.Code(), msg)
}
func SetErrorWithToastAI(code Code, msg string, toast string) Code {
	return AddErrorWithToast(code.Code(), msg, toast)
}
