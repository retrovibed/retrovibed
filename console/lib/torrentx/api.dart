import 'dart:convert';

import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/meta/meta.torrent.pb.dart';

export 'package:retrovibed/meta/meta.torrent.pb.dart';

typedef FnTorrentDiagnostics = Future<TorrentMetricsResponse> Function({List<httpx.Option> options});

abstract class diagnostics {
  static Future<TorrentMetricsResponse> get({
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(httpx.host(), "/diagnostics/torrent/"),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) => TorrentMetricsResponse()..mergeFromProto3Json(jsonDecode(v.body)));
  }
}
