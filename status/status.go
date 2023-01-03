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

func (e *Status) Code() Code {
	return e.code
}

func (e *Status) HTTP() int {
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

func (e *Status) AttachTo(parent *Status) *Status {
	clone := e.clone()
	clone.parent = parent

	return clone
}

func (e *Status) AppendMessage(message string) *Status {
	clone := e.clone()
	clone.message = clone.message + ": " + message

	return clone
}

func (e *Status) WithHTTP(code int) *Status {
	clone := e.clone()
	clone.codeHTTP = code

	return clone
}

func (e *Status) clone() *Status {
	parent := e.parent
	if parent != nil {
		parent = parent.clone()
	}

	return &Status{
		parent:   parent,
		code:     e.code,
		codeHTTP: e.codeHTTP,
		message:  e.message,
	}
}
