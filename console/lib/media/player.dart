import 'dart:async';
import 'package:flutter/material.dart';
import 'package:media_kit/media_kit.dart'; // Provides [Player], [Media], [Playlist] etc.
import 'package:media_kit_video/media_kit_video.dart';
import 'package:console/designkit.dart' as ds;

class VideoScreen extends StatefulWidget {
  final Widget child;
  final Player player;
  const VideoScreen(this.child, this.player, {Key? key}) : super(key: key);

  static _VideoState? of(BuildContext context) {
    return context.findAncestorStateOfType<_VideoState>();
  }

  @override
  State<VideoScreen> createState() => _VideoState();
}

class _VideoState extends State<VideoScreen> {
  Widget _resume = SizedBox();

  // Create a [VideoController] to handle video output from [Player].
  late final controller = VideoController(widget.player);
  late final StreamSubscription<bool> subscription;

  void add(Media m) {
    widget.player.add(m).then((v) {
      widget.player.next();
    });
  }

  @override
  void initState() {
    super.initState();
    subscription = widget.player.stream.playing.listen((state) {
      Widget resumew = SizedBox();

      if (widget.player.state.playlist.medias.length > 0) {
        final _m =
            widget.player.state.playlist.medias[widget
                .player
                .state
                .playlist
                .index];
        final _title = _m.extras?["title"] ?? "";
        resumew = IconButton(
          onPressed: () {
            widget.player.play();
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
        if (!super.mounted) return;
        _resume = resumew;
      });
    });
  }

  @override
  void dispose() {
    subscription.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final themex = ds.Defaults.of(context);

    final m =
        widget.player.state.playing
            ? SizedBox()
            : Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Expanded(
                  child: DecoratedBox(
                    decoration: BoxDecoration(
                      color: theme.scaffoldBackgroundColor.withValues(
                        alpha: themex.opaque?.a ?? 0.0,
                      ),
                    ),
                    child: widget.child,
                  ),
                ),
                DecoratedBox(
                  decoration: BoxDecoration(
                    color: theme.scaffoldBackgroundColor,
                  ),
                  child: _resume,
                ),
              ],
            );

    return Stack(
      children: [ds.Full(Center(child: Video(controller: controller))), m],
    );
  }
}
