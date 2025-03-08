import 'package:mime/mime.dart';
export 'package:mime/mime.dart';

final resolver =
    MimeTypeResolver()..addMagicNumber([0x4F, 0x67, 0x67, 0x53], "video/ogg");

String fromFile(String s, {List<int>? magicbits}) {
  return maybe(resolver.lookup(s, headerBytes: magicbits));
}

String maybe(String? s) {
  return s ?? "application/octet-stream";
}
