package egapplex

import (
	"context"
	"fmt"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

// Upload submits artifactPath to App Store Connect via altool. platformType is passed directly
// as altool's --type value (e.g. "macos", "ios"). Fails the step if altool's own log reports
// "UPLOAD FAILED", since altool can exit 0 on a rejected upload.
func Upload(apikey, issuerid, artifactPath, platformType string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		cmd := shell.Newf(
			"xcrun altool --upload-app --type %s -f %s --apiKey ${APPLE_API_KEY} --apiIssuer ${APPLE_ISSUER_ID} 2>&1 | tee /tmp/altool.log && ! grep -q 'UPLOAD FAILED' /tmp/altool.log",
			platformType, artifactPath,
		).
			Environ("APPLE_API_KEY", apikey).
			Environ("APPLE_ISSUER_ID", issuerid)

		return shell.Run(ctx, cmd)
	}
}

// Notarize submits artifactPath to Apple's notary service and staples the resulting ticket.
// darwin-only, no iOS equivalent. Builds the AuthKey path itself from apikey rather than relying
// on a hardcoded shell string matching AuthKey's internal naming.
func Notarize(apikey, issuerid, artifactPath string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		keypath := fmt.Sprintf("~/.private_keys/AuthKey_%s.p8", apikey)

		return shell.Run(
			ctx,
			shell.Newf("xcrun notarytool submit %s --key %s --key-id %s --issuer %s --wait", artifactPath, keypath, apikey, issuerid),
			shell.Newf("xcrun stapler staple %s", artifactPath),
		)
	}
}
