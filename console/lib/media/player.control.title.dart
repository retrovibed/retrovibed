import 'package:flutter/material.dart';
import './play.queue.dart';

class PlayerControlTitle extends StatelessWidget {
  final ValueNotifier<PlayableMedia?> events;
  const PlayerControlTitle(this.events, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder<PlayableMedia?>(
      valueListenable: events,
      builder: (context, media, _) {
        final description = media?.current.description ?? '';
        return Text(description, maxLines: 1, overflow: TextOverflow.ellipsis);
      },
    );
  }
}
