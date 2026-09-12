import 'dart:convert';
import 'dart:async';
import 'dart:io';
import 'package:fixnum/fixnum.dart' as fixnum;
// aliased: its Stream message would otherwise collide with dart:async's Stream, used unprefixed throughout this file.
import 'package:retrovibed/media/media.remote.control.pb.dart' as rc;
import 'package:retrovibed/media/media.pb.dart' as media;
import 'package:retrovibed/meta/meta.daemon.pb.dart' as meta;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/retrovibed.dart' as retro;

export 'package:retrovibed/media/media.remote.control.pb.dart';

// Sentinel Seek.offset values meaning "skip to next/previous track" rather
// than a relative seek, per media.remote.control.proto's Seek.
abstract class SeekOffset {
  static const int next = 0x7FFFFFFF; // int32 max
  static const int previous = -0x80000000; // int32 min
}

abstract class RemoteControlSocket {
  // placeholder used when we'll no longer be attempting to connect to the socket.
  static final RemoteControlSocket disabled = _NoopRemoteControlSocket();
  // placeholder used before a real connection exists, so callers never need
  // to null-check (e.g. `_socket!`) while waiting to connect.
  static final RemoteControlSocket noop = _NoopRemoteControlSocket();

  Stream<rc.Stream> get messages;
  void send(rc.Stream msg);
  Future<void> close();
}

class _NoopRemoteControlSocket implements RemoteControlSocket {
  @override
  Stream<rc.Stream> get messages => const Stream.empty();

  @override
  void send(rc.Stream msg) {}

  @override
  Future<void> close() => Future.value();
}

class _WebSocketRemoteControlSocket implements RemoteControlSocket {
  final WebSocket _socket;

  _WebSocketRemoteControlSocket(this._socket);

  @override
  Stream<rc.Stream> get messages => _socket.transform(
    StreamTransformer.fromHandlers(
      handleData: (data, sink) {
        if (data is List<int>) {
          final msg = httpx.fromProto3JsonSafe(rc.Stream.create(), jsonDecode(utf8.decode(data)));
          sink.add(msg);
        } else {
          sink.addError('deserialization failed data: $data');
        }
      },
      handleDone: (sink) {
        print("websocket closed: code=${_socket.closeCode} reason=${_socket.closeReason}");
        sink.close();
      },
    ),
  );

  @override
  void send(rc.Stream msg) {
    _socket.add(utf8.encode(jsonEncode(msg.toProto3Json())));
  }

  @override
  Future<void> close() => _socket.close();
}

// Stream.asMedia unwraps the media a Sync current/queue entry carries -
// current/queue are now the original enqueue Stream frame (carrying
// sid/profile_id/session_id alongside the media), not bare Media, so every
// reader needs this one extra hop. Named asMedia (not media) since a
// getter named media would collide with this file's own `media` import
// prefix.
extension StreamMediaX on rc.Stream {
  media.Media get asMedia => queue.media;
}

abstract class syncmut {
  // profileId is deliberately not a param here - it's server-authoritative
  // (stamped by shallows/mediaapi/http.remote.control.go from the caller's
  // JWT), never client-set. This optimistic local entry leaves it unset
  // until the real echoed Sync arrives with the server-stamped value.
  static rc.Sync Function(rc.Sync) queue(media.Media m, {required String sessionId}) {
    return (v) {
      v.queue.add(rc.Stream(sid: uuidx.v7(), sessionId: sessionId, queue: rc.Queue(media: m)));
      return v;
    };
  }

  static rc.Sync Function(rc.Sync) dequeue(media.Media m) {
    return (v) {
      v.queue.removeWhere((s) => s.asMedia.id == m.id);
      return v;
    };
  }
}

abstract class messages {
  static rc.Stream queue(media.Media m, {required String sessionId}) {
    return rc.Stream(
      sid: uuidx.v7(),
      sessionId: sessionId,
      queue: rc.Queue(media: m),
    );
  }

  // ambiguous once the same media can be queued twice by different
  // sessions/profiles (dequeue targets every queue entry with this media
  // id) - each Sync.queue entry now carries its own sid, so a future fix
  // could dequeue by that instead. Not addressed here.
  static rc.Stream dequeue(String id, {required String sessionId}) {
    return rc.Stream(
      sid: uuidx.v7(),
      sessionId: sessionId,
      dequeue: rc.Dequeue(id: id),
    );
  }

  // pause has no payload - each command toggles the receiving device's
  // play/pause state; ordering against concurrent/stale commands is
  // resolved by the receiver using sid as a vector clock, same as Mute.
  static rc.Stream pause({required String sessionId}) {
    return rc.Stream(
      sid: uuidx.v7(),
      sessionId: sessionId,
      pause: rc.Pause(),
    );
  }

  static rc.Stream seek(int offset, {required String sessionId}) {
    return rc.Stream(
      sid: uuidx.v7(),
      sessionId: sessionId,
      seek: rc.Seek(offset: offset),
    );
  }

  static rc.Stream previous({required String sessionId}) {
    return seek(SeekOffset.previous, sessionId: sessionId);
  }

  static rc.Stream next({required String sessionId}) {
    return seek(SeekOffset.next, sessionId: sessionId);
  }

  // relative volume adjustment (offset applied to the receiver's current
  // level, 0-100 scale) - reuses Seek's shape rather than setting an
  // absolute value.
  static rc.Stream volume(int offset, {required String sessionId}) {
    return rc.Stream(
      sid: uuidx.v7(),
      sessionId: sessionId,
      volume: rc.Seek(offset: offset),
    );
  }

  // mute has no payload - each command toggles the receiving device's
  // audio between silent and its prior level; ordering against
  // concurrent/stale commands is resolved by the receiver using sid as a
  // vector clock, same as Fullscreen.
  static rc.Stream mute({required String sessionId}) {
    return rc.Stream(
      sid: uuidx.v7(),
      sessionId: sessionId,
      mute: rc.Mute(),
    );
  }

  // sync with no fields set requests the listener's current library and
  // playback queue; with fields set it reports the listener's current
  // library and playback queue, unsolicited or in reply to a request.
  // current/queue are Stream (not bare Media) so provenance survives the
  // round trip - see rc.Sync's proto comment.
  static rc.Stream sync({
    required String sessionId,
    meta.Daemon? library,
    String token = "",
    fixnum.Int64? expiration,
    int capacity = 0,
    rc.Stream? current,
    List<rc.Stream> queue = const [],
    double volume = 0,
    bool muted = false,
    bool paused = false,
    bool fullscreen = false,
    fixnum.Int64? vid,
    rc.Playback? playback,
  }) {
    return rc.Stream(
      sid: uuidx.v7(),
      sessionId: sessionId,
      vid: vid,
      sync: rc.Sync(
        library: library,
        token: token,
        expiration: expiration,
        capacity: capacity,
        current: current,
        queue: queue,
        volume: volume,
        muted: muted,
        paused: paused,
        fullscreen: fullscreen,
        playback: playback,
      ),
    );
  }

  // reports live position/duration for the session's current pick - a
  // lighter frame than sync (no queue/library/token), sent periodically
  // while something is playing.
  static rc.Stream playback({
    required String sessionId,
    required String profileId,
    required media.Media media,
    required fixnum.Int64 position,
    required fixnum.Int64 duration,
  }) {
    return rc.Stream(
      sid: uuidx.v7(),
      sessionId: sessionId,
      profileId: profileId,
      playback: rc.Playback(media: media, position: position, duration: duration),
    );
  }

  // fullscreen has no payload - each command flips the receiving device's
  // current state; ordering against concurrent/stale commands is resolved
  // by the receiver using sid as a vector clock.
  static rc.Stream fullscreen({required String sessionId}) {
    return rc.Stream(
      sid: uuidx.v7(),
      sessionId: sessionId,
      fullscreen: rc.Fullscreen(),
    );
  }
}

abstract class remotecontrol {
  // listen is always this device's own local frontend: authenticated with
  // the process-local token exposed over the native bridge, never the
  // profile bearer, and never valid outside this process.
  static Future<RemoteControlSocket> listen({List<httpx.Option> options = const []}) async {
    return httpx
        .websocket(
          Uri.https(httpx.localhost(), "/rc/listen", null),
          options: [
            httpx.Request.authorization("bearer ${retro.remote_control_listen_token()}"),
            ...options,
          ],
        )
        .then((socket) {
          socket.pingInterval = Duration(seconds: 10);
          return _WebSocketRemoteControlSocket(socket);
        });
  }

  static Future<RemoteControlSocket> connect({required String host, List<httpx.Option> options = const []}) async {
    return httpx.websocket(Uri.https(host, "/rc/connect", null), options: options).then((socket) {
      socket.pingInterval = Duration(seconds: 10);
      return _WebSocketRemoteControlSocket(socket);
    });
  }
}
