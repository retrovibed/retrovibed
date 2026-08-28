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

abstract class messages {
  static rc.Stream queue(media.Media m) {
    return rc.Stream(
      sid: uuidx.v7(),
      queue: rc.Queue(media: m),
    );
  }

  static rc.Stream dequeue(String id) {
    return rc.Stream(
      sid: uuidx.v7(),
      dequeue: rc.Dequeue(id: id),
    );
  }

  // pause has no payload - each command toggles the receiving device's
  // play/pause state; ordering against concurrent/stale commands is
  // resolved by the receiver using sid as a vector clock, same as Mute.
  static rc.Stream pause() {
    return rc.Stream(
      sid: uuidx.v7(),
      pause: rc.Pause(),
    );
  }

  static rc.Stream seek(int offset) {
    return rc.Stream(
      sid: uuidx.v7(),
      seek: rc.Seek(offset: offset),
    );
  }

  static rc.Stream previous() {
    return seek(SeekOffset.previous);
  }

  static rc.Stream next() {
    return seek(SeekOffset.next);
  }

  // relative volume adjustment (offset applied to the receiver's current
  // level, 0-100 scale) - reuses Seek's shape rather than setting an
  // absolute value.
  static rc.Stream volume(int offset) {
    return rc.Stream(
      sid: uuidx.v7(),
      volume: rc.Seek(offset: offset),
    );
  }

  // mute has no payload - each command toggles the receiving device's
  // audio between silent and its prior level; ordering against
  // concurrent/stale commands is resolved by the receiver using sid as a
  // vector clock, same as Fullscreen.
  static rc.Stream mute() {
    return rc.Stream(
      sid: uuidx.v7(),
      mute: rc.Mute(),
    );
  }

  // sync with no fields set requests the listener's current library and
  // playback queue; with fields set it reports the listener's current
  // library and playback queue, unsolicited or in reply to a request.
  static rc.Stream sync({
    meta.Daemon? library,
    String token = "",
    fixnum.Int64? expiration,
    int capacity = 0,
    media.Media? current,
    List<media.Media> queue = const [],
    double volume = 0,
    bool muted = false,
    bool paused = false,
    bool fullscreen = false,
    fixnum.Int64? vid,
  }) {
    return rc.Stream(
      sid: uuidx.v7(),
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
      ),
    );
  }

  // fullscreen has no payload - each command flips the receiving device's
  // current state; ordering against concurrent/stale commands is resolved
  // by the receiver using sid as a vector clock.
  static rc.Stream fullscreen() {
    return rc.Stream(
      sid: uuidx.v7(),
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
