import 'dart:convert';
import 'package:retrovibed/torrents/torrent.pb.dart' show DiscoverySettings;
import 'package:retrovibed/httpx.dart' as httpx;

export 'package:retrovibed/torrents/torrent.pb.dart' show DiscoverySettings;

abstract class discoveryapi {
  static Future<DiscoverySettings> get({
    List<httpx.Option> options = const [],
  }) async {
    return httpx.get(Uri.https(httpx.host(), "/s/discovery/", {}), options: options).then((v) {
      return Future.value(
        DiscoverySettings.create()..mergeFromProto3Json(jsonDecode(v.body)),
      );
    });
  }

  static Future<DiscoverySettings> create(
    DiscoverySettings req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/s/discovery/", {}),
          options: options,
          body: jsonEncode(req.toProto3Json()),
        )
        .then((v) {
          return Future.value(
            DiscoverySettings.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }
}
