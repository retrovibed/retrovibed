import 'package:http/http.dart' as http;
import 'package:http_parser/http_parser.dart';
import 'package:console/retrovibed.dart' as retro;

var _host = "localhost:9998";

String host() {
  return _host;
}

void set(String uri) {
  _host = uri;
}

// return an identity token for the local device.
String auto_bearer() {
  return "bearer ${retro.bearer_token()}";
}

// return a identity token from the currently connected host.
String auto_bearer_host() {
  return "bearer ${retro.bearer_token()}";
}

abstract class mimetypes {
  static MediaType parse(String s) {
    try {
      return MediaType.parse(s);
    } catch (e) {
      print(
        "failed to parse mimetype ${s} ${e} returning application/octet-stream",
      );
      return MediaType("application", "octet-stream");
    }
  }

  static MediaType maybe(String? s) {
    if (s == null) return MediaType("application", "octet-stream");
    return parse(s);
  }
}

Future<http.MultipartFile> uploadable(
  String path,
  String name,
  String mimetype, {
  String field = 'content',
}) {
  return http.MultipartFile.fromPath(
    field,
    path,
    filename: name,
    contentType: mimetypes.parse(mimetype),
  );
}
