import 'dart:convert';

import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/meta/meta.network.pb.dart';

export 'package:retrovibed/meta/meta.network.pb.dart' show NetworkMetricsResponse;

abstract class NetworkDiagnostics {
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
