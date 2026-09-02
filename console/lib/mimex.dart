import 'package:flutter/material.dart';
import 'package:mime/mime.dart' as mimetype;
export 'package:mime/mime.dart';

const audio = "audio";
const video = "video";
const image = "image";
const metadataarchive = "application/vnd";
const mediaarchive = "application/vnd.retrovibed.media.archive";
const neural = "application/vnd.retrovibed.neural";
const search = "application/vnd.retrovibed.discovery.search.module";
const bittorrent = "application/x-bittorrent";
const binary = "application/octet-stream";
const directory = "inode/directory";
const pdf = "application/pdf";

final resolver = mimetype.MimeTypeResolver()..addMagicNumber([0x4F, 0x67, 0x67, 0x53], "video/ogg");

String fromFile(String s, {List<int>? magicbits}) {
  return maybe(resolver.lookup(s, headerBytes: magicbits));
}

String maybe(String? s) {
  s = s ?? "";
  return s.isNotEmpty ? s : binary;
}

const icomovie = Icons.movie;
const icoaudio = Icons.music_note_outlined;
const icoimage = Icons.image;
const icobinary = Icons.file_open_outlined;
const icometadataarchive = Icons.live_tv;
const iconneural = Icons.psychology;
const icofolder = Icons.folder_outlined;

String ext(String mime) {
  return mimetype.extensionFromMime(mime) ?? ".bin";
}

const List<String> folders = [directory];

const List<String> videos = [
  "video/mp4",
  "video/webm",
  "video/ogg",
  "video/mpeg",
  "video/dvd",
  "video/x-msvideo",
  "video/x-ms-wmv",
  "video/3gpp",
  "video/3gpp2",
  "video/quicktime",
  "video/mp2t",
  "video/x-flv",
  "video/x-matroska",
  "video/x-dv",
  "video/fli",
  "video/x-fli",
  "video/gl",
  "video/x-gl",
  "video/x-ms-asf",
  "video/avi",
  "video/msvideo",
  "video/avs-video",
  "video/dl",
  "video/x-dl",
  "video/animaflex",
];

const List<String> images = [
  "image/jpeg",
  "image/png",
  "image/gif",
  "image/webp",
  "image/svg+xml",
  "image/bmp",
  "image/tiff",
  "image/ico",
  "image/x-icon",
  "image/avif",
  "image/heic",
  "image/heif",
];

const List<String> audios = [
  "audio/x-wav",
  "audio/x-smpte336m",
  "audio/x-pn-wav",
  "audio/x-mpegurl",
  "audio/x-midi",
  "audio/x-m4a",
  "audio/x-gsm",
  "audio/x-aiff",
  "audio/weba",
  "audio/wav",
  "audio/vnd.wave",
  "audio/vnd.rn-realaudio",
  "audio/vnd.qcelp",
  "audio/vnd.nuera.ecelp9600",
  "audio/vnd.nuera.ecelp7470",
  "audio/vnd.nuera.ecelp4800",
  "audio/vnd.dts",
  "audio/vnd.dts.hd",
  "audio/vnd.dolby.mlp",
  "audio/opus",
  "audio/ogg",
  "audio/mpeg",
  "audio/mp4",
  "audio/mp3",
  "audio/midi",
  "audio/ilbc",
  "audio/g729",
  "audio/flac",
  "audio/basic",
  "audio/amr",
  "audio/aiff",
  "audio/adpcm",
  "audio/aac",
  "audio/3gpp2",
  "audio/3gpp",
];

List<String> of(IconData v) {
  if (v == icomovie) return videos;
  if (v == icoimage) return images;
  if (v == icoaudio) return audios;
  return const [];
}

int checksumfor(IconData v) {
  return checksum(of(v));
}

int checksum(List<String> mimes) {
  if (mimes.isEmpty) return -1;
  return Object.hashAllUnordered(mimes);
}

bool isVideo(String mimetype) => mimetype.startsWith('video');
bool isAudio(String mimetype) => mimetype.startsWith('audio');
bool isImage(String mimetype) => mimetype.startsWith('image');

// text the reader can render directly. the structured formats below are text on the wire
// even though their mimetype does not say so.
bool isText(String mimetype) =>
    mimetype.startsWith('text') ||
    const [
      "application/json",
      "application/xml",
      "application/javascript",
      "application/x-yaml",
      "application/yaml",
      "application/toml",
    ].contains(mimetype);

String category(List<String> mimes) {
  final sum = checksum(mimes);
  return switch (sum) {
    _ when sum == checksumfor(icomovie) => "video",
    _ when sum == checksumfor(icoaudio) => "audio",
    _ when sum == checksumfor(icoimage) => "image",
    _ => "",
  };
}

class CategoryOptionsLabel extends StatelessWidget {
  final List<String> mimetypes;
  const CategoryOptionsLabel(this.mimetypes, {super.key});

  static String text(List<String> mimetypes) {
    final sum = checksum(mimetypes);
    return switch (sum) {
      _ when sum == checksumfor(icomovie) => "Movie",
      _ when sum == checksumfor(icoaudio) => "Music",
      _ => "File",
    };
  }

  @override
  Widget build(BuildContext context) {
    final category = text(mimetypes);
    return Text("${category} options", key: ValueKey(category));
  }
}

IconData icon(String mimetype) {
  if (mimetype == directory) {
    return icofolder;
  }

  if (isVideo(mimetype)) {
    return icomovie;
  }

  if (isAudio(mimetype)) {
    return icoaudio;
  }

  if (isImage(mimetype)) {
    return icoimage;
  }

  if (mimetype == metadataarchive) {
    return icometadataarchive;
  }

  return icobinary;
}
