package cmdmeta

import (
	"log"

	"os"

	"github.com/charmbracelet/huh"
	"github.com/mattn/go-isatty"

	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/me"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
)

type CloudMeEdit struct {
	Display string `name:"display" help:"new display name, prompted for when neither --display nor --email is provided"`
	Email   string `name:"email" help:"new email address"`
}

func (t CloudMeEdit) Run(gctx *cmdopts.Global, id *cmdopts.SSHID) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return err
	}

	c, session, err := cloudClient(gctx.Context, gctx, signer)
	if err != nil {
		return err
	}

	// the api replaces the display and email with what is sent, anything not specified keeps its current value.
	display, email := stringsx.First(t.Display, session.Profile.GetDisplay()), stringsx.First(t.Email, session.Profile.GetEmail())

	// nothing specified, prompt with the current values.
	if stringsx.Blank(t.Display) && stringsx.Blank(t.Email) {
		if !isatty.IsTerminal(os.Stdin.Fd()) || !isatty.IsTerminal(os.Stdout.Fd()) {
			return errorsx.New("nothing to edit, specify --display or --email when not running interactively")
		}

		form := huh.NewForm(huh.NewGroup(
			huh.NewInput().Title("display name").Value(&display),
			huh.NewInput().Title("email").Value(&email),
		))

		if err = form.Run(); err != nil {
			return err
		}
	}

	p, err := me.UpdateProfile(gctx.Context, c, display, email)
	if err != nil {
		if me.BadRequest(err) {
			return errorsx.Wrap(err, "the profile was rejected, check the email address")
		}

		return err
	}

	log.Println("updated", p.Id, p.Display, p.Email)

	return nil
}
