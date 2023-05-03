package status

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Code int

func (c Code) HTTP() int {
	if c < 1000 {
		return int(c)
	}

	return int(c / 1000)
}

var (
	_ json.Marshaler = (*Status)(nil)
	_ fmt.Stringer   = (*Status)(nil)
)

type Status struct {
	code    Code
	message string
}

func FromFiber(e *fiber.Error) *Status {
	return New(Code(e.Code), e.Message)
}

func New(code Code, message ...string) *Status {
	_message := strings.Join(message, ": ")
	return &Status{code: code, message: _message}
}

func (e *Status) Code() Code {
	return e.code
}

func (e *Status) Message() string {
	return e.message
}

func (e Status) String() string {
	str := "[" + strconv.Itoa(int(e.code)) + "]"
	if e.message != "" {
		str += " " + e.message
	}

	return str
}

func (e Status) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(e.code))), nil
}

func (e *Status) Describe(format string, args ...any) *Status {
	clone := e.clone()
	message := fmt.Sprintf(format, args...)

	if clone.message == "" {
		clone.message = message
	} else {
		clone.message = clone.message + ": " + message
	}

	return clone
}

func (e *Status) clone() *Status {
	return &Status{
		code:    e.code,
		message: e.message,
	}
}
