import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/torrents.dart' as torrents_pkg;
import 'package:retrovibed/wireguard.dart' as wireguard;
import 'package:retrovibed/storage.dart' as storage_pkg;
import 'package:retrovibed/discovery.dart' as discovery;

class _GridItem extends StatelessWidget {
  final Widget child;
  final List<Widget> leading;
  final List<Widget> trailing;
  const _GridItem(
    this.child, {
    this.leading = const [],
    this.trailing = const [],
  });

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return ds.Card(
      alignment: Alignment.topLeft,
      leading: leading,
      ClipRRect(borderRadius: defaults.borderRadius, child: child),
      trailing: trailing,
    );
  }
}

class GridSettings extends StatelessWidget {
  final Future<torrents_pkg.TorrentSettings> Function({
    List<httpx.Option> options,
  })
  torrents;
  final Future<storage_pkg.StorageSettingsResponse> Function({
    List<httpx.Option> options,
  })
  storage;
  final wireguard.FnWireguardCurrent wgcurrent;
  final wireguard.FnWireguardUpdate wgupdate;

  GridSettings({
    super.key,
    this.torrents = torrents_pkg.api.get,
    this.storage = storage_pkg.api.get,
    this.wgcurrent = wireguard.wireguard.current,
    this.wgupdate = wireguard.wireguard.update,
  });

  @override
  Widget build(BuildContext context) {
    final items = [
      _GridItem(discovery.Settings()),
      _GridItem(
        leading: [ds.Heading(Text("sharing"))],
        torrents_pkg.Settings.future(
          httpx.withRetry(
            () => torrents(options: [authn.request(authn.AuthzCache.meta(context))]),
          ),
          onChange: (v) {
            return torrents_pkg.api.create(
              v,
              options: [authn.request(authn.AuthzCache.meta(context))],
            );
          },
        ),
        trailing: const [],
      ),
      _GridItem(
        storage_pkg.MinimalSettings.future(
          httpx.withRetry(
            () => storage(
              options: [authn.request(authn.AuthzCache.meta(context))],
            ),
          ),
          onChange: (v) {
            return storage_pkg.api.create(
              v,
              options: [authn.request(authn.AuthzCache.meta(context))],
            );
          },
        ),
      ),
      _GridItem(
        leading: [ds.Heading(Text("network"))],
        wireguard.Settings.future(
          wgcurrent().then((r) => r.wireguard),
          onChange: (v) => wgupdate(v, options: [authn.request(authn.AuthzCache.meta(context))]).then((r) => r.wireguard),
        ),
      ),
    ];

    return LayoutBuilder(
      builder: (context, constraints) {
        const minheight = 750.0;
        const minCrossAxisExtent = 400.0;
        const maxCrossAxisExtent = 508.0;

        // returns (width, height) for a grid item at the given max extent,
        // where height is always minheight.
        final calc = (double extent) {
          final columns = (constraints.maxWidth / extent).ceil();
          final w = constraints.maxWidth / columns;
          return (w, minheight);
        };

        final (minW, minH) = calc(minCrossAxisExtent);
        final (maxW, maxH) = calc(maxCrossAxisExtent);
        final maxRatio = maxW / maxH;
        final compact = (minW > minH) && (maxW > maxH);

        if (compact) {
          return Column(mainAxisSize: MainAxisSize.min, children: items);
        }

        return ds.Grid(
          (context, i) => i,
          padding: EdgeInsets.zero,
          maxCrossAxisExtent: maxCrossAxisExtent,
          aspectRatio: maxRatio,
          children: items,
        );
      },
    );
  }
}
