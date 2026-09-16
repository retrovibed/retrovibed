import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/library.dart' as lib;
import 'package:retrovibed/uuidx.dart' as uuidx;

void main() {
  group('WatchHistoryTracker.tick', () {
    test('starts at the nil sentinel before the first tick', () {
      final t = lib.WatchHistoryTracker();
      expect(t.id, uuidx.min());
      expect(t.mediaId, uuidx.min());
      expect(t.watched, Duration.zero);
    });

    test('first tick for a media id mints a fresh id and zero duration', () {
      final t = lib.WatchHistoryTracker();
      final rec = t.tick('media-a', now: DateTime(2026, 1, 1, 0, 0, 0));

      expect(rec.id, isNot(uuidx.min()));
      expect(rec.mediaId, 'media-a');
      expect(rec.duration.toInt(), 0);
      expect(t.id, rec.id);
    });

    test('repeat ticks for the same media id accumulate elapsed time', () {
      final t = lib.WatchHistoryTracker();
      final start = DateTime(2026, 1, 1, 0, 0, 0);
      t.tick('media-a', now: start);
      final rec = t.tick('media-a', now: start.add(const Duration(seconds: 5)));

      expect(rec.duration.toInt(), 5000);
    });

    test('accumulates across multiple ticks', () {
      final t = lib.WatchHistoryTracker();
      final start = DateTime(2026, 1, 1, 0, 0, 0);
      t.tick('media-a', now: start);
      t.tick('media-a', now: start.add(const Duration(seconds: 2)));
      final rec = t.tick('media-a', now: start.add(const Duration(seconds: 5)));

      expect(rec.duration.toInt(), 5000);
    });

    test('paused ticks advance the clock without crediting time', () {
      final t = lib.WatchHistoryTracker();
      final start = DateTime(2026, 1, 1, 0, 0, 0);
      t.tick('media-a', now: start);
      t.tick('media-a', now: start.add(const Duration(seconds: 3)), playing: false);
      final rec = t.tick('media-a', now: start.add(const Duration(seconds: 5)));

      // the paused tick's own gap (0s -> 3s) is dropped; only the gap since
      // that tick (3s -> 5s) is credited.
      expect(rec.duration.toInt(), 2000);
    });

    test('a different media id mints a new id and resets watched to zero', () {
      final t = lib.WatchHistoryTracker();
      final start = DateTime(2026, 1, 1, 0, 0, 0);
      final first = t.tick('media-a', now: start);
      t.tick('media-a', now: start.add(const Duration(seconds: 5)));

      final second = t.tick('media-b', now: start.add(const Duration(seconds: 6)));

      expect(second.id, isNot(first.id));
      expect(second.mediaId, 'media-b');
      expect(second.duration.toInt(), 0);
    });

    test('resuming the earlier media id after a switch starts a new playback', () {
      final t = lib.WatchHistoryTracker();
      final start = DateTime(2026, 1, 1, 0, 0, 0);
      final first = t.tick('media-a', now: start);
      t.tick('media-b', now: start.add(const Duration(seconds: 1)));
      final rewatch = t.tick('media-a', now: start.add(const Duration(seconds: 2)));

      expect(rewatch.id, isNot(first.id), reason: 'a rewatch must get its own history row');
      expect(rewatch.duration.toInt(), 0);
    });
  });
}
