import 'dart:convert';
import 'package:retrovibed/torrents/torrent.pb.dart' show DHTSettings;
import 'package:retrovibed/httpx.dart' as httpx;

export 'package:retrovibed/torrents/torrent.pb.dart' show DHTSettings, GeneratedMessageGenericExtensions;

abstract class api {
  static Future<DHTSettings> get({
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(httpx.host(), "/s/dht/", {}),
          options: [httpx.Content.urlencoded, httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(DHTSettings.create(), jsonDecode(v.body)),
          );
        });
  }

  static Future<DHTSettings> create(
    DHTSettings req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/s/dht/", {}),
          options: [httpx.Content.json, httpx.Accept.json, ...options],
          body: jsonEncode(req.toProto3Json()),
        )
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(DHTSettings.create(), jsonDecode(v.body)),
          );
        });
  }
}
