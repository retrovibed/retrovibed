package meta

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/atomicx"
	"github.com/retrovibed/retrovibed/retroapi/fsx"
	"github.com/retrovibed/retrovibed/shallows/dnscache"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/netx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
	"github.com/retrovibed/retrovibed/shallows/internal/wireguardx"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
)

type dialerproxy interface {
	Store(d netx.Dialer)
}

type dnsproxy interface {
	dnscache.Resolver
	Store(c dnscache.Resolver)
}

func connect(ctx context.Context, path string, wg Wireguard, dp dialerproxy, dnsp dnsproxy) (wgdev *device.Device, err error) {
	var (
		wgnet *netstack.Net
	)

	wcfg, err := wireguardx.Parse(path)
	if err != nil {
		return nil, errorsx.Wrapf(err, "unable to parse wireguard configuration: %s", path)
	}

	log.Println("loaded wireguard configuration", path)

	if wgnet, wgdev, err = torrentx.WireguardSocket(ctx, wcfg); err != nil {
		return nil, errorsx.Wrap(err, "unable to setup wireguard tunnel")
	}

	dp.Store(wireguardx.DefaultDialer(wgnet, dnsp))
	dnsp.Store(
		dnscache.New(
			wireguardx.HostLookupAdapter(wgnet),
			dnscache.CacheOptionRateLimit(wg.RateLimitDNS),
		),
	)

	return wgdev, nil
}

func NewWireguardNetworkBackgroundAuto(ctx context.Context, dir string, ring WireguardRing, q sqlx.Queryer, dp dialerproxy, dnsp dnsproxy) error {
	a := asyncx.NewWakeup(ctx)

	if err := fsx.Watch(ctx, a.Cond, dir); err != nil {
		return err
	}

	firstrun := func() (wgdev *device.Device, err error) {
		var (
			wg Wireguard
		)

		if err := WireguardCurrent(ctx, q, uint32(ring)).Scan(&wg); errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else if err != nil {
			return nil, err
		}

		path := fsx.DirVirtual(dir).Path(uuid.FromStringOrNil(wg.ID).String())
		if !fsx.Exists(path) {
			return nil, nil
		}

		return connect(ctx, path, wg, dp, dnsp)
	}

	initdev, err := firstrun()
	if err != nil {
		return err
	}

	devcache := atomicx.PointerPtr(initdev)
	asyncx.Background(ctx, a, func(_ctx context.Context) error {
		var (
			wg Wireguard
		)

		if err := WireguardCurrent(_ctx, q, uint32(ring)).Scan(&wg); errors.Is(err, sql.ErrNoRows) {
			return nil
		} else if err != nil {
			return err
		}

		path := fsx.DirVirtual(dir).Path(uuid.FromStringOrNil(wg.ID).String())
		if !fsx.Exists(path) {
			return nil
		}

		wgdev, err := connect(ctx, path, wg, dp, dnsp)
		if err != nil {
			return errorsx.Wrap(err, "failed to setup wireguard")
		}

		if old := devcache.Swap(wgdev); old != nil {
			old.Close()
		}

		log.Println("ring updated", ring, spew.Sdump(wg))
		return nil
	})

	return nil
}
