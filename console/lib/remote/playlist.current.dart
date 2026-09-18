import 'package:flutter/material.dart';
import 'package:retrovibed/media/media.row.display.dart' as rowdisplay;
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart' as remote;

class PlaylistCurrent extends StatelessWidget {
  final remote.Stream current;
  // when set, rows this Connect session itself enqueued (session_id
  // matches) are visually highlighted.
  final String sessionId;
  const PlaylistCurrent(this.current, {this.sessionId = "", Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ds.Container(
      rowdisplay.RowDisplay(
        media: current.asMedia,
        leading: const [Icon(Icons.play_arrow_rounded)],
      ),
    );
  }
}
