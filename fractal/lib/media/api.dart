import 'package:fractal/media/media.pb.dart';
import 'package:http/http.dart' as http;
import 'package:fractal/httpx.dart' as httpx;
import 'dart:convert';

Future<MediaSearchResponse> recent() async {
  final client = http.Client();
  return client.get(Uri.https(httpx.host(), "/m/recent")).then((v) {
    return Future.value(
      MediaSearchResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
    );
  });
}

Future<MediaSearchResponse> discovered() async {
  final client = http.Client();
  return client.get(Uri.https(httpx.host(), "/m/discovered")).then((v) {
    return Future.value(
      MediaSearchResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
    );
  });
}
