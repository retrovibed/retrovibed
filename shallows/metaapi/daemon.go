package metaapi

import (
	"github.com/retrovibed/retrovibed/internal/grpcx"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/internal/timex"
	"github.com/retrovibed/retrovibed/meta"
)

func NewDaemonFromMetaDaemon(mp meta.Daemon) (_ *Daemon, err error) {
	var p Daemon

	if err = grpcx.JSONDecode(langx.Clone(mp, timex.JSONSafeEncodeOption, timex.UTCEncodeOption), &p); err != nil {
		return nil, err
	}

	return &p, nil
}

func NewMetadaemonFromDaemon[T ~func(*meta.Daemon)](v *Daemon, options ...T) (e meta.Daemon, err error) {
	if err = grpcx.JSONEncode(v, &e); err != nil {
		return e, err
	}

	return langx.Clone(e, options...), nil
}
