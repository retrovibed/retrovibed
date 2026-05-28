import 'dart:async';
import 'package:flutter/material.dart';
import 'package:media_kit/media_kit.dart'; // Provides [Player], [Media], [Playlist] etc.
import 'package:media_kit_video/media_kit_video.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/uuidx.dart' as uuidx;
import './playlist.dart' as internal;
import './player.control.previous.dart';
import './player.control.next.dart';
import './player.control.title.dart';
import './player.control.settings.dart';
import './player.control.filedrop.dart';
import './player.control.fullscreen.dart';
import './player.control.resume.dart';

class VideoScreen extends StatefulWidget {
  final Widget child;
  final Player player;
  final FocusNode focus;
  const VideoScreen(this.child, this.player, this.focus, {Key? key}) : super(key: key);

  static _VideoState? of(BuildContext context) {
    return context.findAncestorStateOfType<_VideoState>();
  }

  @override
  State<VideoScreen> createState() => _VideoState(player);
}

class _VideoState extends State<VideoScreen> {
  bool _playing = false;
  FocusScopeNode _selffocus = FocusScopeNode(
    debugLabel: "focus.video.player.scope",
  );
  final controller;
  late final StreamSubscription<Tracks> sub0;
  late final StreamSubscription<bool> sub1;

  _VideoState(Player player)
    : controller = VideoController(
        player,
        configuration: const bool.fromEnvironment('EG_VM')
            ? const VideoControllerConfiguration(enableHardwareAcceleration: false)
            : const VideoControllerConfiguration(),
      );

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void add(Media m) {
    widget.player.add(m).then((v) {
      widget.player.next();
    });
  }

  @override
  void initState() {
    super.initState();

    _playing = widget.player.state.playing;
    sub0 = widget.player.stream.tracks.listen((state) {
      setState(() {});
    });
    sub1 = widget.player.stream.playing.listen((playing) {
      setState(() {
        _playing = playing;
      });
    });

    widget.player.stream.log.listen((log) => print(log));
  }

  @override
  void dispose() {
    sub0.cancel();
    sub1.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Builder(
      builder: (context) {
        final theme = Theme.of(context);
        final defaults = ds.Defaults.of(context);
        final plist = internal.Playlist.of(context)!;
        final current = plist.known;
        final compact = defaults.isCompact;
        final controls = [
          PlayerControlPrevious(),
          MaterialPlayOrPauseButton(),
          PlayerControlNext(),
          MaterialDesktopVolumeButton(),
          SizedBox.square(dimension: defaults.spacing),
          Expanded(child: PlayerControlTitle(plist.current)),
          SizedBox.square(dimension: defaults.spacing),
          MaterialPositionIndicator(),
          SizedBox.square(dimension: defaults.spacing),
          PlayerControlSettings(widget.player),
          SizedBox.square(dimension: defaults.spacing),
          PlayerControlFiledrop(widget.player),
          SizedBox.square(dimension: defaults.spacing),
          PlayerControlFullscreen(widget.player),
        ];

        return FocusScope(
          node: _selffocus,
          child: Stack(
            fit: StackFit.passthrough,
            children: [
              MaterialDesktopVideoControlsTheme(
                normal: MaterialDesktopVideoControlsThemeData(
                  modifyVolumeOnScroll: false,
                  bottomButtonBar: controls,
                  keyboardShortcuts: {},
                ),
                fullscreen: MaterialDesktopVideoControlsThemeData(
                  modifyVolumeOnScroll: false,
                  bottomButtonBar: controls,
                  keyboardShortcuts: {},
                ),
                child: Video(focusNode: widget.focus, controller: controller),
              ),
              Visibility(
                maintainState: true,
                maintainFocusability: true,
                visible: !_playing,
                child: Column(
                  mainAxisSize: MainAxisSize.max,
                  mainAxisAlignment: MainAxisAlignment.start,
                  verticalDirection: compact ? VerticalDirection.up : VerticalDirection.down,
                  children: [
                    Expanded(
                      child: DecoratedBox(
                        decoration: BoxDecoration(
                          color: theme.scaffoldBackgroundColor.withValues(
                            alpha: defaults.opaque.a,
                          ),
                        ),
                        child: widget.child,
                      ),
                    ),
                    Visibility(
                      visible: !uuidx.isMin(uuidx.fromString(current.id)),
                      child: PlayerControlResume(widget.player, current),
                    ),
                  ],
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}
