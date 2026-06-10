import 'dart:convert';

import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/meta/meta.network.pb.dart';

export 'package:retrovibed/meta/meta.network.pb.dart';

typedef FnNetworkDiagnostics = Future<NetworkMetricsResponse> Function({List<httpx.Option> options});

abstract class network {
  static Future<NetworkMetricsResponse> get({
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(httpx.host(), "/diagnostics/network"),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) => NetworkMetricsResponse()..mergeFromProto3Json(jsonDecode(v.body)));
  }
}
