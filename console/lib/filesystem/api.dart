import 'dart:convert';
import 'package:qs_dart/qs_dart.dart' as qs;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media/media.filesystem.pb.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/uuidx.dart' as uuidx;

export 'package:retrovibed/media/media.filesystem.pb.dart';

typedef FnFilesystemSearch =
    Future<FilesystemSearchResponse> Function(
      FilesystemSearchRequest req, {
      String? host,
      List<httpx.Option> options,
    });

typedef FnFilesystemCreate =
    Future<FilesystemCreateResponse> Function(
      FilesystemCreateRequest req, {
      List<httpx.Option> options,
    });

typedef FnFilesystemMove =
    Future<FilesystemMoveResponse> Function(
      String id,
      FilesystemMoveRequest req, {
      List<httpx.Option> options,
    });

typedef FnFilesystemDelete =
    Future<FilesystemDeleteResponse> Function(
      String id, {
      List<httpx.Option> options,
    });

abstract class filesystem {
  static FilesystemSearchRequest request({
    int limit = 0,
    String query = "",
    String directory = "",
    List<String> mimetypes = const [],
  }) => FilesystemSearchRequest(
    limit: ds.Int64(limit),
    query: query,
    mimetypes: mimetypes,
    directoryId: directory.isEmpty ? uuidx.min() : directory,
  );

  static FilesystemSearchResponse response({FilesystemSearchRequest? next}) =>
      FilesystemSearchResponse(next: next ?? request(limit: 100), items: [], breadcrumb: []);

  static Future<FilesystemSearchResponse> emptysearch(
    FilesystemSearchRequest req, {
    String? host,
    List<httpx.Option> options = const [],
  }) async {
    return Future.value(FilesystemSearchResponse(items: [], breadcrumb: [], next: req));
  }

  static Future<FilesystemSearchResponse> search(
    FilesystemSearchRequest req, {
    String? host,
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(host ?? httpx.host(), "/fs/").replace(query: qs.encode(req.toProto3Json())),
          options: [httpx.Content.urlencoded, httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(FilesystemSearchResponse.create(), jsonDecode(v.body)),
          );
        });
  }

  static Future<FilesystemCreateResponse> create(
    FilesystemCreateRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/fs/"),
          body: jsonEncode(req.toProto3Json()),
          options: [httpx.Accept.json, httpx.Content.json, ...options],
        )
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(FilesystemCreateResponse.create(), jsonDecode(v.body)),
          );
        });
  }

  static Future<FilesystemMoveResponse> move(
    String id,
    FilesystemMoveRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/fs/${id}"),
          body: jsonEncode(req.toProto3Json()),
          options: [httpx.Accept.json, httpx.Content.json, ...options],
        )
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(FilesystemMoveResponse.create(), jsonDecode(v.body)),
          );
        });
  }

  static Future<FilesystemDeleteResponse> delete(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .delete(
          Uri.https(httpx.host(), "/fs/${id}"),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(FilesystemDeleteResponse.create(), jsonDecode(v.body)),
          );
        });
  }
}
