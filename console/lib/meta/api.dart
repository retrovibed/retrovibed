import 'dart:convert';
import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:http/http.dart' as http;
import 'package:console/httpx.dart' as httpx;
import './meta.daemon.pb.dart';
import './meta.daemon.pbenum.dart';
import './meta.daemon.pbjson.dart';

abstract class daemons {
  static final client = http.Client();
  static DaemonSearchRequest request({int limit = 0}) =>
      DaemonSearchRequest(limit: fixnum.Int64(limit));
  static DaemonSearchResponse response({DaemonSearchRequest? next}) =>
      DaemonSearchResponse(next: next ?? request(limit: 128), items: []);

  static Future<DaemonSearchResponse> get(DaemonSearchRequest req) async {
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
            DaemonSearchResponse.create()
              ..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }
}
