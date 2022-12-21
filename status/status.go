package status

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Code int

var (
	_ json.Marshaler = (*Status)(nil)
	_ fmt.Stringer   = (*Status)(nil)
)

type Status struct {
	parent   *Status
	code     Code
	codeHTTP int
	message  string
}

// New creates a new Status with the given Code and optional message.
func New(code Code, message ...string) *Status {
	_message := ""
	if len(message) > 1 {
		panic("too many arguments")
	} else if len(message) == 1 {
		_message = message[0]
	}

	return &Status{code: code, message: _message}
}

func (e *Status) Is(s *Status) bool {
	if e == nil || s == nil {
		return false
	} else if e.code == s.code {
		return true
	}

	for p := e.parent; p != nil; p = p.parent {
		if p.code == s.code {
			return true
		}
	}

	return false
}

func (e *Status) Parent() *Status {
	return e.parent
}

func (e *Status) AttachTo(parent *Status) *Status {
	e.parent = parent
	return e
}

func (e *Status) Code() Code {
	return e.code
}

func (e *Status) GetHTTP() int {
	if e.codeHTTP != 0 {
		return e.codeHTTP
	}

	for p := e.parent; p != nil; p = p.parent {
		if p.codeHTTP != 0 {
			return p.codeHTTP
		}
	}

	panic("no status code found")
}

func (e *Status) SetHTTP(code int) *Status {
	e.codeHTTP = code
	return e
}

func (e *Status) Message() string {
	return e.message
}

func (e *Status) Copy() *Status {
	return &Status{
		parent:   e.parent.Copy(),
		code:     e.code,
		codeHTTP: e.codeHTTP,
		message:  e.message,
	}
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
