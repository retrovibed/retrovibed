package cmdopts

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/mattn/go-isatty"

	"github.com/retrovibed/retrovibed/retroapi/authn"
)

// ChooseProfile interactively selects the profile to use when the ssh key is associated with multiple.
// requires a terminal, otherwise no profile is selected and the caller reports the ambiguity.
func ChooseProfile(ctx context.Context, profiles []*authn.Authn) (*authn.Authn, error) {
	if !isatty.IsTerminal(os.Stdin.Fd()) || !isatty.IsTerminal(os.Stdout.Fd()) {
		return nil, nil
	}

	options := make([]huh.Option[*authn.Authn], 0, len(profiles))
	for _, p := range profiles {
		options = append(options, huh.NewOption(fmt.Sprintf("%s (%s)", p.Account.GetDescription(), p.Account.GetId()), p))
	}

	var selected *authn.Authn
	if err := huh.NewSelect[*authn.Authn]().Title("select the account to use").Options(options...).Value(&selected).WithTheme(huh.ThemeBase()).Run(); err != nil {
		return nil, err
	}

	return selected, nil
}
