import 'package:flutter/material.dart';
import 'package:media_kit/media_kit.dart';
import 'package:retrovibed/designkit.dart' as ds;

class PlayerControlFullscreen extends StatelessWidget {
  final Player player;
  const PlayerControlFullscreen(this.player, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final current = ds.Full.nochrome(context);
    return IconButton(
      onPressed: () => ds.Full.toggle(context),
      icon: Icon(
        current ? Icons.fit_screen : Icons.crop_free,
      ),
      tooltip: current ? 'Exit Full Screen' : 'Full Screen',
    );
  }
}
