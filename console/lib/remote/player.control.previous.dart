import 'package:flutter/material.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'api.dart' as remote;

class PlayerControlPrevious extends StatelessWidget {
  final remote.RemoteControlSocket socket;

  const PlayerControlPrevious({Key? key, required this.socket}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return IconButton(
      onPressed: () {
        // int32 min is a sentinel meaning "skip to previous track", see media.remote.control.proto's Seek.
        socket.send(remote.Stream(sid: uuidx.random(), seek: remote.Seek(offset: -0x80000000)));
      },
      icon: Icon(Icons.skip_previous_rounded),
    );
  }
}
