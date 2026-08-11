import 'package:flutter/material.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'api.dart' as remote;

class PlayerControlFastForward extends StatelessWidget {
  final remote.RemoteControlSocket socket;
  final Duration step;

  const PlayerControlFastForward({Key? key, required this.socket, this.step = const Duration(seconds: 10)})
    : super(key: key);

  @override
  Widget build(BuildContext context) {
    return IconButton(
      onPressed: () {
        socket.send(remote.Stream(sid: uuidx.random(), seek: remote.Seek(offset: step.inMilliseconds)));
      },
      icon: Icon(Icons.fast_forward_rounded),
    );
  }
}
