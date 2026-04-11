import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/timex.dart' as timex;

void main() {
  group('Range', () {
    final a = DateTime.utc(2024, 1, 1);
    final b = DateTime.utc(2024, 6, 15);
    final c = DateTime.utc(2025, 1, 1);

    test('equal when begin and end are identical', () {
      final r1 = timex.Range(a, b);
      final r2 = timex.Range(a, b);
      expect(r1, equals(r2));
    });

    test('not equal when begin differs', () {
      final r1 = timex.Range(a, b);
      final r2 = timex.Range(c, b);
      expect(r1, isNot(equals(r2)));
    });

    test('not equal when end differs', () {
      final r1 = timex.Range(a, b);
      final r2 = timex.Range(a, c);
      expect(r1, isNot(equals(r2)));
    });

    test('not equal when both differ', () {
      final r1 = timex.Range(a, b);
      final r2 = timex.Range(b, c);
      expect(r1, isNot(equals(r2)));
    });

    test('identical instance is equal', () {
      final r1 = timex.Range(a, b);
      expect(r1, equals(r1));
    });

    test('hashCode matches for equal ranges', () {
      final r1 = timex.Range(a, b);
      final r2 = timex.Range(a, b);
      expect(r1.hashCode, equals(r2.hashCode));
    });

    test('hashCode differs for unequal ranges', () {
      final r1 = timex.Range(a, b);
      final r2 = timex.Range(a, c);
      expect(r1.hashCode, isNot(equals(r2.hashCode)));
    });

    test('works as map key', () {
      final r1 = timex.Range(a, b);
      final r2 = timex.Range(a, b);
      final map = {r1: 'value'};
      expect(map[r2], equals('value'));
    });

    test('works in set deduplication', () {
      final r1 = timex.Range(a, b);
      final r2 = timex.Range(a, b);
      final r3 = timex.Range(a, c);
      final set = {r1, r2, r3};
      expect(set.length, equals(2));
    });
  });

  group('Range.latest', () {
    test('end is after begin', () {
      final r = timex.Range.latest(const Duration(days: 30));
      expect(r.end.isAfter(r.begin), isTrue);
    });

    test('duration between begin and end matches input', () {
      const d = Duration(days: 30);
      final r = timex.Range.latest(d);
      final diff = r.end.difference(r.begin);
      expect(diff, equals(d));
    });

    test('zero duration produces begin equal to end', () {
      final r = timex.Range.latest(Duration.zero);
      expect(r.begin, equals(r.end));
    });

    test('end is approximately now', () {
      final before = DateTime.now().toUtc();
      final r = timex.Range.latest(const Duration(hours: 1));
      final after = DateTime.now().toUtc();
      expect(r.end.isAfter(before) || r.end.isAtSameMomentAs(before), isTrue);
      expect(r.end.isBefore(after) || r.end.isAtSameMomentAs(after), isTrue);
    });

    test('begin is end minus duration', () {
      const d = Duration(days: 7);
      final r = timex.Range.latest(d);
      expect(r.begin, equals(r.end.subtract(d)));
    });
  });
}
