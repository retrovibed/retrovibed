import 'package:fractal/media/media.pb.dart';
import 'package:http/http.dart' as http;
import 'package:fractal/httpx.dart' as httpx;

Future<MediaResponse> recent() async {
  final client = http.Client();
  return client
      .get(Uri.https(httpx.host(), "/m/recent"))
      .then((v) => Future.value(MediaResponse.fromJson(v.body)));
}

Future<MediaResponse> discovered() async {
  final client = http.Client();
  return client
      .get(Uri.https(httpx.host(), "/m/discovered"))
      .then((v) => Future.value(MediaResponse.fromJson(v.body)));
}
