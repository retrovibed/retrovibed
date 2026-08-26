import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart' as remote;

class PlayerControlFullscreen extends StatelessWidget {
  final remote.RemoteControlSocket socket;
  final bool current;

  const PlayerControlFullscreen({Key? key, required this.socket, required this.current}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final help = current ? 'Exit Full Screen' : 'Full Screen';
    return ds.LoadingIconButton(
      onPressed: ds.LoadingIconButton.convert(() => socket.send(remote.messages.fullscreen())),
      icon: Icon(current ? Icons.fit_screen : Icons.crop_free),
      tooltip: help,
      help: ds.Hint(Text(help)),
    );
  }
}
