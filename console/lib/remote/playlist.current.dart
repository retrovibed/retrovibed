import 'package:flutter/material.dart';
import 'package:retrovibed/media/media.pb.dart' as media;
import 'package:retrovibed/media/media.row.display.dart' as rowdisplay;

class PlaylistCurrent extends StatelessWidget {
  final media.Media current;
  const PlaylistCurrent(this.current, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    if (!current.hasId()) return const SizedBox.shrink();
    return rowdisplay.RowDisplay(
      media: current,
      leading: const [Icon(Icons.play_arrow_rounded)],
    );
  }
}
