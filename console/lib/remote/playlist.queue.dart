import 'package:flutter/material.dart';
import 'package:retrovibed/media/media.pb.dart' as media;
import 'package:retrovibed/media/media.row.display.dart' as rowdisplay;

class PlaylistQueue extends StatelessWidget {
  final List<media.Media> queue;
  const PlaylistQueue(this.queue, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    if (queue.isEmpty) return const SizedBox.shrink();
    return Column(
      children: queue
          .map((m) => rowdisplay.RowDisplay(media: m, leading: const [Icon(Icons.queue_music)]))
          .toList(),
    );
  }
}
