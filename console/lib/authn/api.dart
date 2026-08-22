import 'dart:convert';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/meta/meta.authn.pb.dart';

export 'package:retrovibed/meta/meta.account.pb.dart';
export 'package:retrovibed/meta/meta.authn.pb.dart';

Future<Session> current(String token) async {
  return httpx
      .get(
        Uri.https(httpx.metaendpoint(), "/authn/current"),
        options: [httpx.Request.authorization(httpx.bearer(token))],
      )
      .then((v) {
        return httpx.fromProto3JsonSafe(Session.create(), jsonDecode(v.body));
      });
}

Future<Session> otp({List<httpx.Option> options = const []}) async {
  return httpx
      .get(
        Uri.https(httpx.metaendpoint(), "/authn/otp"),
        options: options,
      )
      .then((v) {
        return httpx.fromProto3JsonSafe(Session.create(), jsonDecode(v.body));
      });
}

Future<Authed> ssh() async {
  final token = httpx.deeppool_oauth2_bearer();
  return httpx
      .post(
        Uri.https(httpx.metaendpoint(), "/authn/ssh"),
        options: [httpx.Request.authorization(token)],
      )
      .then((v) {
        return httpx.fromProto3JsonSafe(Authed.create(), jsonDecode(v.body));
      });
}

Future<Session> signup() {
  return httpx
      .post(
        Uri.https(httpx.metaendpoint(), "/authn/signup"),
        options: [httpx.Request.authorization(httpx.deeppool_oauth2_bearer())],
      )
      .then((v) {
        return httpx.fromProto3JsonSafe(Session.create(), jsonDecode(v.body));
      });
}

extension sessions on Session {
  bool get isZero {
    return !profile.hasId();
  }
}
