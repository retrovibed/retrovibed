package ddiscapi

import (
	"context"
	"fmt"

	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

func AuthzPermPeer(ctx context.Context, cause error) (_ context.Context, token *metaapi.Token, err error) {
	if cause != nil {
		return ctx, nil, errorsx.Authorization(fmt.Errorf("not authorized"))
	}

	if token, err = metaapi.FromContext(ctx); err != nil {
		return ctx, token, errorsx.Wrap(errorsx.Authorization(err), "not authorized")
	} else if !token.Usermanagement {
		return ctx, token, errorsx.Authorization(errorsx.WithStack(fmt.Errorf("not authorized: permission denied")))
	}

	return ctx, token, nil
}
