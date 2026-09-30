import 'dart:convert';
import 'package:qs_dart/qs_dart.dart' as qs;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media/media.autoimport.pb.dart';

export 'package:retrovibed/media/media.autoimport.pb.dart';

// monitored directories are paths on this machine, every request goes to the local daemon
// (httpx.localhost()), never to httpx.host() which follows the currently selected library.

typedef FnAutoimportSearch =
    Future<AutoimportDirectorySearchResponse> Function(
      AutoimportDirectorySearchRequest req, {
      List<httpx.Option> options,
    });

typedef FnAutoimportCreate =
    Future<AutoimportDirectoryCreateResponse> Function(
      AutoimportDirectoryCreateRequest req, {
      List<httpx.Option> options,
    });

typedef FnAutoimportUpdate =
    Future<AutoimportDirectoryUpdateResponse> Function(
      String id,
      AutoimportDirectoryUpdateRequest req, {
      List<httpx.Option> options,
    });

typedef FnAutoimportDelete =
    Future<AutoimportDirectoryDeleteResponse> Function(
      String id, {
      List<httpx.Option> options,
    });

const int autoimportModeCopy = 0;
const int autoimportModeMove = 1;

String autoimportModeLabel(int mode) {
  switch (mode) {
    case autoimportModeCopy:
      return "copy";
    case autoimportModeMove:
      return "move";
    default:
      return "unknown";
  }
}

abstract class autoimport {
  static AutoimportDirectorySearchRequest request({String query = "", int limit = 32}) {
    return AutoimportDirectorySearchRequest(
      query: query,
      offset: ds.Int64(0),
      limit: ds.Int64(limit),
    );
  }

  static AutoimportDirectorySearchResponse response({AutoimportDirectorySearchRequest? next}) {
    return AutoimportDirectorySearchResponse(next: next ?? request(), items: []);
  }

  // a new monitored directory with the daemon's defaults: an hour debounce, copying files.
  static AutoimportDirectory directory() {
    return AutoimportDirectory(debounce: ds.Int64(3600), mode: autoimportModeCopy);
  }

  static Future<AutoimportDirectorySearchResponse> search(
    AutoimportDirectorySearchRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(httpx.localhost(), "/autoimport/", qs.decode(qs.encode(req.toProto3Json()))),
          options: [httpx.Content.urlencoded, httpx.Accept.json, ...options],
        )
        .then((v) => httpx.fromProto3JsonSafe(AutoimportDirectorySearchResponse.create(), jsonDecode(v.body)));
  }

  static Future<AutoimportDirectoryCreateResponse> create(
    AutoimportDirectoryCreateRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.localhost(), "/autoimport/"),
          body: jsonEncode(req.toProto3Json()),
          options: [httpx.Content.json, httpx.Accept.json, ...options],
        )
        .then((v) => httpx.fromProto3JsonSafe(AutoimportDirectoryCreateResponse.create(), jsonDecode(v.body)));
  }

  static Future<AutoimportDirectoryUpdateResponse> update(
    String id,
    AutoimportDirectoryUpdateRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.localhost(), "/autoimport/${id}"),
          body: jsonEncode(req.toProto3Json()),
          options: [httpx.Content.json, httpx.Accept.json, ...options],
        )
        .then((v) => httpx.fromProto3JsonSafe(AutoimportDirectoryUpdateResponse.create(), jsonDecode(v.body)));
  }

  static Future<AutoimportDirectoryDeleteResponse> delete(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .delete(
          Uri.https(httpx.localhost(), "/autoimport/${id}"),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) => httpx.fromProto3JsonSafe(AutoimportDirectoryDeleteResponse.create(), jsonDecode(v.body)));
  }
}
