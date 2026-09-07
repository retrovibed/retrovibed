import 'dart:io';
import 'dart:convert';
import 'package:crypto/crypto.dart';
import 'package:path/path.dart' as path;

import 'caching.dart';
import 'dirs.dart';

/// Filesystem-backed cache with no TTL. Entries persist until [clear] is called.
///
/// Usage:
///   final c = FSCache(someDirectory);
///   final v = c.maybe<bool>('my-key', () => false);
///   c.clear();
class Dir {
  static Dir cacheroot = Dir(Directory(global().cache));
  final Directory dir;
  final Codec codec;

  Dir(this.dir, {this.codec = const jsoncodec()});

  /// Returns the cached value for [key] if present, otherwise calls [fn],
  /// persists the result to disk, and returns it.
  T maybe<T>(String key, T Function() fn) {
    final file = pathFor(key);

    if (file.existsSync()) {
      try {
        return codec.decode<T>(file.readAsBytesSync());
      } catch (_) {
        // corrupt entry — fall through and regenerate
      }
    }

    final v = fn();

    try {
      dir.createSync(recursive: true);
      file.writeAsBytesSync(codec.encode(v));
    } catch (_) {
      // best-effort write; still return the value
    }

    return v;
  }

  /// Writes [v] to the cache under [key], overwriting any existing entry.
  T write<T>(String key, T v) {
    final file = pathFor(key);
    dir.createSync(recursive: true);
    file.writeAsBytesSync(codec.encode(v));
    return v;
  }

  /// Returns the file backing [key], whether or not it has been written yet.
  /// Callers that need the raw file instead of a decoded value should derive
  /// it from here; reimplementing the key derivation lets the read and write
  /// paths drift onto different files.
  File pathFor(String key) => File(path.join(dir.path, _keyfile(key)));

  /// Removes all cached entries by deleting and recreating the cache directory.
  void clear() {
    if (dir.existsSync()) dir.deleteSync(recursive: true);
    dir.createSync(recursive: true);
  }

  String _keyfile(String key) => md5.convert(utf8.encode(key)).toString();
}

/// Creates an [Dir] rooted at [dirPath].
Dir newFSCache(String dirPath, {Codec codec = const jsoncodec()}) {
  return Dir(Directory(dirPath), codec: codec);
}
