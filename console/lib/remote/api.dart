import 'dart:convert';
import 'dart:async';
import 'dart:io';
// aliased: its Stream message would otherwise collide with dart:async's Stream, used unprefixed throughout this file.
import 'package:retrovibed/media/media.remote.control.pb.dart' as rc;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/retrovibed.dart' as retro;

export 'package:retrovibed/media/media.remote.control.pb.dart';

class RemoteControlSocket {
  final WebSocket _socket;

  RemoteControlSocket._(this._socket);

  Stream<rc.Stream> get messages => _socket.transform(
    StreamTransformer.fromHandlers(
      handleData: (data, sink) {
        if (data is List<int>) {
          final msg = rc.Stream.create()..mergeFromProto3Json(jsonDecode(utf8.decode(data)));
          sink.add(msg);
        } else {
          sink.addError('deserialization failed data: $data');
        }
      },
    ),
  );

  void send(rc.Stream msg) {
    _socket.add(utf8.encode(jsonEncode(msg.toProto3Json())));
  }

  Future<void> close() => _socket.close();
}

abstract class remotecontrol {
  // listen is always this device's own local frontend: authenticated with
  // the process-local token exposed over the native bridge, never the
  // profile bearer, and never valid outside this process.
  static Future<RemoteControlSocket> listen({List<httpx.Option> options = const []}) async {
    return httpx
        .websocket(
          Uri.https(httpx.host(), "/rc/listen", null),
          options: [
            httpx.Request.authorization("bearer ${retro.remote_control_listen_token()}"),
            ...options,
          ],
        )
        .then((socket) {
          socket.pingInterval = Duration(seconds: 10);
          return RemoteControlSocket._(socket);
        });
  }

  static Future<RemoteControlSocket> connect({List<httpx.Option> options = const []}) async {
    return httpx
        .websocket(Uri.https(httpx.host(), "/rc/connect", null), options: options)
        .then((socket) {
          socket.pingInterval = Duration(seconds: 10);
          return RemoteControlSocket._(socket);
        });
  }
}
