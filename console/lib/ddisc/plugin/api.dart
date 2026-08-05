import 'dart:convert';
import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:http/http.dart' as http;
import 'package:retrovibed/httpx.dart' as httpx;
import './searchplugin.pb.dart';

export './searchplugin.pb.dart';

abstract class plugins {
  static PluginSearchRequest request({int limit = 0}) => PluginSearchRequest(limit: fixnum.Int64(limit));
  static PluginSearchResponse response({PluginSearchRequest? next}) =>
      PluginSearchResponse(next: next ?? request(limit: 100), items: []);

  static Future<PluginSearchResponse> search(
    PluginSearchRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(httpx.host(), "/ddisc/plugin/", httpx.params(req.toProto3Json())),
          options: options,
        )
        .then((v) {
          return Future.value(
            PluginSearchResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<PluginFindResponse> find(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx.get(Uri.https(httpx.host(), "/ddisc/plugin/${id}"), options: options).then((v) {
      return Future.value(
        PluginFindResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
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

  static Future<PluginCreateResponse> upload(
    String name,
    http.MultipartRequest Function(http.MultipartRequest req) mkreq, {
    List<httpx.Option> options = const [],
  }) async {
    final client = http.Client();
    final r0 = mkreq(
      http.MultipartRequest("POST", Uri.https(httpx.host(), "/ddisc/plugin/")),
    );
    r0.fields["name"] = name;

    return httpx.request(options).then((r) {
      r0.headers.addAll(r.headers);
      return client.send(r0).then(httpx.auto_error).then((v) {
        return v.stream.bytesToString().then((s) {
          return Future.value(
            PluginCreateResponse.create()..mergeFromProto3Json(jsonDecode(s)),
          );
        });
      });
    });
  }

  static Future<PluginDeleteResponse> delete(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx.delete(Uri.https(httpx.host(), "/ddisc/plugin/${id}"), options: options).then((v) {
      return Future.value(
        PluginDeleteResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
      );
    });
  }
}

abstract class environment {
  static Future<String> get(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx.get(Uri.https(httpx.host(), "/ddisc/plugin/environment/${id}"), options: options).then((v) => v.body);
  }

  static Future<String> update(
    String id,
    String content, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/ddisc/plugin/environment/${id}"),
          options: options,
          body: content,
        )
        .then((v) => v.body);
  }

  static Future<void> delete(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx.delete(Uri.https(httpx.host(), "/ddisc/plugin/environment/${id}"), options: options).then((_) {});
  }
}
