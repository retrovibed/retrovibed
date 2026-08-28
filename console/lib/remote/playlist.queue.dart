import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media/media.row.display.dart' as rowdisplay;
import 'api.dart' as remote;

class PlaylistQueue extends StatelessWidget {
  final remote.Sync current;
  final remote.RemoteControlSocket socket;
  final void Function(remote.Sync Function(remote.Sync c) mutated) onChange;
  const PlaylistQueue(this.current, this.socket, {Key? key, this.onChange = _noop}) : super(key: key);

  static void _noop(remote.Sync Function(remote.Sync q) mutated) {}

  @override
  Widget build(BuildContext context) {
    if (current.queue.isEmpty) return const SizedBox.shrink();
    final defaults = ds.Defaults.of(context);
    return SingleChildScrollView(
      child: Column(
        verticalDirection: defaults.isCompact ? VerticalDirection.up : VerticalDirection.down,
        children: current.queue
            .map(
              (m) => rowdisplay.RowDisplay(
                media: m,
                leading: const [Icon(Icons.queue_music)],
                trailing: [
                  ds.LoadingIconButton.remove(
                    onPressed: () async {
                      onChange(remote.syncmut.dequeue(m));
                      socket.send(remote.messages.dequeue(m.id));
                    },
                  ),
                ],
              ),
            )
            .toList(),
      ),
    );
  }
}
