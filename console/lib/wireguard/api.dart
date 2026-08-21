import 'dart:convert';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:http/http.dart' as http;
import 'package:retrovibed/httpx.dart' as httpx;
import './meta.wireguard.pb.dart';

export './meta.wireguard.pb.dart';

typedef FnWireguardSearch = Future<WireguardSearchResponse> Function(WireguardSearchRequest req);

typedef FnUploadRequest =
    Future<WireguardUploadResponse> Function(
      http.MultipartRequest Function(http.MultipartRequest req) mkreq,
    );

typedef FnWireguardCurrent = Future<WireguardCurrentResponse> Function();
typedef FnWireguardUpdate = Future<WireguardUpdateResponse> Function(Wireguard wg, {List<httpx.Option> options});

abstract class wireguard {
  static WireguardSearchRequest request({int limit = 0, String query = ""}) =>
      WireguardSearchRequest(limit: ds.Int64(limit));
  static WireguardSearchResponse response({WireguardSearchRequest? next}) =>
      WireguardSearchResponse(next: next ?? request(limit: 100), items: []);

  static Future<WireguardSearchResponse> get(WireguardSearchRequest req) async {
    final client = http.Client();

    return client
        .get(
          Uri.https(
            httpx.host(),
            "/wireguard/",
            jsonDecode(jsonEncode(req.toProto3Json())),
          ),
          headers: {"Authorization": httpx.auto_bearer_host()},
        )
        .then(httpx.auto_error)
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(WireguardSearchResponse.create(), jsonDecode(v.body)),
          );
        });
  }

  static Future<WireguardDeleteResponse> delete(String id) async {
    final client = http.Client();
    return client
        .delete(
          Uri.https(httpx.host(), "/wireguard/${id}"),
          headers: {"Authorization": httpx.auto_bearer_host()},
        )
        .then(httpx.auto_error)
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(WireguardDeleteResponse.create(), jsonDecode(v.body)),
          );
        });
  }

  static Future<http.MultipartFile> uploadable(
    String path,
    String name,
    String mimetype,
  ) {
    return httpx.uploadable(path, name, mimetype);
  }

  static Future<WireguardUploadResponse> upload(
    http.MultipartRequest Function(http.MultipartRequest req) mkreq,
  ) async {
    final client = http.Client();
    final req = mkreq(
      http.MultipartRequest("POST", Uri.https(httpx.host(), "/wireguard/")),
    );
    req.headers["Authorization"] = httpx.auto_bearer_host();
    return client.send(req).then(httpx.auto_error).then((v) {
      return v.stream.bytesToString().then((s) {
        return Future.value(
          httpx.fromProto3JsonSafe(WireguardUploadResponse.create(), jsonDecode(s)),
        );
      });
    });
  }

  static Future<WireguardUpdateResponse> update(
    Wireguard wg, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .patch(
          Uri.https(httpx.host(), "/wireguard/${wg.id}"),
          body: jsonEncode(
            WireguardUpdateRequest(wireguard: wg).toProto3Json(),
          ),
          options: options,
        )
        .then(
          (v) => httpx.fromProto3JsonSafe(WireguardUpdateResponse.create(), jsonDecode(v.body)),
        );
  }

  // activate the specified wireguard configuration.
  static Future<WireguardTouchResponse> touch(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx.put(Uri.https(httpx.host(), "/wireguard/${id}"), options: options).then(httpx.auto_error).then((v) {
      return Future.value(
        httpx.fromProto3JsonSafe(WireguardTouchResponse.create(), jsonDecode(v.body)),
      );
    });
  }

  static Future<WireguardCurrentResponse> current() async {
    final client = http.Client();
    return client
        .get(
          Uri.https(httpx.host(), "/wireguard/current"),
          headers: {"Authorization": httpx.auto_bearer_host()},
        )
        .then(httpx.auto_error)
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(WireguardCurrentResponse.create(), jsonDecode(v.body)),
          );
        });
  }
}
