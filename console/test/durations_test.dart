import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/timex.dart' as timex;

void main() {
  group('durations.tryParse', () {
    test('parses hours', () {
      expect(timex.durations.tryParse('PT2H'), const Duration(hours: 2));
    });

    test('parses minutes', () {
      expect(timex.durations.tryParse('PT30M'), const Duration(minutes: 30));
    });

    test('parses seconds', () {
      expect(timex.durations.tryParse('PT45S'), const Duration(seconds: 45));
    });

    test('parses hours and minutes', () {
      expect(
        timex.durations.tryParse('PT1H30M'),
        const Duration(hours: 1, minutes: 30),
      );
    });

    test('parses days', () {
      expect(timex.durations.tryParse('P7D'), const Duration(days: 7));
    });

    test('parses days and hours', () {
      expect(
        timex.durations.tryParse('P1DT2H'),
        const Duration(days: 1, hours: 2),
      );
    });

    test('parses zero duration', () {
      expect(timex.durations.tryParse('PT0S'), Duration.zero);
    });

    test('invalid string returns default fallback', () {
      expect(timex.durations.tryParse('not-a-duration'), const Duration());
    });

    test('empty string returns default fallback', () {
      expect(timex.durations.tryParse(''), const Duration());
    });

    test('invalid string returns custom fallback', () {
      const custom = Duration(seconds: 99);
      expect(
        timex.durations.tryParse('not-a-duration', fallback: custom),
        custom,
      );
    });

    test('invalid string with null fallback returns null', () {
      expect(
        timex.durations.tryParse('not-a-duration', fallback: null),
        isNull,
      );
    });
  });

  group('durations.iso8601', () {
    test('formats hours', () {
      expect(timex.durations.iso8601(const Duration(hours: 2)), 'PT2H');
    });

    test('formats minutes', () {
      expect(timex.durations.iso8601(const Duration(minutes: 30)), 'PT30M');
    });

    test('formats seconds', () {
      expect(timex.durations.iso8601(const Duration(seconds: 45)), 'PT45S');
    });

    test('formats hours and minutes', () {
      expect(
        timex.durations.iso8601(const Duration(hours: 1, minutes: 30)),
        'PT1H30M',
      );
    });

    test('formats days', () {
      expect(timex.durations.iso8601(const Duration(days: 7)), 'P7D');
    });

    test('formats days and hours', () {
      expect(
        timex.durations.iso8601(const Duration(days: 1, hours: 2)),
        'P1DT2H',
      );
    });

    test('formats zero duration', () {
      expect(timex.durations.iso8601(Duration.zero), 'PT0S');
    });

    test('roundtrips through tryParse', () {
      const original = Duration(days: 2, hours: 3, minutes: 15, seconds: 30);
      final encoded = timex.durations.iso8601(original);
      expect(timex.durations.tryParse(encoded), original);
    });
  });
}
