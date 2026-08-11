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
        if (msg == null || msg.whichCommand() != remote.Stream_Command.queue) {
          return const SizedBox.shrink();
        }
        return media.RowDisplay(
          media: msg.queue.media,
          leading: const [Icon(Icons.play_arrow_rounded)],
        );
      },
    );
  }
}
