import 'package:flutter/material.dart';
import 'package:media_kit/media_kit.dart' as mediakit;
import 'package:fractal/mimex.dart' as mimex;
import './api.dart' as api;
import './media.pb.dart';
import './player.dart';

class ButtonPlay extends StatelessWidget {
  final Media current;
  const ButtonPlay({super.key, required this.current});

  @override
  Widget build(BuildContext context) {
    switch (mimex.icon(current.mimetype)) {
      case mimex.movie:
      case mimex.audio:
        final vscreen = VideoScreen.of(context);
        return IconButton(
          icon: Icon(Icons.play_circle_outline_rounded),
          onPressed:
              vscreen == null
                  ? null
                  : () {
                    vscreen.add(
                      mediakit.Media(api.media.download_uri(current.id)),
                    );
                  },
        );
      default:
        return Container();
    }
  }
}
