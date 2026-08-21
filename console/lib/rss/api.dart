import 'dart:convert';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/timex.dart' as timex;
import 'package:retrovibed/uuidx.dart' as uuidx;
import './rss.pb.dart';
export './rss.pb.dart';

typedef FnSearch = Future<FeedSearchResponse> Function(FeedSearchRequest req, {List<httpx.Option> options});

Future<FeedSearchResponse> search(
  FeedSearchRequest req, {
  List<httpx.Option> options = const [],
}) async {
  return httpx
      .get(
        Uri.https(
          httpx.host(),
          "/rss/",
          jsonDecode(jsonEncode(req.toProto3Json())),
        ),
        options: options,
      )
      .then((v) {
        return Future.value(
          httpx.fromProto3JsonSafe(FeedSearchResponse.create(), jsonDecode(v.body)),
        );
      });
}

Future<FeedCreateResponse> create(
  FeedCreateRequest req, {
  List<httpx.Option> options = const [],
}) async {
  return httpx
      .post(
        Uri.https(httpx.host(), "/rss/", null),
        body: jsonEncode(req.toProto3Json()),
        options: options,
      )
      .then((v) {
        return httpx.fromProto3JsonSafe(FeedCreateResponse.create(), jsonDecode(v.body));
      });
}

Future<FeedCreateResponse> refresh(
  FeedCreateRequest req, {
  List<httpx.Option> options = const [],
}) async {
  req.feed..nextCheck = timex.formatISO8601(timex.now());
  req.feed..digest = uuidx.min();
  return httpx
      .post(
        Uri.https(httpx.host(), "/rss/", null),
        body: jsonEncode(req.toProto3Json()),
        options: options,
      )
      .then((v) {
        return httpx.fromProto3JsonSafe(FeedCreateResponse.create(), jsonDecode(v.body));
      });
}

Future<FeedDeleteResponse> delete(
  String id, {
  List<httpx.Option> options = const [],
}) async {
  return httpx.delete(Uri.https(httpx.host(), "/rss/${id}"), options: options).then((v) {
    return httpx.fromProto3JsonSafe(FeedDeleteResponse.create(), jsonDecode(v.body));
  });
}
