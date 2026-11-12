package meta

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
)

type DaemonOption func(*Daemon)

func DaemonOptionMaybeID(v *Daemon) {
	v.ID = stringsx.First(v.ID, md5x.String(v.Hostname))
}

func DaemonOptionEnsureDescription(v *Daemon) {
	v.Description = stringsx.First(v.Description, v.Hostname, v.ID)
}

func DaemonOptionTestDefaults(v *Daemon) {
	v.ID = ""
	v.Hostname = fmt.Sprintf("%s.lan", uuidx.NewV7().String())
	v.Downloads = false
}

func DaemonOptionDownload(v *Daemon) {
	v.Downloads = true
}

func DaemonSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) DaemonScanner {
	return NewDaemonScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func DaemonSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(DaemonScannerStaticColumns)...).From("meta_daemons")
}

func DaemonFromHost() Daemon {
	return Daemon{
		CreatedAt:   time.Now(),
		Description: "local device",
		Hostname:    stringsx.First(fmt.Sprintf("%s:9998", errorsx.Zero(os.Hostname())), "localhost"),
		ID:          uuidx.Max().String(),
		UpdatedAt:   time.Now(),
	}
}
