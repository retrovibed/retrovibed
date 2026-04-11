package errors

import "fmt"

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Error uint

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	ErrSuccess Error = iota
	ErrBadParameter
	ErrDuplicateEntry
	ErrUnexpectedResponse
	ErrNotFound
	ErrNotModified
	ErrNotImplemented
	ErrNotAuthorized
	ErrExpired
	ErrInternalAppError
	ErrChannelBlocked
	ErrOutOfOrder
)

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e Error) Error() string {
	switch e {
	case ErrSuccess:
		return "ErrSuccess"
	case ErrBadParameter:
		return "ErrBadParameter"
	case ErrDuplicateEntry:
		return "ErrDuplicateEntry"
	case ErrUnexpectedResponse:
		return "ErrUnexpectedResponse"
	case ErrNotFound:
		return "ErrNotFound"
	case ErrNotModified:
		return "ErrNotModified"
	case ErrNotImplemented:
		return "ErrNotImplemented"
	case ErrNotAuthorized:
		return "ErrNotAuthorized"
	case ErrExpired:
		return "ErrExpired"
	case ErrInternalAppError:
		return "ErrInternalAppError"
	case ErrChannelBlocked:
		return "ErrChannelBlocked"
	case ErrOutOfOrder:
		return "ErrOutOfOrder"
	default:
		return "[?? Invalid Error value]"
	}
}

func (e Error) With(args ...interface{}) error {
	return fmt.Errorf("%w: %s", e, fmt.Sprint(args...))
}

func (e Error) Withf(format string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", e, fmt.Sprintf(format, args...))
}
