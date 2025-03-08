import 'package:flutter/material.dart';
import 'package:media_kit/media_kit.dart'; // Provides [Player], [Media], [Playlist] etc.
import 'package:media_kit_video/media_kit_video.dart';
import 'package:fractal/designkit.dart' as ds;

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
      setState(() {});
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
            ? ds.Full(Center(child: Video(controller: controller)))
            : null;
    return ds.Overlay(child: widget.child, overlay: m);
  }
}
