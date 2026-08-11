import 'package:flutter/material.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'api.dart' as remote;

class PlayerControlNext extends StatelessWidget {
  final remote.RemoteControlSocket socket;

  const PlayerControlNext({Key? key, required this.socket}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return IconButton(
      onPressed: () {
        // int32 max is a sentinel meaning "skip to next track", see media.remote.control.proto's Seek.
        socket.send(remote.Stream(sid: uuidx.random(), seek: remote.Seek(offset: 0x7FFFFFFF)));
      },
      icon: Icon(Icons.skip_next_rounded),
    );
  }
}
