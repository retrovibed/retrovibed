package cmdmeta

import (
	"cmp"
	"fmt"
	"io"
	"os"
	"slices"
	"text/tabwriter"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/me"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
)

type CloudMeLs struct{}

func (t CloudMeLs) Run(gctx *cmdopts.Global, id *cmdopts.SSHID) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return err
	}

	c, session, err := cloudClient(gctx.Context, gctx, signer)
	if err != nil {
		return err
	}

	profiles, err := me.Profiles(gctx.Context, c)
	if err != nil {
		return errorsx.Wrap(err, "unable to list profiles")
	}

	return t.render(os.Stdout, session.Profile.GetId(), profiles)
}

// render the profiles, the current profile first and marked with a *.
func (t CloudMeLs) render(w io.Writer, current string, profiles []*authn.Profile) error {
	slices.SortStableFunc(profiles, func(a, b *authn.Profile) int {
		switch {
		case a.GetId() == current && b.GetId() != current:
			return -1
		case b.GetId() == current && a.GetId() != current:
			return 1
		}

		return cmp.Compare(a.GetId(), b.GetId())
	})

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "\tID\tACCOUNT\tNAME\tEMAIL\tSTATUS")
	for _, p := range profiles {
		mark := ""
		if p.GetId() == current {
			mark = "*"
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n", mark, p.GetId(), p.GetAccountId(), p.GetDisplay(), p.GetEmail(), authn.ProfileState(p))
	}

	return tw.Flush()
}
