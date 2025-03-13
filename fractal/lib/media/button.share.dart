import 'package:flutter/material.dart';
import './media.pb.dart';

class ButtonShare extends StatelessWidget {
  final Media current;
  const ButtonShare({super.key, required this.current});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return IconButton(
      icon: Icon(Icons.share_outlined),
      color: current.torrentId == "" ? theme.disabledColor : null,
      onPressed: () {
        print("DERP DERP ${current.torrentId}");
      },
    );
  }
}
