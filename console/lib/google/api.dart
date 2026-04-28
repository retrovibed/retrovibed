import 'dart:convert';

import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/community/community.pb.dart';

export 'package:retrovibed/community/community.pb.dart' show YouTubeStatus;

class YouTube {
  static Future<YouTubeStatus> status({
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(httpx.host(), "/integrations/youtube/status"),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) => YouTubeStatus()..mergeFromProto3Json(jsonDecode(v.body)));
  }

  static Future<void> unlink({
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .delete(
          Uri.https(httpx.host(), "/integrations/youtube/token"),
          options: [httpx.Accept.json, ...options],
        )
        .then((_) {});
  }

  static Uri authUri({required String token}) {
    return Uri.https(httpx.metaendpoint(), "/oauth2/proxy/google/auth", {
      "token": token,
      "redirect_uri": Uri.https(httpx.host(), "/integrations/youtube/callback").toString(),
    });
  }
}
