import 'dart:convert';
import 'package:retrovibed/torrents/torrent.pb.dart';
import 'package:retrovibed/httpx.dart' as httpx;

export 'package:retrovibed/torrents/torrent.pb.dart';

abstract class api {
  static Future<TorrentSettings> get({
    List<httpx.Option> options = const [],
  }) async {
    return httpx.get(Uri.https(httpx.host(), "/s/torrents/", {}), options: options).then((v) {
      return Future.value(
        httpx.fromProto3JsonSafe(TorrentSettings.create(), jsonDecode(v.body)),
      );
    });
  }

  static Future<TorrentSettings> create(
    TorrentSettings req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/s/torrents/", {}),
          options: options,
          body: jsonEncode(req.toProto3Json()),
        )
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(TorrentSettings.create(), jsonDecode(v.body)),
          );
        });
  }

  static Future<TorrentSettings> delete({
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .delete(
          Uri.https(httpx.host(), "/s/torrents/", {}),
          options: options,
        )
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(TorrentSettings.create(), jsonDecode(v.body)),
          );
        });
  }
}
