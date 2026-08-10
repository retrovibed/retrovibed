import 'package:flutter/material.dart';
import 'package:retrovibed/media/playlist.dart' as internal;

class PlayerControlNext extends StatelessWidget {
  const PlayerControlNext({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return IconButton(
      onPressed: () => internal.Playlist.of(context)?.next(),
      icon: Icon(Icons.skip_next_rounded),
    );
  }
}
