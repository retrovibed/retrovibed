package metaapi

import (
	"github.com/retrovibed/retrovibed/internal/grpcx"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/internal/timex"
	"github.com/retrovibed/retrovibed/meta"
)

func NewWireguardFromMetaWireguard(mp meta.Wireguard) (_ *Wireguard, err error) {
	var p Wireguard

	if err = grpcx.JSONDecode(langx.Clone(mp, timex.JSONSafeEncodeOption, timex.UTCEncodeOption), &p); err != nil {
		return nil, err
	}

	return &p, nil
}

func NewMetaWireguardFromWireguard[T ~func(*meta.Wireguard)](v *Wireguard, options ...T) (e meta.Wireguard, err error) {
	if err = grpcx.JSONEncode(v, &e); err != nil {
		return e, err
	}

	return langx.Clone(e, options...), nil
}
