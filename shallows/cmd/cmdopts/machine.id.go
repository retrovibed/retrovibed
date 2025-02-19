package cmdopts

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/james-lawrence/deeppool/internal/x/errorsx"
	"github.com/james-lawrence/deeppool/internal/x/userx"
)

func MachineID() string {
	var (
		err error
		raw []byte
	)

	midpath := filepath.Join(userx.DefaultCacheDirectory(), "machine-id")

	if err = os.MkdirAll(filepath.Dir(midpath), 0700); err != nil {
		panic(errorsx.Wrapf(err, "unable to ensure cache directory for machine id %s", midpath))
	}

	if raw, err = os.ReadFile(midpath); err == nil {
		return strings.TrimSpace(string(raw))
	}

	// log.Println("failed to read a valid machine id, generating a random uuid", err)
	uid := uuid.Must(uuid.NewV7()).String()
	if err = os.WriteFile(midpath, []byte(uid), 0600); err == nil {
		return strings.TrimSpace(uid)
	}

	panic(errorsx.Wrapf(err, "failed to generate a machine id at %s", midpath))
}
