import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart' as remote;

class PlayerControlSeek extends StatelessWidget {
  final remote.RemoteControlSocket socket;
  final int offset;
  final IconData icon;
  final String help;

  const PlayerControlSeek._({
    Key? key,
    required this.socket,
    required this.offset,
    required this.icon,
    required this.help,
  }) : super(key: key);

  factory PlayerControlSeek.forward({
    Key? key,
    required remote.RemoteControlSocket socket,
    Duration step = const Duration(seconds: 10),
  }) {
    return PlayerControlSeek._(
      key: key,
      socket: socket,
      offset: step.inMilliseconds,
      icon: Icons.fast_forward_rounded,
      help: "seek forward ${step.inSeconds}s on the remote device",
    );
  }

  factory PlayerControlSeek.backward({
    Key? key,
    required remote.RemoteControlSocket socket,
    Duration step = const Duration(seconds: 10),
  }) {
    return PlayerControlSeek._(
      key: key,
      socket: socket,
      offset: -step.inMilliseconds,
      icon: Icons.fast_rewind_rounded,
      help: "seek backward ${step.inSeconds}s on the remote device",
    );
  }

  factory PlayerControlSeek.next({Key? key, required remote.RemoteControlSocket socket}) {
    return PlayerControlSeek._(
      key: key,
      socket: socket,
      offset: remote.SeekOffset.next,
      icon: Icons.skip_next_rounded,
      help: "skip to the next track on the remote device",
    );
  }

  factory PlayerControlSeek.prev({Key? key, required remote.RemoteControlSocket socket}) {
    return PlayerControlSeek._(
      key: key,
      socket: socket,
      offset: remote.SeekOffset.previous,
      icon: Icons.skip_previous_rounded,
      help: "skip to the previous track on the remote device",
    );
  }

  @override
  Widget build(BuildContext context) {
    return ds.LoadingIconButton(
      onPressed: ds.LoadingIconButton.convert(() => socket.send(remote.messages.seek(offset))),
      icon: Icon(icon),
      tooltip: help,
      help: ds.Hint(Text(help)),
    );
  }
}
