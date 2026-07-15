import 'dart:convert';

import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/meta/meta.discovery.pb.dart';
import 'ddisc.discovery.pb.dart';

export 'ddisc.discovery.pb.dart';
export 'package:retrovibed/meta/meta.discovery.pb.dart';

typedef FnDiscoveryDiagnostics = Future<DiscoveryMetricsResponse> Function({List<httpx.Option> options});

abstract class sources {
  static const String discovered = "538416cf-3bc5-9332-670a-f4cae9485ebe"; // md5 of RecommendationSourceDiscovered
}

abstract class api {
  static Future<DiscoveryDownloadResponse> download(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/ddisc/discovery/${id}", {}),
          options: options,
          body: jsonEncode(DiscoveryDownloadRequest().toProto3Json()),
        )
        .then((v) {
          return Future.value(
            DiscoveryDownloadResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }
}

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
