import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/timex.dart' as timex;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/library/library.watch.pb.dart';

// Per-playback watch-history bookkeeping, shared by the local player
// (media/playlist.dart) and the remote-control listener (remote/connect.dart)
// so both heartbeat paths mint/accumulate identically. id is a fresh
// uuidx.v7() minted whenever the tracked media changes, so a rewatch later
// gets its own history row instead of overwriting the last one; watched is
// the wall-clock time accumulated between ticks (not derived from player
// position, so seeking doesn't distort it).
class WatchHistoryTracker {
  String id = uuidx.min();
  String mediaId = uuidx.min();
  Duration watched = Duration.zero;
  DateTime _lastTick = timex.neginf;

  // Advances bookkeeping for a heartbeat against mediaId and returns the
  // resulting record. Mints a fresh id and resets watched when mediaId
  // differs from the last tick; otherwise accumulates the wall-clock time
  // elapsed since the last tick, unless playing is false (e.g. paused), in
  // which case the timestamp still advances but no time is credited.
  // now is a test seam - defaults to the real clock.
  WatchHistoryRecord tick(String mediaId, {bool playing = true, DateTime? now}) {
    now ??= DateTime.now();
    if (mediaId != this.mediaId) {
      this.mediaId = mediaId;
      id = uuidx.v7();
      watched = Duration.zero;
    } else if (_lastTick != timex.neginf && playing) {
      watched += now.difference(_lastTick);
    }
    _lastTick = now;

    return WatchHistoryRecord(id: id, mediaId: this.mediaId, duration: ds.Int64(watched.inMilliseconds));
  }
}
