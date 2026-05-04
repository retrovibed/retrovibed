import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/env.dart';

void main() {
  group('lookup', () {
    final env = <String, String?>{
      'FIRST': 'one',
      'SECOND': 'two',
      'EMPTY': '',
    };

    test('returns value for the first matching key', () {
      expect(lookup(['SECOND', 'FIRST'], env, 'x'), 'two');
    });

    test('returns fallback when no key exists', () {
      expect(lookup(['MISSING', 'ALSO_MISSING'], env, 'default'), 'default');
    });

    test('returns empty string when key maps to empty', () {
      expect(lookup(['EMPTY'], env, 'x'), '');
    });

    test('handles single-element list', () {
      expect(lookup(['FIRST'], env, 'x'), 'one');
    });

    test('handles empty key list', () {
      expect(lookup(<String>[], env, 'fallback'), 'fallback');
    });
  });

  group('string', () {
    test('returns empty string for unset env var', () {
      expect(string(['RETROVIBED_NONEXISTENT']), '');
    });

    test('returns fallback when no key exists', () {
      expect(string(['RETROVIBED_NONEXISTENT'], fallback: 'prod'), 'prod');
    });
  });
}
