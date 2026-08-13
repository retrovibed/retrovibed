import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart' as remote;

class PlayerControlSync extends StatelessWidget {
  static const String help = "request the remote device's current library, token, and playback queue";

  final remote.RemoteControlSocket socket;

  const PlayerControlSync({Key? key, required this.socket}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ds.LoadingIconButton(
      onPressed: ds.LoadingIconButton.convert(() => socket.send(remote.messages.sync())),
      icon: const Icon(Icons.sync_rounded),
      tooltip: help,
      help: ds.Hint(const Text(help)),
    );
  }
}
