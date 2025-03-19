import 'package:flutter/material.dart';
import 'package:media_kit/media_kit.dart' as mediakit;
import 'package:console/mimex.dart' as mimex;
import './api.dart' as api;
import './media.pb.dart';
import './player.dart';

mediakit.Media PlayableMedia(Media current) {
  return mediakit.Media(
    api.media.download_uri(current.id),
    extras: Map.of(<String, String>{
      "id": current.id,
      "title": current.description,
    }),
  );
}

void Function()? PlayAction(BuildContext context, Media current) {
  switch (mimex.icon(current.mimetype)) {
    case mimex.movie:
    case mimex.audio:
      final vscreen = VideoScreen.of(context);
      return vscreen == null
          ? null
          : () {
            vscreen.add(PlayableMedia(current));
          };
    default:
      return null;
  }
}

class ButtonPlay extends StatelessWidget {
  final Media current;
  const ButtonPlay({super.key, required this.current});

  @override
  Widget build(BuildContext context) {
    switch (mimex.icon(current.mimetype)) {
      case mimex.movie:
      case mimex.audio:
        return IconButton(
          icon: Icon(Icons.play_circle_outline_rounded),
          onPressed: PlayAction(context, current),
        );
      default:
        return Container();
    }
  }
}
