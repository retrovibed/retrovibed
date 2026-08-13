import 'package:flutter/material.dart';
import 'package:retrovibed/media/media.row.display.dart' as media;
import 'api.dart' as remote;

class PlaylistCurrent extends StatelessWidget {
  final Stream<remote.Stream> status;
  const PlaylistCurrent(this.status, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return StreamBuilder<remote.Stream>(
      stream: status,
      builder: (context, snapshot) {
        final msg = snapshot.data;
        if (msg == null) return const SizedBox.shrink();

        switch (msg.whichCommand()) {
          case remote.Stream_Command.queue:
            return media.RowDisplay(
              media: msg.queue.media,
              leading: const [Icon(Icons.play_arrow_rounded)],
            );
          case remote.Stream_Command.sync:
            if (!msg.sync.hasCurrent()) return const SizedBox.shrink();
            return media.RowDisplay(
              media: msg.sync.current,
              leading: const [Icon(Icons.play_arrow_rounded)],
            );
          default:
            return const SizedBox.shrink();
        }
      },
    );
  }
}
