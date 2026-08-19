package eggpgx

import (
	"context"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

const (
	EnvName  = "EG_GPG_KEYRING_NAME"
	EnvEmail = "EG_GPG_KEYRING_EMAIL"
	EnvSeed  = "EG_GPG_KEYRING_SEED"
)

type Option func(*option)
type options []Option

type option struct {
	name  string
	email string
	seed  string
}

// Generates default options from the environment
func Options() options {
	return options(nil).
		Name(egenv.String("", EnvName)).
		Email(egenv.String("", EnvEmail)).
		Seed(egenv.String("", EnvSeed))
}

func (t options) Name(v string) options {
	return append(t, func(o *option) {
		o.name = v
	})
}

func (t options) Email(v string) options {
	return append(t, func(o *option) {
		o.email = v
	})
}

func (t options) Seed(v string) options {
	return append(t, func(o *option) {
		o.seed = v
	})
}

func runtime(options ...Option) shell.Command {
	var opts option
	for _, opt := range options {
		opt(&opts)
	}

	return shell.Runtime().
		Environ(EnvEmail, opts.email).
		Environ(EnvName, opts.name).
		Environ(EnvSeed, opts.seed)
}

func Debug(options ...Option) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		runtime := runtime(options...)
		return shell.Op(
			runtime.New("env | grep -i 'EG_GPG_'"),
			runtime.New("gpg --list-keys"),
		)(ctx, o)
	}
}

func Seed(options ...Option) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		runtime := runtime(options...)
		return shell.Op(
			runtime.New("eg gpg keyring --name=\"${EG_GPG_KEYRING_NAME}\" --email=\"${EG_GPG_KEYRING_EMAIL}\" --seed=\"${EG_GPG_KEYRING_SEED}\""),
			runtime.New("GNUPGHOME=${HOME}/.gnupg gpg --import ${HOME}/.gnupg/private.asc"),
		)(ctx, o)
	}
}
