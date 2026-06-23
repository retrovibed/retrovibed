import 'dart:convert';

import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/meta/meta.dht.pb.dart';

export 'package:retrovibed/meta/meta.dht.pb.dart';

typedef FnDHTDiagnostics = Future<DHTMetricsResponse> Function({List<httpx.Option> options});

abstract class diagnostics {
  static Future<DHTMetricsResponse> get({
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(httpx.host(), "/diagnostics/dht/"),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) => DHTMetricsResponse()..mergeFromProto3Json(jsonDecode(v.body)));
  }
}
