import 'package:flutter/material.dart';
import 'package:mime/mime.dart' as mimetype;
export 'package:mime/mime.dart';

final resolver =
    mimetype.MimeTypeResolver()
      ..addMagicNumber([0x4F, 0x67, 0x67, 0x53], "video/ogg");

String fromFile(String s, {List<int>? magicbits}) {
  return maybe(resolver.lookup(s, headerBytes: magicbits));
}

String maybe(String? s) {
  return s ?? "application/octet-stream";
}

IconData icon(String mimetype) {
  if (mimetype.startsWith('video/')) {
    return Icons.movie;
  }

  if (mimetype.startsWith('audio/')) {
    return Icons.music_note_outlined;
  }

  return Icons.file_open_outlined;
}
