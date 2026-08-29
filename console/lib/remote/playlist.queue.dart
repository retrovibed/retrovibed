import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'api.dart' as remote;

class PlaylistQueue extends StatelessWidget {
  final remote.Sync current;
  final remote.RemoteControlSocket socket;
  final void Function(remote.Sync Function(remote.Sync c) mutated) onChange;
  final Widget empty;
  const PlaylistQueue(
    this.current,
    this.socket, {
    Key? key,
    this.onChange = _noop,
    this.empty = ds.Empty,
  }) : super(key: key);

  static void _noop(remote.Sync Function(remote.Sync q) mutated) {}

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    if (current.queue.isEmpty) return empty;
    return SingleChildScrollView(
      child: Column(
        verticalDirection: defaults.isCompact ? VerticalDirection.up : VerticalDirection.down,
        children: current.queue
            .map(
              (m) => media.RowDisplay(
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
