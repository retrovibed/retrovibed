import 'package:flutter/material.dart';
import 'package:retrovibed/media/media.row.display.dart' as rowdisplay;
import 'api.dart' as remote;

class PlaylistCurrent extends StatelessWidget {
  final remote.Stream current;
  // when set, rows this Connect session itself enqueued (session_id
  // matches) are visually highlighted.
  final String? mySessionId;
  const PlaylistCurrent(this.current, {this.mySessionId, Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    if (!current.asMedia.hasId()) return const SizedBox.shrink();
    return rowdisplay.RowDisplay(
      media: current.asMedia,
      leading: const [Icon(Icons.play_arrow_rounded)],
      highlighted: mySessionId != null && current.sessionId == mySessionId,
    );
  }
}
