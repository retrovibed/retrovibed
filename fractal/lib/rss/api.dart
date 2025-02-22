import 'package:fractal/rss/rss.pb.dart';
import 'package:http/http.dart' as http;
import 'package:fractal/httpx.dart' as httpx;
import 'dart:convert';

typedef FnSearch = Future<FeedSearchResponse> Function(FeedSearchRequest req);

// Future<FeedSearchResponse> searchfake(
//   FeedSearchRequest req, {
//   Duration delay = const Duration(seconds: 3),
// }) async {
//   return Future.delayed(
//     delay,
//     () => FeedSearchResponse(
//       next: req,
//       items: [
//         Feed(
//           id: "1",
//           description: "feed1",
//           url: "https://example1",
//           autodownload: false,
//         ),
//         Feed(
//           id: "2",
//           description: "feed2",
//           url: "https://example2",
//           autodownload: false,
//         ),
//       ],
//     ),
//   );
// }

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
