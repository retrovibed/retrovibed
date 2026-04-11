import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/lucene.dart' as lucene;

void main() {
  group('Elapsed', () {
    final field = lucene.Elapsed.auto('runtime', Duration.zero, (_) {});

    group('parse', () {
      test(
        'minutes',
        () => expect(field.parse('PT90M'), Duration(minutes: 90)),
      );
      test('hours', () => expect(field.parse('PT2H'), Duration(hours: 2)));
      test(
        'hours and minutes',
        () => expect(field.parse('PT1H30M'), Duration(hours: 1, minutes: 30)),
      );
      test('empty returns defaultValue', () => expect(field.parse(''), Duration.zero));
      test('invalid returns defaultValue', () => expect(field.parse('abc'), Duration.zero));
    });

    group('suggestions', () {
      test('empty partial returns all presets', () {
        final s = field.suggestions('');
        expect(s, hasLength(4));
        expect(s.map((s) => s.label), containsAll(['30 min', '60 min', '90 min', '2 hours']));
      });

      test('partial filters by label', () {
        final s = field.suggestions('2');
        expect(s, hasLength(1));
        expect(s.first.label, '2 hours');
      });

      test('no match returns empty', () {
        expect(field.suggestions('xyz'), isEmpty);
      });

      test('completions are prefixed with field name', () {
        final s = field.suggestions('');
        expect(s.every((s) => s.completion.startsWith('runtime:')), isTrue);
      });

      test('completion round-trips through parse', () {
        for (final s in field.suggestions('')) {
          final raw = s.completion.substring(s.completion.indexOf(':') + 1);
          expect(field.parse(raw), isNot(Duration.zero),
              reason: '${s.label} completion "$raw" did not parse to a non-zero Duration');
        }
      });
    });
  });
}
