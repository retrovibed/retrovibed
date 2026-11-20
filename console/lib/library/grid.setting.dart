import 'package:material_ui/material_ui.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/torrents.dart' as libtorrents;
import 'package:retrovibed/storage.dart' as libstorage;
import 'package:retrovibed/discovery.dart' as discovery;
import 'package:retrovibed/library/network.settings.dart';
import 'package:retrovibed/dht.dart' as libdht;

class _GridItem extends StatelessWidget {
  final Widget child;
  const _GridItem(this.child);

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return ds.Card(
      alignment: Alignment.topLeft,
      ClipRRect(borderRadius: defaults.borderRadius, child: child),
    );
  }
}

class GridSettings extends StatelessWidget {
  final Future<libtorrents.TorrentSettings> Function({
    List<httpx.Option> options,
  })
  torrents;
  final Future<libdht.DHTSettings> Function({
    List<httpx.Option> options,
  })
  dht;
  final Future<libstorage.StorageSettingsResponse> Function({
    List<httpx.Option> options,
  })
  storage;

  GridSettings({
    super.key,
    this.torrents = libtorrents.api.get,
    this.dht = libdht.api.get,
    this.storage = libstorage.api.get,
  });

  @override
  Widget build(BuildContext context) {
    final items = [
      _GridItem(discovery.Settings()),
      _GridItem(
        NetworkSettings(
          torrent: ({options = const []}) => httpx.withRetry(() => torrents(options: options)),
          dhts: ({options = const []}) => httpx.withRetry(() => dht(options: options)),
        ),
      ),
      _GridItem(
        libstorage.MinimalSettings.future(
          httpx.withRetry(
            () => storage(
              options: [authn.request(authn.AuthzCache.meta(context))],
            ),
          ),
          onChange: (v) {
            return libstorage.api.create(
              v,
              options: [authn.request(authn.AuthzCache.meta(context))],
            );
          },
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
