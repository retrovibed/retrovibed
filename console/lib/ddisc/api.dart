import 'dart:async';
import 'dart:convert';

import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media/ddisc.locate.pb.dart';
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
    Discovery? discovery,
    bool autodownload = true,
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/ddisc/discovery/download", {}),
          options: options,
          body: jsonEncode(
            DiscoveryDownloadRequest(
              discovery: discovery ?? (Discovery.create()..id = id),
              autodownload: autodownload,
            ).toProto3Json(),
          ),
        )
        .then((v) {
          return Future.value(
            DiscoveryDownloadResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  // locate runs a live, ephemeral network search and streams each ranked
  // candidate as it's found; the server always re-sends the best candidate
  // as the final message before closing.
  static Future<Stream<Discovery>> locate(
    Locate req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .websocket(
          Uri.https(httpx.host(), "/ddisc/discovery/locate", httpx.params(req.toProto3Json())),
          options: options,
        )
        .then((socket) {
          return socket.transform(
            StreamTransformer.fromHandlers(
              handleData: (data, sink) {
                if (data is List<int>) {
                  sink.add(Discovery.create()..mergeFromProto3Json(jsonDecode(utf8.decode(data))));
                } else {
                  sink.addError('deserialization failed data: $data');
                }
              },
            ),
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
