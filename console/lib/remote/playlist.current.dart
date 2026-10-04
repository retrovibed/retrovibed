import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:retrovibed/media/media.row.display.dart' as rowdisplay;
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart' as remote;
import 'player.control.playback.dart';

class PlaylistCurrent extends StatelessWidget {
  final remote.Stream current;
  final remote.RemoteControlSocket socket;
  final ValueListenable<remote.Playback> playback;
  // when set, rows this Connect session itself enqueued (session_id
  // matches) are visually highlighted.
  final String sessionId;
  const PlaylistCurrent(
    this.current, {
    required this.socket,
    required this.playback,
    this.sessionId = "",
    Key? key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return ds.Container(
      decoration: BoxDecoration(border: defaults.border),
      Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          rowdisplay.RowDisplay(
            padding: defaults.padding / 2,
            media: current.asMedia,
            leading: const [Icon(Icons.play_arrow_rounded)],
            border: const Border(),
          ),
          PlayerControlPlayback(
            padding: defaults.padding.copyWith(top: 0) / 2,
            socket: socket,
            sessionId: sessionId,
            current: playback,
          ),
        ],
      ),
    );
  }
}
