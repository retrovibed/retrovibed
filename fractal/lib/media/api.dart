import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:fractal/media/media.pb.dart';
import 'package:http/http.dart' as http;
import 'package:fractal/httpx.dart' as httpx;
import 'dart:convert';
export 'package:fractal/media/media.pb.dart';

typedef FnMediaSearch =
    Future<MediaSearchResponse> Function(MediaSearchRequest req);

typedef FnDownloadSearch =
    Future<DownloadSearchResponse> Function(DownloadSearchRequest req);

Future<MediaSearchResponse> recent() async {
  final client = http.Client();
  return client.get(Uri.https(httpx.host(), "/m/recent")).then((v) {
    return Future.value(
      MediaSearchResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
    );
  });
}

abstract class mediasearch {
  static MediaSearchRequest request({int limit = 0}) =>
      MediaSearchRequest(limit: fixnum.Int64(limit));
  static MediaSearchResponse response({MediaSearchRequest? next}) =>
      MediaSearchResponse(next: next ?? request(limit: 100), items: []);
}

abstract class discoveredsearch {
  static DownloadSearchRequest request({int limit = 0}) =>
      DownloadSearchRequest(limit: fixnum.Int64(limit));
}

abstract class discovered {
  static Future<MediaSearchResponse> available(MediaSearchRequest req) async {
    final client = http.Client();
    return client
        .get(
          Uri.https(
            httpx.host(),
            "/d/available",
            jsonDecode(jsonEncode(req.toProto3Json())),
          ),
        )
        .then((v) {
          return Future.value(
            MediaSearchResponse.create()
              ..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<DownloadSearchResponse> downloading(
    DownloadSearchRequest req,
  ) async {
    final client = http.Client();
    return client
        .get(
          Uri.https(
            httpx.host(),
            "/d/downloading",
            jsonDecode(jsonEncode(req.toProto3Json())),
          ),
        )
        .then((v) {
          return Future.value(
            DownloadSearchResponse.create()
              ..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<DownloadBeginResponse> download(String id) async {
    final client = http.Client();
    return client
        .post(Uri.https(httpx.host(), "/d/${id}", null), body: jsonEncode({}))
        .then((v) {
          return DownloadBeginResponse.create()
            ..mergeFromProto3Json(jsonDecode(v.body));
        });
  }

  static Future<DownloadPauseResponse> pause(String id) async {
    final client = http.Client();
    return client
        .delete(Uri.https(httpx.host(), "/d/${id}", null), body: jsonEncode({}))
        .then((v) {
          return Future.value(
            DownloadPauseResponse.create()
              ..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }
}
