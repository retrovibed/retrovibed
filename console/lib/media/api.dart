import 'dart:convert';
import 'dart:async';
import 'package:flutter/foundation.dart';
import 'package:qs_dart/qs_dart.dart' as qs;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media/media.pb.dart';
import 'package:retrovibed/media/content.addressable.storage.pb.dart' as cas;
import 'package:http/http.dart' as http;
import 'package:retrovibed/httpx.dart' as httpx;

export 'package:retrovibed/media/media.pb.dart';

typedef FnMediaSearch =
    Future<MediaSearchResponse> Function(
      MediaSearchRequest req, {
      List<httpx.Option> options,
    });

typedef FnMediaRandom =
    Future<MediaFindResponse> Function(
      MediaSearchRequest req, {
      List<httpx.Option> options,
    });

typedef FnDownloadSearch =
    Future<DownloadSearchResponse> Function(
      DownloadSearchRequest req, {
      List<httpx.Option> options,
    });

typedef FnDownloadWatch =
    Future<Stream<Download>> Function(
      String id, {
      List<httpx.Option> options,
    });

typedef FnUploadRequest =
    Future<MediaUploadResponse> Function(
      http.MultipartRequest Function(http.MultipartRequest req) mkreq,
    );

abstract class media {
  static MediaSearchRequest request({
    int limit = 0,
    String query = "",
    List<String> mimetypes = const [],
  }) => MediaSearchRequest(limit: ds.Int64(limit), mimetypes: mimetypes);
  static MediaSearchResponse response({MediaSearchRequest? next}) =>
      MediaSearchResponse(next: next ?? request(limit: 100), items: []);

  static Future<MediaSearchResponse> search(
    MediaSearchRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(
            httpx.host(),
            "/m/",
          ).replace(query: qs.encode(req.toProto3Json())),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            MediaSearchResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<MediaUpdateResponse> metadatasync(
    String id,
    Media upd, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/m/${id}/metadatasync"),
          body: jsonEncode(MediaUpdateRequest(media: upd).toProto3Json()),
          options: [httpx.Accept.json, httpx.Content.json, ...options],
        )
        .then((v) {
          return Future.value(
            MediaUpdateResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<MediaFindResponse> random(
    MediaSearchRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(
            httpx.host(),
            "/m/random",
          ).replace(query: qs.encode(req.toProto3Json())),
          options: options,
        )
        .then((v) {
          return Future.value(
            MediaFindResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<http.StreamedResponse> download(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx.send(Uri.https(httpx.host(), "/m/${id}"), options: options);
  }

  static String download_uri(String id) {
    return Uri.https(httpx.host(), "/m/${id}").toString();
  }

  static Future<MediaDeleteResponse> delete(String id) async {
    final client = http.Client();
    return client
        .delete(
          Uri.https(httpx.host(), "/m/${id}"),
          headers: {"Authorization": httpx.auto_bearer_host()},
        )
        .then(httpx.auto_error)
        .then((v) {
          return Future.value(
            MediaDeleteResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<MediaUpdateResponse> update(
    String id,
    Media upd, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/m/${id}"),
          options: options,
          body: jsonEncode(MediaUpdateRequest(media: upd).toProto3Json()),
        )
        .then((v) {
          return Future.value(
            MediaUpdateResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<cas.MediaDeleteResponse> unarchive(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx.delete(Uri.https(httpx.metaendpoint(), "/m/${id}"), options: options).then((v) {
      return Future.value(
        cas.MediaDeleteResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
      );
    });
  }

  static Future<http.MultipartFile> uploadable(
    String path,
    String name,
    String mimetype, {
    ValueNotifier<int>? progress,
  }) {
    return httpx.uploadable(path, name, mimetype, progress: progress);
  }

  static Future<MediaUploadResponse> upload(http.MultipartRequest Function(http.MultipartRequest req) mkreq) async {
    final client = http.Client();
    final req = mkreq(
      http.MultipartRequest("POST", Uri.https(httpx.host(), "/m/")),
    );
    req.headers["Authorization"] = httpx.auto_bearer_host();

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
  static DownloadSearchRequest request({int limit = 0}) => DownloadSearchRequest(limit: ds.Int64(limit));
  static DownloadSearchResponse response({DownloadSearchRequest? next}) =>
      DownloadSearchResponse(next: next ?? request(limit: 100), items: []);
}

abstract class discovered {
  static Future<DownloadSearchResponse> available(
    DownloadSearchRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(
            httpx.host(),
            "/d/available",
            httpx.params(req.toProto3Json()),
          ),
          options: options,
        )
        .then((v) {
          return Future.value(
            DownloadSearchResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<DownloadSearchResponse> downloading(
    DownloadSearchRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(
            httpx.host(),
            "/d/downloading",
            httpx.params(req.toProto3Json()),
          ),
          options: options,
        )
        .then((v) {
          return Future.value(
            DownloadSearchResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<MagnetCreateResponse> magnet(
    MagnetCreateRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/d/magnet"),
          body: jsonEncode(req.toProto3Json()),
          options: options,
        )
        .then((v) {
          return Future.value(
            MagnetCreateResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<MediaUploadResponse> upload(
    http.MultipartRequest Function(http.MultipartRequest req) mkreq,
  ) async {
    final client = http.Client();
    final req = mkreq(
      http.MultipartRequest("POST", Uri.https(httpx.host(), "/d/")),
    );
    req.headers["Authorization"] = httpx.auto_bearer_host();

    return client.send(req).then((v) {
      return v.stream.bytesToString().then((s) {
        return Future.value(
          MediaUploadResponse.create()..mergeFromProto3Json(jsonDecode(s)),
        );
      });
    });
  }

  static Future<DownloadBeginResponse> download(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/d/${id}", null),
          body: jsonEncode({}),
          options: options,
        )
        .then(httpx.auto_error)
        .then((v) {
          return DownloadBeginResponse.create()..mergeFromProto3Json(jsonDecode(v.body));
        });
  }

  static Future<MetadataSyncResponse> metadatasync(
    String id,
    Media upd, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/d/${id}/metadatasync"),
          body: jsonEncode(MetadataSyncRequest(media: upd).toProto3Json()),
          options: [httpx.Accept.json, httpx.Content.json, ...options],
        )
        .then((v) {
          return Future.value(
            MetadataSyncResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<Stream<Download>> watch(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .websocket(
          Uri.https(httpx.host(), "/d/${id}/socket", null),
          options: options,
        )
        .then((socket) {
          socket.pingInterval = Duration(seconds: 10);
          return socket.transform(
            StreamTransformer.fromHandlers(
              handleData: (data, sink) {
                if (data is List<int>) {
                  // Inlined deserialization logic
                  final download = Download.create()..mergeFromProto3Json(jsonDecode(utf8.decode(data)));
                  sink.add(download);
                } else {
                  sink.addError('deserialization failed data: $data');
                }
              },
            ),
          );
        });
  }

  static Future<DownloadMetadataResponse> get(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx.get(Uri.https(httpx.host(), "/d/${id}", null), options: options).then((v) {
      return DownloadMetadataResponse.create()..mergeFromProto3Json(jsonDecode(v.body));
    });
  }

  static Future<DownloadUpdateResponse> update(
    String id,
    Download upd, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .put(
          Uri.https(httpx.host(), "/d/${id}"),
          body: jsonEncode(DownloadUpdateRequest(download: upd).toProto3Json()),
          options: [httpx.Accept.json, httpx.Content.json, ...options],
        )
        .then((v) {
          return Future.value(
            DownloadUpdateResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<DownloadPauseResponse> pause(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .delete(
          Uri.https(httpx.host(), "/d/${id}/pause", null),
          body: jsonEncode({}),
          options: options,
        )
        .then(httpx.auto_error)
        .then((v) {
          return Future.value(
            DownloadPauseResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<DownloadTuneResponse> tune(
    String id,
    DownloadTuneRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/d/${id}/tune", null),
          body: jsonEncode(req.toProto3Json()),
          options: options,
        )
        .then((v) {
          return Future.value(
            DownloadTuneResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<DownloadDeleteResponse> reset(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .delete(
          Uri.https(httpx.host(), "/d/${id}"),
          body: jsonEncode({}),
          options: options,
        )
        .then((v) {
          return Future.value(
            DownloadDeleteResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }
}
