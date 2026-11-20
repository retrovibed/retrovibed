import 'package:material_ui/material_ui.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/torrents.dart' as torrents;
import 'package:retrovibed/dht.dart' as dht;

class NetworkSettings extends StatelessWidget {
  final Future<torrents.TorrentSettings> Function({List<httpx.Option> options}) torrent;
  final Future<dht.DHTSettings> Function({List<httpx.Option> options}) dhts;

  const NetworkSettings({
    super.key,
    this.torrent = torrents.api.get,
    this.dhts = dht.api.get,
  });

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return Column(
      mainAxisSize: MainAxisSize.min,
      spacing: defaults.spacing / 2.5,
      children: [
        ds.Heading(Text("sharing")),
        torrents.Settings.future(
          torrent(options: [authn.request(authn.AuthzCache.meta(context))]),
          onChange: (v) => torrents.api.create(v, options: [authn.request(authn.AuthzCache.meta(context))]),
        ),
        ds.Heading(Text("dht")),
        dht.Settings.future(
          dhts(options: [authn.request(authn.AuthzCache.meta(context))]),
          onChange: (v) => dht.api.create(v, options: [authn.request(authn.AuthzCache.meta(context))]),
        ),
      ],
    );
  }
}
