import 'package:flutter/material.dart';
import 'package:retrovibed/media/playlist.dart' as internal;

class PlayerControlPrevious extends StatelessWidget {
  const PlayerControlPrevious({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return IconButton(
      onPressed: () => internal.Playlist.of(context)?.previous(),
      icon: Icon(Icons.skip_previous_rounded),
    );
  }
}
