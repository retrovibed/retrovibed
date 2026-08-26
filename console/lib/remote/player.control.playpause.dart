import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart' as remote;

class PlayerControlPlayPause extends StatelessWidget {
  final remote.RemoteControlSocket socket;
  final bool paused;

  const PlayerControlPlayPause({Key? key, required this.socket, required this.paused}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final playing = !paused;
    final help = playing ? "pause playback on the remote device" : "resume playback on the remote device";
    return ds.LoadingIconButton(
      onPressed: ds.LoadingIconButton.convert(() => socket.send(remote.messages.pause())),
      icon: Icon(playing ? Icons.pause_rounded : Icons.play_arrow_rounded),
      tooltip: help,
      help: ds.Hint(Text(help)),
    );
  }
}
