import 'dart:convert';
import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:http/http.dart' as http;
import 'package:console/httpx.dart' as httpx;
import './meta.daemon.pb.dart';
import './meta.authz.pb.dart';

export './meta.daemon.pb.dart';
export './meta.authz.pb.dart';

Future<http.Response> healthz() async {
  return http.Client().get(Uri.https(httpx.host(), "/healthz"));
}

Future<AuthzResponse> authz(Token current) {
  return http.Client()
      .get(
        Uri.https(httpx.host(), "/meta/authz/"),
        headers: {"Authorization": httpx.auto_bearer_host()},
      )
      .then(httpx.auto_error)
      .then((v) {
        return AuthzResponse.create()..mergeFromProto3Json(jsonDecode(v.body));
      });
}

abstract class daemons {
  static DaemonSearchRequest request({int limit = 0}) =>
      DaemonSearchRequest(limit: fixnum.Int64(limit));
  static DaemonSearchResponse response({DaemonSearchRequest? next}) =>
      DaemonSearchResponse(next: next ?? request(limit: 128), items: []);

  static Future<DaemonSearchResponse> search(DaemonSearchRequest req) async {
    return http.Client()
        .get(
          Uri.https(
            httpx.host(),
            "/meta/d/",
            jsonDecode(jsonEncode(req.toProto3Json())),
          ),
          headers: {"Authorization": httpx.auto_bearer()},
        )
        .then(httpx.auto_error)
        .then((v) {
          return DaemonSearchResponse.create()
            ..mergeFromProto3Json(jsonDecode(v.body));
        });
  }

  static Future<DaemonCreateResponse> create(DaemonCreateRequest req) async {
    return http.Client()
        .post(
          Uri.https(httpx.host(), "/meta/d/"),
          headers: {"Authorization": httpx.auto_bearer()},
          body: jsonEncode(req.toProto3Json()),
        )
        .then(httpx.auto_error)
        .then((v) {
          return Future.value(
            DaemonCreateResponse.create()
              ..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<DaemonLookupResponse> latest() async {
    return http.Client()
        .get(
          Uri.https(httpx.host(), "/meta/d/latest"),
          headers: {"Authorization": httpx.auto_bearer()},
        )
        .then((v) {
          return Future.value(
            DaemonLookupResponse.create()
              ..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }
}
