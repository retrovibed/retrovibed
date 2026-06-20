import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './play.queue.dart';
import './media.row.display.dart' as media;

class PlayListHistory extends StatelessWidget {
  final BoxConstraints? constraints;
  final void Function(PlayableMedia media)? onTap;
  final List<PlayableMedia> history;
  final List<Widget> leading;
  final List<Widget> trailing;
  const PlayListHistory(
    this.history, {
    this.onTap,
    this.leading = const [],
    this.trailing = const [],
    this.constraints,
  });

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final played = history.reversed.toList();
    return ds.Container(
      padding: defaults.padding,
      margin: defaults.margin,
      constraints: constraints,
      Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          ...leading,
          if (leading.isNotEmpty) const Divider(),
          ...[
            for (final item in played)
              media.RowDisplay(
                media: item.current,
                leading: const [Icon(Icons.history)],
                onTap: onTap == null ? null : () async => onTap!(item),
              ),
          ],
          if (trailing.isNotEmpty) const Divider(),
          ...trailing,
        ],
      ),
    );
  }
}
