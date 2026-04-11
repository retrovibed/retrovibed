import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/lucene.dart' as lucene;
import 'package:retrovibed/timex.dart' as timex;

void main() {
  group('Timestamp', () {
    final field = lucene.DateRange.auto(
      'date',
      timex.Range(timex.epoch, timex.inf),
      (_) {},
    );

    group('suggestions', () {
      test('empty partial returns all three presets', () {
        final s = field.suggestions('');
        expect(s.length, 3);
      });

      test('parse round-trips a suggestion completion into its range', () {
        final s = field.suggestions('last 7').first;
        final afterColon = s.completion.substring(s.completion.indexOf(':') + 1);
        final parsed = field.parse(afterColon);
        final age = parsed.max.difference(parsed.min);
        expect(age, equals(const Duration(days: 7)));
      });

      test('label partial filters to matching presets', () {
        final s = field.suggestions('last 3');
        expect(s.length, 1);
        expect(s.first.label, 'last 30 days');
      });

      test('no match returns empty', () {
        expect(field.suggestions('xyz'), isEmpty);
      });

      test('completion is a lucene range with field name prefix', () {
        final s = field.suggestions('last 7');
        expect(s.first.completion, startsWith('date:['));
        expect(s.first.completion, contains(' TO '));
        expect(s.first.completion, endsWith(']'));
      });
    });

    group('from', () {
      test('finite range serialises both bounds as ISO 8601', () {
        final node = field.from((
          min: DateTime.utc(2025, 1, 1),
          max: DateTime.utc(2026, 1, 1),
        )) as lucene.Range;
        expect(node.min, '2025-01-01T00:00:00.000Z');
        expect(node.max, '2026-01-01T00:00:00.000Z');
        expect(node.field, 'date');
      });

      test('open-ended range (inf end) serialises max as null', () {
        final node = field.from((
          min: DateTime.utc(2025, 1, 1),
          max: timex.inf,
        )) as lucene.Range;
        expect(node.min, '2025-01-01T00:00:00.000Z');
        expect(node.max, isNull);
      });
    });

    group('of', () {
      test('of round-trips the value into field.current', () {
        final value = (min: DateTime.utc(2025, 1, 1), max: timex.inf);
        final result = field.of(value);
        expect((result as dynamic).field.current, value);
      });

      test('of preserves the field name', () {
        final value = (min: DateTime.utc(2025, 1, 1), max: timex.inf);
        final result = field.of(value);
        expect((result as dynamic).field.name, 'date');
      });

      test('of preserves the defaultValue', () {
        final defaultRange = timex.Range(timex.epoch, timex.inf);
        final f = lucene.DateRange.auto('date', defaultRange, (_) {});
        final result = f.of((min: DateTime.utc(2025, 1, 1), max: timex.inf));
        expect(
          (result as dynamic).field.defaultValue,
          (min: defaultRange.begin, max: defaultRange.end),
        );
      });
    });

    group('apply', () {
      test('range node calls setter with parsed min and max', () {
        ({DateTime min, DateTime max})? result;
        final f = lucene.DateRange.auto(
          'date',
          timex.Range(timex.epoch, timex.inf),
          (r) => result = (min: r.begin, max: r.end),
        );
        f.apply(lucene.Range('date', min: '2025-01-01T00:00:00Z', max: '2026-01-01T00:00:00Z'));
        expect(result?.min, DateTime.utc(2025, 1, 1));
        expect(result?.max, DateTime.utc(2026, 1, 1));
      });

      test('unrelated field does not call setter', () {
        bool called = false;
        final f = lucene.DateRange.auto(
          'date',
          timex.Range(timex.epoch, timex.inf),
          (_) => called = true,
        );
        f.apply(lucene.Range('other', min: '2025-01-01T00:00:00Z'));
        expect(called, isFalse);
      });
    });
  });
}
