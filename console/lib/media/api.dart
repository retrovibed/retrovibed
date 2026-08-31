import 'dart:convert';
import 'dart:async';
import 'package:flutter/foundation.dart';
import 'package:qs_dart/qs_dart.dart' as qs;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media/media.pb.dart';
import 'package:retrovibed/media/content.addressable.storage.pb.dart' as cas;
import 'package:http/http.dart' as http;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/timex.dart' as timex;

export 'package:retrovibed/media/media.pb.dart';
export 'search.state.dart';

typedef FnMediaSearch =
    Future<MediaSearchResponse> Function(
      MediaSearchRequest req, {
      String? host,
      List<httpx.Option> options,
    });

typedef FnMediaFind =
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

typedef FnMkdir =
    Future<MediaUploadResponse> Function(
      String description,
      String parent, {
      List<httpx.Option> options,
    });

typedef FnMediaDownload =
    Future<http.StreamedResponse> Function(
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

  static Future<MediaSearchResponse> emptysearch(
    MediaSearchRequest req, {
    String? host,
    List<httpx.Option> options = const [],
  }) async {
    return Future.value(
      MediaSearchResponse(items: [], next: req),
    );
  }

  static Future<MediaSearchResponse> search(
    MediaSearchRequest req, {
    String? host,
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(
            host ?? httpx.host(),
            "/m/",
          ).replace(query: qs.encode(req.toProto3Json())),
          options: [httpx.Content.urlencoded, httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(MediaSearchResponse.create(), jsonDecode(v.body)),
          );
        });
  }

  static FnMediaSearch searchendpoint(String host, List<httpx.Option> options) {
    final endpointHost = host;
    final endpointOptions = options;
    return (req, {String? host, List<httpx.Option> options = const []}) =>
        media.search(req, host: endpointHost, options: endpointOptions);
  }

  static FnMediaFind randomendpoint(String host, List<httpx.Option> options) {
    final endpointHost = host;
    final endpointOptions = options;
    return (req, {List<httpx.Option> options = const []}) =>
        media.random(req, host: endpointHost, options: endpointOptions);
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
            httpx.fromProto3JsonSafe(MediaUpdateResponse.create(), jsonDecode(v.body)),
          );
        });
  }

  static Future<MediaFindResponse> random(
    MediaSearchRequest req, {
    String? host,
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(
            host ?? httpx.host(),
            "/m/random",
          ).replace(query: qs.encode(req.toProto3Json())),
          options: options,
        )
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(MediaFindResponse.create(), jsonDecode(v.body)),
          );
        });
  }

  static Future<MediaFindResponse> similar(
    String mediaId,
    MediaSearchRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(httpx.host(), "/similar/$mediaId").replace(query: qs.encode(req.toProto3Json())),
          options: options,
        )
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(MediaFindResponse.create(), jsonDecode(v.body)),
          );
        });
  }

  // acoustic returns a random function anchored to seed: each call fetches one
  // acoustically similar track from /similar, excluding what's already been played,
  // and falls back to random once the similarity engine has nothing left (cold start,
  // below threshold, or library too small).
  // the seed is the last excluded id - by construction (see PlayQueue.recent)
  // that's always whatever's currently playing, so no separate seed parameter
  // is needed here.
  static Future<MediaFindResponse> acoustic(
    MediaSearchRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    // if excluded is empty, its a bug.
    return similar(req.excluded.last, req, options: options).catchError((cause) {
      print("acoustic similarity lookup failed, falling back to random: $cause");
      return random(req, options: options);
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

  static Future<MediaDeleteResponse> delete(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .delete(
          Uri.https(httpx.host(), "/m/${id}"),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(MediaDeleteResponse.create(), jsonDecode(v.body)),
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
            httpx.fromProto3JsonSafe(MediaUpdateResponse.create(), jsonDecode(v.body)),
          );
        });
  }

  static Future<cas.MediaDeleteResponse> unarchive(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx.delete(Uri.https(httpx.metaendpoint(), "/m/${id}"), options: options).then((v) {
      return Future.value(
        httpx.fromProto3JsonSafe(cas.MediaDeleteResponse.create(), jsonDecode(v.body)),
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

  // a folder shares the upload endpoint because it shares the destination and the
  // response shape. the directory mimetype is what tells the daemon there are no bytes
  // coming.
  static Future<MediaUploadResponse> mkdir(
    String description,
    String parent, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/m/"),
          body: {
            "mimetype": mimex.directory,
            "description": description,
            "parent_id": parent,
          },
          options: [httpx.Accept.json, httpx.Content.urlencoded, ...options],
        )
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(MediaUploadResponse.create(), jsonDecode(v.body)),
          );
        });
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
          httpx.fromProto3JsonSafe(MediaUploadResponse.create(), jsonDecode(s)),
        );
      });
    });
  }
}

abstract class download {
  static bool completed(Download d) => timex.iso8601(d.completedAt).isBefore(timex.inf);
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
            httpx.fromProto3JsonSafe(DownloadSearchResponse.create(), jsonDecode(v.body)),
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
            httpx.fromProto3JsonSafe(DownloadSearchResponse.create(), jsonDecode(v.body)),
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
            httpx.fromProto3JsonSafe(MagnetCreateResponse.create(), jsonDecode(v.body)),
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
          httpx.fromProto3JsonSafe(MediaUploadResponse.create(), jsonDecode(s)),
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
          return httpx.fromProto3JsonSafe(DownloadBeginResponse.create(), jsonDecode(v.body));
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
            httpx.fromProto3JsonSafe(MetadataSyncResponse.create(), jsonDecode(v.body)),
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
                  final download = httpx.fromProto3JsonSafe(Download.create(), jsonDecode(utf8.decode(data)));
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
      return httpx.fromProto3JsonSafe(DownloadMetadataResponse.create(), jsonDecode(v.body));
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
            httpx.fromProto3JsonSafe(DownloadUpdateResponse.create(), jsonDecode(v.body)),
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
            httpx.fromProto3JsonSafe(DownloadPauseResponse.create(), jsonDecode(v.body)),
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
            httpx.fromProto3JsonSafe(DownloadTuneResponse.create(), jsonDecode(v.body)),
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
            httpx.fromProto3JsonSafe(DownloadDeleteResponse.create(), jsonDecode(v.body)),
          );
        });
  }
}
