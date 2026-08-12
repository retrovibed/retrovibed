import 'package:flutter/material.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'api.dart' as remote;

class PlayerControlSeek extends StatelessWidget {
  final remote.RemoteControlSocket socket;
  final int offset;
  final IconData icon;

  const PlayerControlSeek._({Key? key, required this.socket, required this.offset, required this.icon}) : super(key: key);

  factory PlayerControlSeek.forward({Key? key, required remote.RemoteControlSocket socket, Duration step = const Duration(seconds: 10)}) {
    return PlayerControlSeek._(key: key, socket: socket, offset: step.inMilliseconds, icon: Icons.fast_forward_rounded);
  }

  factory PlayerControlSeek.backward({Key? key, required remote.RemoteControlSocket socket, Duration step = const Duration(seconds: 10)}) {
    return PlayerControlSeek._(key: key, socket: socket, offset: -step.inMilliseconds, icon: Icons.fast_rewind_rounded);
  }

  factory PlayerControlSeek.next({Key? key, required remote.RemoteControlSocket socket}) {
    return PlayerControlSeek._(key: key, socket: socket, offset: remote.SeekOffset.next, icon: Icons.skip_next_rounded);
  }

  factory PlayerControlSeek.prev({Key? key, required remote.RemoteControlSocket socket}) {
    return PlayerControlSeek._(key: key, socket: socket, offset: remote.SeekOffset.previous, icon: Icons.skip_previous_rounded);
  }

  @override
  Widget build(BuildContext context) {
    return IconButton(
      onPressed: () {
        socket.send(
          remote.Stream(
            sid: uuidx.random(),
            seek: remote.Seek(offset: offset),
          ),
        );
      },
      icon: Icon(icon),
    );
  }
}
