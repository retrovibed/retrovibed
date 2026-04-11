import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:retrovibed/httpx.dart' as httpx;
import './meta.authz.pb.dart';

export './meta.daemon.pb.dart';
export './meta.authz.pb.dart';

Future<http.Response> healthz() {
  return httpx.get(Uri.https(httpx.metaendpoint(), "/healthz"));
}

Future<AuthzResponse> authz({List<Future<httpx.Request> Function(httpx.Request)> options = const []}) {
  return httpx.get(Uri.https(httpx.metaendpoint(), "/m/authz/"), options: options).then((r) {
    return AuthzResponse.create()..mergeFromProto3Json(jsonDecode(r.body));
  });
}
