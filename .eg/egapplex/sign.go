package egapplex

import (
	"context"
	"fmt"
	"strings"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

// SignOption returns extra codesign CLI flag text to splice into the invocation built by Sign.
type SignOption func() string

// SignDeep adds --deep, recursively signing nested code.
func SignDeep() SignOption {
	return func() string { return "--deep" }
}

// SignRuntime adds --options runtime, enabling the hardened runtime (required for notarization).
func SignRuntime() SignOption {
	return func() string { return "--options runtime" }
}

// SignEntitlements adds --entitlements path.
func SignEntitlements(path string) SignOption {
	return func() string { return fmt.Sprintf("--entitlements %s", path) }
}

// Sign unlocks keychain then codesigns target with identity, applying any extra flags from opts.
func Sign(identity, keychain, target string, opts ...SignOption) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		flags := make([]string, 0, len(opts))
		for _, opt := range opts {
			flags = append(flags, opt())
		}

		cmd := shell.Newf(
			"codesign --force %s --sign %q --keychain %s %s",
			strings.Join(flags, " "), identity, keychain, target,
		)

		if err := UnlockKeychain(keychain)(ctx, o); err != nil {
			return err
		}

		return shell.Run(ctx, cmd)
	}
}
