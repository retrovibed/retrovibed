import 'package:flutter/material.dart';
import 'package:retrovibed/media/playlist.dart' as internal;

class PlayerControlFastForward extends StatelessWidget {
  final Duration step;
  const PlayerControlFastForward({Key? key, this.step = const Duration(seconds: 10)}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return IconButton(
      onPressed: () {
        final playlist = internal.Playlist.of(context);
        if (playlist == null) return;
        playlist.player.seek(playlist.player.state.position + step);
      },
      icon: Icon(Icons.fast_forward_rounded),
    );
  }
}
