import 'dart:convert';
import 'dart:io';

import 'package:crypto/crypto.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:path/path.dart' as path;
import 'package:retrovibed/caching.dart' as caching;
import 'package:retrovibed/design.kit/image.cache.dart' as imagecache;

void main() {
  late Directory tmp;
  late Directory images;

  // the cache Dir inside image.cache.dart is a top-level final, so it resolves
  // caching.global() lazily on the first cached() call and pins that path for
  // the rest of the process. the global therefore has to be set once up front,
  // with every test sharing the directory rather than getting its own.
  setUpAll(() {
    tmp = Directory.systemTemp.createTempSync('image_cache_test_');
    images = Directory(path.join(tmp.path, 'images'));
    caching.setglobal(caching.DirsWellKnown(cache: tmp.path));
  });

  tearDownAll(() {
    if (tmp.existsSync()) tmp.deleteSync(recursive: true);
  });

  setUp(() {
    if (images.existsSync()) images.deleteSync(recursive: true);
    images.createSync(recursive: true);
  });

  group('cached', () {
    test('returned file exists on disk after a fetch', () async {
      final file = await imagecache.cached(
        'https://example.com/a.png',
        () async => [1, 2, 3],
      );

      expect(file.existsSync(), isTrue);
    });

    test('returned file holds the fetched bytes', () async {
      final file = await imagecache.cached(
        'https://example.com/b.png',
        () async => [4, 5, 6, 7],
      );

      expect(file.readAsBytesSync(), equals([4, 5, 6, 7]));
    });

    test('writes under md5 of the url, not md5 of md5', () async {
      const url = 'https://example.com/c.png';
      await imagecache.cached(url, () async => [8]);

      final key = md5.convert(utf8.encode(url)).toString();
      expect(
        images.listSync().map((e) => path.basename(e.path)),
        equals([key]),
      );
    });

    test('second call is a hit and does not refetch', () async {
      const url = 'https://example.com/d.png';
      await imagecache.cached(url, () async => [9]);

      var calls = 0;
      await imagecache.cached(url, () async {
        calls++;
        return [10];
      });

      expect(calls, equals(0));
    });

    test('second call returns the same path with the original bytes', () async {
      const url = 'https://example.com/e.png';
      final first = await imagecache.cached(url, () async => [11, 12]);
      final second = await imagecache.cached(url, () async => [13, 14]);

      expect(second.path, equals(first.path));
      expect(second.readAsBytesSync(), equals([11, 12]));
    });

    test('a file written by a prior run is served without fetching', () async {
      const url = 'https://example.com/f.png';
      File(path.join(images.path, md5.convert(utf8.encode(url)).toString()))
          .writeAsBytesSync([15, 16]);

      var calls = 0;
      final file = await imagecache.cached(url, () async {
        calls++;
        return [17];
      });

      expect(calls, equals(0));
      expect(file.readAsBytesSync(), equals([15, 16]));
    });

    test('distinct urls get distinct files', () async {
      final a = await imagecache.cached(
        'https://example.com/g.png',
        () async => [18],
      );
      final b = await imagecache.cached(
        'https://example.com/h.png',
        () async => [19],
      );

      expect(a.path, isNot(equals(b.path)));
      expect(a.readAsBytesSync(), equals([18]));
      expect(b.readAsBytesSync(), equals([19]));
    });

    test('a failed fetch propagates and caches nothing', () async {
      const url = 'https://example.com/i.png';

      await expectLater(
        imagecache.cached(url, () async => throw const SocketException('boom')),
        throwsA(isA<SocketException>()),
      );

      expect(
        File(path.join(images.path, md5.convert(utf8.encode(url)).toString()))
            .existsSync(),
        isFalse,
      );
    });
  });
}
