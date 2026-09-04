import 'dart:io';
import 'package:path/path.dart' as path;
import 'package:retrovibed/caching.dart' as caching;
import 'package:retrovibed/uuidx.dart' as uuidx;

final caching.Dir _images = caching.newFSCache(
  path.join(caching.global().cache, 'images'),
  codec: const caching.bytescodec(),
);

File _pathFor(String k) => File(path.join(_images.dir.path, k));

// Returns the disk-cached file for url. If it isn't cached yet, invokes
// fetch to populate it first.
Future<File> cached(String url, Future<List<int>> Function() fetch) async {
  final key = uuidx.md5x(url);
  final file = _pathFor(key);
  if (file.existsSync()) return file;

  final bytes = await fetch();
  _images.write(key, bytes);
  return file;
}
