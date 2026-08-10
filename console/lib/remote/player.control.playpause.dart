import 'package:flutter/material.dart';
import 'package:retrovibed/media/playlist.dart' as internal;

class PlayerControlPlayPause extends StatelessWidget {
  const PlayerControlPlayPause({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final playlist = internal.Playlist.of(context);
    if (playlist == null) return const SizedBox.shrink();
    return StreamBuilder<bool>(
      stream: playlist.player.stream.playing,
      initialData: playlist.player.state.playing,
      builder: (context, snapshot) {
        final playing = snapshot.data ?? false;
        return IconButton(
          onPressed: () => playlist.player.playOrPause(),
          icon: Icon(playing ? Icons.pause_rounded : Icons.play_arrow_rounded),
        );
      },
    );
  }
}
