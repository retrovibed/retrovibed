import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/timex.dart' as timex;

void main() {
  group('download.completed', () {
    test('true when completedAt is a past timestamp', () {
      final media.Download d = media.Download(completedAt: '2025-06-01T12:00:00Z');
      expect(media.download.completed(d), isTrue);
    });

    test('false when completedAt is timex.inf', () {
      final media.Download d = media.Download(completedAt: timex.formatISO8601(timex.inf));
      expect(media.download.completed(d), isFalse);
    });

    test('false when completedAt is unset', () {
      final media.Download d = media.Download();
      expect(d.completedAt, isEmpty);
      expect(media.download.completed(d), isFalse);
    });

    test('false when completedAt is explicitly empty', () {
      final media.Download d = media.Download(completedAt: '');
      expect(media.download.completed(d), isFalse);
    });
  });
}
