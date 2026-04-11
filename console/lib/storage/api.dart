import 'dart:convert';
import 'package:retrovibed/httpx.dart' as httpx;
import './storage.pb.dart';
export './storage.pb.dart';

abstract class api {
    static Future<StorageSettingsResponse> get({
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(Uri.https(httpx.host(), "/s/storage/", {}), options: options)
        .then((v) {
          return Future.value(
            StorageSettingsResponse.create()
              ..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<StorageSettingsResponse> create(
    StorageSettingsResponse req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/s/storage/", {}),
          options: options,
          body: jsonEncode(req.toProto3Json()),
        ).then((v) {
          return Future.value(
            StorageSettingsResponse.create()
              ..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<StorageSettingsResponse> delete({
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .delete(
          Uri.https(httpx.host(), "/s/storage/", {}),
          options: options,
        ).then((v) {
          return Future.value(
            StorageSettingsResponse.create()
              ..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }
}