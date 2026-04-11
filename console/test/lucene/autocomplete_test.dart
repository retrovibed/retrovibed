import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/lucene.dart' as lucene;
import 'package:retrovibed/timex.dart' as timex;

void main() {
  late List<lucene.Field> fields;

  setUp(() {
    fields = [
      lucene.Boolean.auto('hd', false, (_) {}),
      lucene.DateRange.auto(
        'date',
        timex.Range(timex.epoch, timex.inf),
        (_) {},
      ),
      lucene.Number.auto('peers', 0, (_) {}),
    ];
  });

  group('complete — field name', () {
    test('empty partial returns all field names', () {
      final s = lucene.complete('', fields);
      expect(s.map((s) => s.label), containsAll(['hd', 'date', 'peers']));
    });

    test('partial matches prefix', () {
      final s = lucene.complete('da', fields);
      expect(s, hasLength(1));
      expect(s.first.label, 'date');
    });

    test('no match returns empty', () {
      final s = lucene.complete('xyz', fields);
      expect(s, isEmpty);
    });

    test('completion includes trailing colon', () {
      final s = lucene.complete('hd', fields);
      expect(s.first.completion, 'hd:');
    });
  });

  group('complete — field values', () {
    test('after colon delegates to field suggestions', () {
      final s = lucene.complete('hd:', fields);
      expect(s.map((s) => s.label), containsAll(['on', 'off']));
    });

    test('after colon with partial filters options', () {
      final s = lucene.complete('hd:of', fields);
      expect(s, hasLength(1));
      expect(s.first.label, 'off');
    });

    test('unknown field after colon returns empty', () {
      final s = lucene.complete('unknown:', fields);
      expect(s, isEmpty);
    });

    test('boolean field suggests on/off', () {
      final s = lucene.complete('hd:', fields);
      expect(s.map((s) => s.label), containsAll(['on', 'off']));
    });

    test('timestamp field suggests presets', () {
      final s = lucene.complete('date:', fields);
      expect(s.length, greaterThan(0));
      expect(s.every((s) => s.completion.startsWith('date:')), isTrue);
    });

    test('timestamp partial filters presets', () {
      final s = lucene.complete('date:7', fields);
      expect(
        s.every((s) => s.label.contains('7') || s.completion.contains('7')),
        isTrue,
      );
    });
  });

  group('Boolean suggestions', () {
    late lucene.Boolean field;

    setUp(() {
      field = lucene.Boolean.auto('hd', false, (_) {});
    });

    test('empty partial returns on and off', () {
      final s = field.suggestions('');
      expect(s.map((s) => s.label), containsAll(['on', 'off']));
    });

    test('partial "of" returns only off', () {
      final s = field.suggestions('of');
      expect(s, hasLength(1));
      expect(s.first.label, 'off');
    });

    test('no match returns empty', () {
      final s = field.suggestions('xyz');
      expect(s, isEmpty);
    });
  });

  group('Elapsed suggestions', () {
    late lucene.Elapsed field;

    setUp(() {
      field = lucene.Elapsed.auto('runtime', Duration.zero, (_) {});
    });

    test('empty partial returns all duration ranges', () {
      expect(field.suggestions(''), isNotEmpty);
    });

    test('completions are prefixed with field name', () {
      final s = field.suggestions('');
      expect(s.every((s) => s.completion.startsWith('runtime:')), isTrue);
    });
  });
}
