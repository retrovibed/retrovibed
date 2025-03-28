import 'dart:convert';
import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:console/media/media.pb.dart';
import 'package:http/http.dart' as http;
import 'package:console/httpx.dart' as httpx;
export 'package:console/media/media.pb.dart';

typedef FnMediaSearch =
    Future<MediaSearchResponse> Function(MediaSearchRequest req);

typedef FnDownloadSearch =
    Future<DownloadSearchResponse> Function(DownloadSearchRequest req);

typedef FnUploadRequest =
    Future<MediaUploadResponse> Function(
      http.MultipartRequest Function(http.MultipartRequest req) mkreq,
    );

Future<MediaSearchResponse> recent() async {
  final client = http.Client();
  return client.get(Uri.https(httpx.host(), "/m/recent")).then((v) {
    return Future.value(
      MediaSearchResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
    );
  });
}

abstract class media {
  static final client = http.Client();
  static MediaSearchRequest request({int limit = 0}) =>
      MediaSearchRequest(limit: fixnum.Int64(limit));
  static MediaSearchResponse response({MediaSearchRequest? next}) =>
      MediaSearchResponse(next: next ?? request(limit: 100), items: []);

  static Future<MediaSearchResponse> get(MediaSearchRequest req) async {
    return client
        .get(
          Uri.https(
            httpx.host(),
            "/m/",
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

  static String download_uri(String id) {
    return Uri.https(httpx.host(), "/m/${id}").toString();
  }

  static Future<MediaDeleteResponse> delete(String id) async {
    return client.delete(Uri.https(httpx.host(), "/m/${id}")).then((v) {
      return Future.value(
        MediaDeleteResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
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

  static Future<MediaUploadResponse> upload(
    http.MultipartRequest Function(http.MultipartRequest req) mkreq,
  ) async {
    final req = mkreq(
      http.MultipartRequest("POST", Uri.https(httpx.host(), "/m/")),
    );

    return client.send(req).then((v) {
      return v.stream.bytesToString().then((s) {
        return Future.value(
          MediaUploadResponse.create()..mergeFromProto3Json(jsonDecode(s)),
        );
      });
    });
  }
}

abstract class discoveredsearch {
  static DownloadSearchRequest request({int limit = 0}) =>
      DownloadSearchRequest(limit: fixnum.Int64(limit));
}

abstract class discovered {
  static final client = http.Client();
  static Future<MediaSearchResponse> available(MediaSearchRequest req) async {
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

  static Future<MediaUploadResponse> upload(
    http.MultipartRequest Function(http.MultipartRequest req) mkreq,
  ) async {
    final req = mkreq(
      http.MultipartRequest("POST", Uri.https(httpx.host(), "/d/")),
    );

    return client.send(req).then((v) {
      return v.stream.bytesToString().then((s) {
        return Future.value(
          MediaUploadResponse.create()..mergeFromProto3Json(jsonDecode(s)),
        );
      });
    });
  }

  static Future<DownloadBeginResponse> download(String id) async {
    return client
        .post(Uri.https(httpx.host(), "/d/${id}", null), body: jsonEncode({}))
        .then((v) {
          return DownloadBeginResponse.create()
            ..mergeFromProto3Json(jsonDecode(v.body));
        });
  }

  static Future<DownloadPauseResponse> pause(String id) async {
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
