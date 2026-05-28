import 'package:flutter/material.dart';
import 'package:media_kit/media_kit.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './player.settings.dart';

class PlayerControlSettings extends StatelessWidget {
  final Player player;
  const PlayerControlSettings(this.player, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return ds.Help(
      IconButton(
        onPressed: () {
          ds.modals.push(
            context,
            PlayerSettings(
              constraints: BoxConstraints(maxWidth: 256),
              padding: defaults.padding,
              current: player,
            ),
          );
        },
        icon: Icon(Icons.tune),
      ),
      ds.Hint(const Text("open settings to select audio, video, or subtitle tracks")),
    );
  }
}
