import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/ddisc/plugin/envfile.dart';

const _example = "FOO=\"bar\" # derp 0\n# derp 1\nBAR=\"baz\"\nBIZ=\"BAN\"\n# derp 2\n";

void main() {
  group('parseEnv', () {
    test('comment block and inline comment both attach as hints', () {
      expect(parseEnv(_example), equals([
        Variable('FOO', 'bar', 'derp 0'),
        Variable('BAR', 'baz', 'derp 1'),
        Variable('BIZ', 'BAN', ''),
      ]));
    });

    test('empty content yields no variables', () {
      expect(parseEnv(''), isEmpty);
    });

    test('unquoted value with no comment', () {
      expect(parseEnv('FOO=bar\n'), equals([Variable('FOO', 'bar', '')]));
    });

    test('blank line clears pending comment', () {
      expect(parseEnv('# orphaned\n\nFOO=bar\n'), equals([Variable('FOO', 'bar', '')]));
    });
  });

  group('serializeEnv', () {
    test('empty list yields empty content', () {
      expect(serializeEnv([]), equals(''));
    });

    test('serializes only the given variables, dropping anything else', () {
      final got = serializeEnv([Variable('FOO', 'bar', ''), Variable('B', '', '')]);
      expect(got, equals('FOO=bar\nB=\n'));
    });

    test('quotes values and preserves hints as inline comments', () {
      final got = serializeEnv([Variable('FOO', 'has space', 'derp')]);
      expect(got, equals('FOO="has space" # derp\n'));
    });

    test('round trip: parseEnv(serializeEnv(variables)) reflects exactly the given variables', () {
      final variables = [Variable('FOO', 'bar', 'derp 0'), Variable('BIZ', 'BAN', '')];
      expect(parseEnv(serializeEnv(variables)), equals(variables));
    });
  });
}
