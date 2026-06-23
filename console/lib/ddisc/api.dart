import 'dart:convert';

import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/meta/meta.discovery.pb.dart';

export 'package:retrovibed/meta/meta.discovery.pb.dart';

typedef FnDiscoveryDiagnostics = Future<DiscoveryMetricsResponse> Function({List<httpx.Option> options});

abstract class diagnostics {
  static Future<DiscoveryMetricsResponse> get({
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(httpx.host(), "/diagnostics/discovery/"),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) => DiscoveryMetricsResponse()..mergeFromProto3Json(jsonDecode(v.body)));
  }
}
