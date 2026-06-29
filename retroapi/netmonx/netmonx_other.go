//go:build !(linux || darwin || android)

package netmonx

import (
	"context"
	"errors"
)

func startWatcher(_ context.Context, _ chan<- struct{}) error {
	return errors.ErrUnsupported
}

func platformMetered(_ string) bool { return false }

func defaultRouteInterface() string {
	return ""
}
