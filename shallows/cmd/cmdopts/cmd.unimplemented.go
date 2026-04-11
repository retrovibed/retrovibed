package cmdopts

import (
	"errors"
)

// utility command for placeholders
type CmdNotImplemented struct{}

func (t CmdNotImplemented) Run() (err error) {
	return errors.ErrUnsupported
}
