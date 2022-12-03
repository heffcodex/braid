package braid

import (
	"strconv"
)

type ErrorCodeInt int

const (
	ECINone ErrorCodeInt = iota
	eciReserved1
	eciReserved2
	eciReserved3
	eciReserved4
	eciReserved5
	ECIInternal
	ECICSRFTokenMismatch
	eciReserved8
	eciReserved9
	ECIValidationFail
	eciReserved11
	eciReserved12
	ECIInvalidQuery
	ECIInvalidPayload

	// ECICustom must be used by default for custom errors made with NewErrorCode
	ECICustom = -1
)

var (
	ErrorCodeNone              = ErrorCode{code: ECINone, message: ""}
	ErrorCodeInternal          = ErrorCode{code: ECIInternal, message: "Internal error"}
	ErrorCodeCSRFTokenMismatch = ErrorCode{code: ECICSRFTokenMismatch, message: "CSRF token mismatch"}
	ErrorCodeValidationFail    = ErrorCode{code: ECIValidationFail, message: "Validation fail"}
	ErrorCodeInvalidQuery      = ErrorCode{code: ECIInvalidQuery, message: "Invalid query"}
	ErrorCodeInvalidPayload    = ErrorCode{code: ECIInvalidPayload, message: "Invalid payload"}
)

type ErrorCode struct {
	code    ErrorCodeInt
	message string
}

func NewErrorCode(code ErrorCodeInt, message string) ErrorCode {
	return ErrorCode{code: code, message: message}
}

func (e *ErrorCode) Code() ErrorCodeInt {
	return e.code
}

func (e *ErrorCode) GetMessage() string {
	return e.message
}

func (e *ErrorCode) SetMessage(message string) *ErrorCode {
	e.message = message
	return e
}

func (e *ErrorCode) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(e.code))), nil
}
