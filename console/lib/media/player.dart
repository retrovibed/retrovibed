import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:media_kit/media_kit.dart'; // Provides [Player], [Media], [Playlist] etc.
import 'package:media_kit_video/media_kit_video.dart';
import 'package:console/designkit.dart' as ds;

class VideoScreen extends StatefulWidget {
  final Widget child;
  const VideoScreen(this.child, {Key? key}) : super(key: key);
  static _VideoState? of(BuildContext context) {
    return context.findAncestorStateOfType<_VideoState>();
  }

  @override
  State<VideoScreen> createState() => _VideoState();
}

class _VideoState extends State<VideoScreen> {
  final FocusNode _focusNode = FocusNode();
  Widget _resume = SizedBox();

  // Create a [Player] to control playback.
  final player = Player();

  // Create a [VideoController] to handle video output from [Player].
  late final controller = VideoController(player);

  void add(Media m) {
    player.add(m).then((v) {
      player.next();
    });
  }

  @override
  void initState() {
    super.initState();
    player.stream.playing.listen((state) {
      Widget resumew = SizedBox();
      if (player.state.playlist.medias.length > 0) {
        final _m = player.state.playlist.medias[player.state.playlist.index];
        final _title = _m.extras?["title"] ?? "";
        resumew = IconButton(
          onPressed: () {
            player.play();
          },
          icon: Row(
            spacing: 10.0,
            children: [
              Icon(Icons.play_circle_outline_rounded),
              Text("Resume ${_title}"),
            ],
          ),
        );
      }

      setState(() {
        _resume = resumew;
      });
    });
  }

  @override
  void dispose() {
    player.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final m =
        player.state.playing
            ? null
            : Column(
              mainAxisSize: MainAxisSize.min,
              children: [Expanded(child: widget.child), _resume],
            );

    return KeyboardListener(
      focusNode: _focusNode,
      autofocus: player.state.playing,
      onKeyEvent: (event) {
        if (event is KeyDownEvent) {
          if (event.logicalKey == LogicalKeyboardKey.space) {
            player.playOrPause();
          }
        }
      },
      child: ds.Overlay(
        child: ds.Full(Center(child: Video(controller: controller))),
        overlay: m,
      ),
    );
  }
}
