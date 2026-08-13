import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart' as remote;

class PlayerControlPlayPause extends StatelessWidget {
  final remote.RemoteControlSocket socket;
  final Stream<remote.Stream> status;

  const PlayerControlPlayPause({Key? key, required this.socket, required this.status}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return StreamBuilder<remote.Stream>(
      stream: status,
      builder: (context, snapshot) {
        final msg = snapshot.data;
        final playing = msg == null || msg.whichCommand() != remote.Stream_Command.playpause || !msg.playpause.paused;
        final help = playing ? "pause playback on the remote device" : "resume playback on the remote device";
        return ds.LoadingIconButton(
          onPressed: ds.LoadingIconButton.convert(() => socket.send(remote.messages.playpause(playing))),
          icon: Icon(playing ? Icons.pause_rounded : Icons.play_arrow_rounded),
          tooltip: help,
          help: ds.Hint(Text(help)),
        );
      },
    );
  }
}
