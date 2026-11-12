package cmdmeta

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	shallowsenv "github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
)

type CloudAccountEdit struct {
	Display string `flag:"" name:"display" help:"set the account's display name" required:"true"`
}

func (t CloudAccountEdit) Run(gctx *cmdopts.Global) (err error) {
	ctx, done := context.WithTimeout(gctx.Context, 30*time.Second)
	defer done()

	id, err := sshx.Load(shallowsenv.PrivateKeyPath())
	if err != nil {
		return err
	}

	session, err := authn.AutoRegistrationOnce(ctx, id, gctx.AccountID, false, cmdopts.ChooseProfile)
	if err != nil {
		return errorsx.Wrap(err, "unable to retrieve cloud session")
	}

	c, err := authn.AutoJWTClient(ctx, id, session.Account.GetId())
	if err != nil {
		return errorsx.Wrap(err, "unable to authenticate")
	}

	return t.run(ctx, os.Stdout, c, fmt.Sprintf("https://%s", env.Deeppool()), session.Account)
}

// run updates the display name of the account. the service replaces the contact details with
// whatever is sent, so the current account is carried forward for everything but the display name.
func (t CloudAccountEdit) run(ctx context.Context, w io.Writer, c *http.Client, endpoint string, current *authn.Account) error {
	zero := time.Time{}.Format(time.RFC3339Nano)

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, fmt.Sprintf("%s/a/", endpoint), nil)
	if err != nil {
		return errorsx.Wrap(err, "unable to create request")
	}

	if err = httpx.EncodeJSON(req, &authn.AccountUpdateRequest{
		Account: &authn.Account{
			Id:          current.GetId(),
			Description: t.Display,
			Email:       current.GetEmail(),
			Phone:       current.GetPhone(),
			// the service parses every timestamp before applying the update, even though only
			// the display name, email and phone are written.
			CreatedAt:  stringsx.First(current.GetCreatedAt(), zero),
			UpdatedAt:  stringsx.First(current.GetUpdatedAt(), zero),
			DisabledAt: stringsx.First(current.GetDisabledAt(), zero),
		},
	}); err != nil {
		return errorsx.Wrap(err, "unable to encode request")
	}

	resp, err := httpx.AsError(c.Do(req))
	if err != nil {
		return errorsx.Wrap(err, "unable to update account")
	}
	defer resp.Body.Close()

	var result authn.AccountLookupResponse
	if err = jsonx.UnmarshalRead(resp.Body, &result); err != nil {
		return errorsx.Wrap(err, "unable to decode response")
	}

	_, err = fmt.Fprintf(w, "account     %s\ndisplay     %s\n", result.Account.GetId(), result.Account.GetDescription())
	return err
}
