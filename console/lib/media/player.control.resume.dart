import 'package:flutter/material.dart';
import 'package:media_kit/media_kit.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/library.dart' as _library;

class PlayerControlResume extends StatelessWidget {
  final Player player;
  final ValueNotifier<bool> overlay;
  final _library.Known current;
  const PlayerControlResume(this.player, this.current, this.overlay, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);
    return ValueListenableBuilder<bool>(
      valueListenable: overlay,
      builder: (context, overlay, child) {
        return Material(
          color: theme.colorScheme.surface,
          child: InkWell(
            onTap: () => player.play(),
            child: Padding(
              padding: defaults.padding / 2,
              child: Row(
                spacing: defaults.spacing,
                children: [
                  Icon(Icons.play_circle_outline_rounded),
                  player.state.playing ? Text(current.description) : Text("Resume ${current.description}"),
                ],
              ),
            ),
          ),
        );
      },
    );
  }
}
