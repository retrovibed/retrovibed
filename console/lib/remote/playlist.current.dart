import 'package:flutter/material.dart';
import 'package:retrovibed/media/media.row.display.dart' as media;
import 'package:retrovibed/media/play.queue.dart';

class PlaylistCurrent extends StatelessWidget {
  final ValueNotifier<PlayableMedia?> current;
  const PlaylistCurrent(this.current, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder<PlayableMedia?>(
      valueListenable: current,
      builder: (context, playing, _) {
        if (playing == null) return const SizedBox.shrink();
        return media.RowDisplay(
          media: playing.current,
          leading: const [Icon(Icons.play_arrow_rounded)],
        );
      },
    );
  }
}
