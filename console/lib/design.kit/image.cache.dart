import 'dart:io';
import 'package:path/path.dart' as path;
import 'package:retrovibed/caching.dart' as caching;

final caching.Dir _images = caching.newFSCache(
  path.join(caching.global().cache, 'images'),
  codec: const caching.bytescodec(),
);

// Returns the disk-cached file for url. If it isn't cached yet, invokes
// fetch to populate it first.
Future<File> cached(String url, Future<List<int>> Function() fetch) async {
  final file = _images.pathFor(url);
  if (file.existsSync()) return file;

  return fetch().then((bytes) {
    _images.write(url, bytes);
    return file;
  });
}
