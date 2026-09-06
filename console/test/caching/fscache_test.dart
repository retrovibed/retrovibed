import 'dart:convert';
import 'dart:io';

import 'package:crypto/crypto.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/caching/fscache.dart';

void main() {
  late Directory tmp;
  late Dir cache;

  setUp(() {
    tmp = Directory.systemTemp.createTempSync('fscache_test_');
    cache = newFSCache(tmp.path);
  });

  tearDown(() {
    if (tmp.existsSync()) tmp.deleteSync(recursive: true);
  });

  group('maybe', () {
    test('calls fn and returns value on cache miss', () {
      final v = cache.maybe<bool>('key', () => true);
      expect(v, isTrue);
    });

    test('returns cached value on hit without calling fn', () {
      cache.maybe<bool>('key', () => true);

      var calls = 0;
      final v = cache.maybe<bool>('key', () {
        calls++;
        return false;
      });

      expect(v, isTrue);
      expect(calls, 0);
    });

    test('fn is only called once across repeated misses for same key', () {
      var calls = 0;
      cache.maybe<int>('key', () => ++calls);
      cache.maybe<int>('key', () => ++calls);
      cache.maybe<int>('key', () => ++calls);

      expect(calls, 1);
    });

    test('different keys are independent', () {
      cache.maybe<int>('a', () => 1);
      cache.maybe<int>('b', () => 2);

      expect(cache.maybe<int>('a', () => 99), equals(1));
      expect(cache.maybe<int>('b', () => 99), equals(2));
    });

    test('persists across separate Dir instances rooted at same path', () {
      cache.maybe<String>('key', () => 'hello');

      final cache2 = newFSCache(tmp.path);
      var calls = 0;
      final v = cache2.maybe<String>('key', () {
        calls++;
        return 'world';
      });

      expect(v, equals('hello'));
      expect(calls, 0);
    });

    test('falls through to fn on corrupt entry', () {
      final file = File('${tmp.path}/${_md5('key')}');
      file.writeAsStringSync('not valid json{{{');

      final v = cache.maybe<bool>('key', () => true);
      expect(v, isTrue);
    });
  });

  group('write', () {
    test('returns the written value', () {
      expect(cache.write<bool>('key', false), isFalse);
    });

    test('subsequent maybe returns written value', () {
      cache.write<int>('key', 42);

      final v = cache.maybe<int>('key', () => 0);
      expect(v, equals(42));
    });

    test('overwrites an existing cached entry', () {
      cache.maybe<int>('key', () => 1);
      cache.write<int>('key', 99);

      expect(cache.maybe<int>('key', () => 0), equals(99));
    });
  });

  group('pathFor', () {
    test('resolves to md5 of the key under the cache dir', () {
      expect(cache.pathFor('key').path, equals('${tmp.path}/${_md5('key')}'));
    });

    test('agrees with the file write actually creates', () {
      cache.write<int>('key', 42);

      expect(cache.pathFor('key').existsSync(), isTrue);
      expect(
        cache.pathFor('key').readAsBytesSync(),
        equals(File('${tmp.path}/${_md5('key')}').readAsBytesSync()),
      );
    });

    test('returns a path for a key that was never written', () {
      expect(cache.pathFor('absent').existsSync(), isFalse);
    });

    test('distinct keys resolve to distinct files', () {
      expect(cache.pathFor('a').path, isNot(equals(cache.pathFor('b').path)));
    });
  });

  group('clear', () {
    test('causes next maybe to call fn again', () {
      cache.maybe<bool>('key', () => true);
      cache.clear();

      var calls = 0;
      cache.maybe<bool>('key', () {
        calls++;
        return false;
      });

      expect(calls, equals(1));
    });

    test('removes all keys', () {
      cache.maybe<int>('a', () => 1);
      cache.maybe<int>('b', () => 2);
      cache.clear();

      var aCalls = 0, bCalls = 0;
      cache.maybe<int>('a', () => ++aCalls);
      cache.maybe<int>('b', () => ++bCalls);

      expect(aCalls, equals(1));
      expect(bCalls, equals(1));
    });

    test('cache is usable after clear', () {
      cache.clear();
      final v = cache.maybe<String>('key', () => 'after clear');
      expect(v, equals('after clear'));
    });
  });
}

// mirror of Dir._keyfile for test assertions
String _md5(String key) => md5.convert(utf8.encode(key)).toString();
