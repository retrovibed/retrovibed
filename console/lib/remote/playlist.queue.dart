import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media/media.pb.dart' as media;
import 'package:retrovibed/media/media.row.display.dart' as rowdisplay;
import 'api.dart' as remote;

class PlaylistQueue extends StatelessWidget {
  final List<media.Media> queue;
  final remote.RemoteControlSocket socket;
  final void Function(media.Media queue) onChange;
  const PlaylistQueue(this.queue, this.socket, {Key? key, this.onChange = _noop}) : super(key: key);

  static void _noop(media.Media queue) {}

  @override
  Widget build(BuildContext context) {
    if (queue.isEmpty) return const SizedBox.shrink();
    final defaults = ds.Defaults.of(context);
    return Column(
      verticalDirection: defaults.isCompact ? VerticalDirection.up : VerticalDirection.down,
      children: queue
          .map(
            (m) => rowdisplay.RowDisplay(
              media: m,
              leading: const [Icon(Icons.queue_music)],
              trailing: [
                ds.LoadingIconButton.delete(
                  onPressed: () async {
                    socket.send(remote.messages.dequeue(m.id));
                    onChange(m);
                  },
                ),
              ],
            ),
          )
          .toList(),
    );
  }
}
