package eppgo

import "fmt"

const (
	// ErrorCodeRegisteringWrongType is used when you are trying to register a
	// child element with the wrong type. It should be a struct or a pointer to a
	// struct.
	ErrorCodeRegisteringWrongType ErrorCode = iota

	// ErrorCodeUnknownElement is used when an unknown element is found in the
	// XML.
	ErrorCodeUnknownElement
)

// ErrorCode stores the type of the error.
type ErrorCode int

// String translates the error code into a human readable text.
func (e ErrorCode) String() string {
	switch e {
	case ErrorCodeRegisteringWrongType:
		return "registering wrong type"
	case ErrorCodeUnknownElement:
		return "unknown element found in the XML"
	}

	return ""
}

// Error is the type used when something goes wrong in the library. It contains
// the necessary information to find where the problem is.
type Error struct {
	Code      ErrorCode
	Reference string
}

// Error translates the error into a human readable text.
func (e Error) Error() string {
	if e.Reference != "" {
		return fmt.Sprintf("[eppgo] %s : %s", e.Code, e.Reference)
	}

	return fmt.Sprintf("[eppgo] %s", e.Code)
}
