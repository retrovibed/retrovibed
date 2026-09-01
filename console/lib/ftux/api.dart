import 'dart:convert';
import 'package:retrovibed/httpx.dart' as httpx;
import 'ftux.pb.dart';

export 'ftux.pb.dart';

class ftux {
  static Future<CommunitySuggestions> defaults({List<httpx.Option> options = const []}) async {
    return httpx
        .get(Uri.https(httpx.host(), "/ftux/communities"), options: [httpx.Accept.json, ...options])
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(CommunitySuggestions.create(), jsonDecode(v.body)),
          );
        });
  }

  static Future<SubscribeCommunitiesResponse> subscribe(
    List<String> ids, {
    List<httpx.Option> options = const [],
  }) async {
    final req = SubscribeCommunitiesRequest(communityId: ids);
    return httpx
        .post(
          Uri.https(httpx.host(), "/ftux/communities"),
          body: jsonEncode(req.toProto3Json()),
          options: [httpx.Accept.json, httpx.Content.json, ...options],
        )
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(SubscribeCommunitiesResponse.create(), jsonDecode(v.body)),
          );
        });
  }
}
