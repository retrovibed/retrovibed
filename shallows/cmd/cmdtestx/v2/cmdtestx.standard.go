package cmdtestx

import (
	"path/filepath"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
)

// Standard returns the options that mirror the real cli's root: the shared tls, ssh identity and
// endpoint flags are embedded into the command and bound for injection, with the standard
// vars and mappers applied.
func Standard(t *testing.T) kong.Option {
	t.Helper()

	var (
		tls      = cmdopts.TLSConfig{PoolDir: filepath.Join(t.TempDir(), "tls.d")}
		sshid    cmdopts.SSHID
		endpoint cmdopts.Endpoint
	)

	// kong doesn't run hooks for embedded structs, the trust pool has to be created up front.
	if err := tls.AfterApply(); err != nil {
		t.Fatal(err)
	}

	return kong.OptionFunc(func(k *kong.Kong) error {
		for _, opt := range []kong.Option{
			kong.Embed(&tls),
			kong.Embed(&sshid),
			kong.Embed(&endpoint),
			kong.Bind(&tls),
			kong.Bind(&sshid),
			kong.Bind(&endpoint),
			kong.Vars{
				"vars_private_key":                  env.PrivateKeyPath(),
				"vars_user_configuration_directory": t.TempDir(),
				"vars_user_cache_directory":         t.TempDir(),
				"env_http_endpoint":                 env.Endpoint,
			},
			kong.NamedMapper("durationinf", kong.MapperFunc(cmdopts.ParseDurationInf)),
			kong.NamedMapper("envvar", kong.MapperFunc(cmdopts.ParseEnviron)),
		} {
			if err := opt.Apply(k); err != nil {
				return err
			}
		}

		return nil
	})
}
