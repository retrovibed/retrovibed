import 'package:http/http.dart' as http;
import 'package:fractal/httpx.dart' as httpx;
import 'dart:convert';
import './rss.pb.dart';
export './rss.pb.dart';

typedef FnSearch = Future<FeedSearchResponse> Function(FeedSearchRequest req);

Future<FeedSearchResponse> search(FeedSearchRequest req) async {
  final client = http.Client();
  return client
      .get(
        Uri.https(
          httpx.host(),
          "/rss/",
          jsonDecode(jsonEncode(req.toProto3Json())),
        ),
      )
      .then((v) {
        return Future.value(
          FeedSearchResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
        );
      });
}

Future<FeedCreateResponse> create(FeedCreateRequest req) async {
  final client = http.Client();
  return client
      .post(
        Uri.https(httpx.host(), "/rss/", null),
        body: jsonEncode(req.toProto3Json()),
      )
      .then((v) {
        return FeedCreateResponse.create()
          ..mergeFromProto3Json(jsonDecode(v.body));
      });
}
