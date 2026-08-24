import 'dart:async';
import 'package:flutter/material.dart';
import 'package:media_kit/media_kit.dart'; // Provides [Player], [Media], [Playlist] etc.
import 'package:media_kit_video/media_kit_video.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/debug.dart' as debug;
import 'package:retrovibed/media/player.control.shuffle.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;
import './playlist.dart' as internal;
import './player.control.previous.dart';
import './player.control.next.dart';
import './player.control.title.dart';
import './player.control.settings.dart';
import './player.control.filedrop.dart';
import './player.control.stream.dart';
import './player.control.fullscreen.dart';
import './player.control.resume.dart';

class VideoScreen extends StatefulWidget {
  final Widget child;
  final Player player;
  final FocusNode focus;
  final ValueNotifier<bool> overlay;
  const VideoScreen(this.child, this.player, this.focus, this.overlay, {Key? key}) : super(key: key);

  static _VideoState? of(BuildContext context) {
    return context.findAncestorStateOfType<_VideoState>();
  }

  @override
  State<VideoScreen> createState() => _VideoState(player);
}

class _VideoState extends State<VideoScreen> {
  FocusScopeNode _selffocus = FocusScopeNode(
    debugLabel: "focus.video.player.scope",
  );
  final controller;
  late final StreamSubscription<Tracks> subtracks;
  late final StreamSubscription<String> suberrors;
  late final StreamSubscription<bool> subplaying;
  final ValueNotifier<Widget> errors = ValueNotifier<Widget>(ds.Error.zero);
  // Stable list so MaterialDesktopVideoControlsThemeData identity doesn't
  // change on every setState, which would cause VideoControlsThemeDataInjector
  // to deactivate mid-frame and abort Impeller on macOS.
  List<Widget> _controls = const [];

  _VideoState(Player player) : controller = VideoController(player);

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

    subtracks = widget.player.stream.tracks.listen((state) {
      setState(() {});
    });

    widget.player.stream.log.listen((log) => print(log));

    suberrors = widget.player.stream.error.listen((err) {
      errors.value = ds.Error.unknown(
        err,
        onTap: () => errors.value = ds.Error.zero,
        decoration: ds.ErrorDecorations.info,
      );
    });

    subplaying = widget.player.stream.playing.listen((playing) {
      if (playing) errors.value = ds.Error.zero;
    });
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    final defaults = ds.Defaults.of(context);
    final plist = internal.Playlist.of(context)!;
    _controls = [
      PlayerControlShuffle(plist.current),
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
      if (authn.developer(context).debug) PlayerControlStream(widget.player),
      if (authn.developer(context).debug) SizedBox.square(dimension: defaults.spacing),
      PlayerControlFullscreen(widget.player),
    ];
  }

  @override
  void dispose() {
    subtracks.cancel();
    suberrors.cancel();
    subplaying.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return debug.Lifecycle(
      message: "player lifecycle",
      Builder(
        builder: (context) {
            final theme = Theme.of(context);
            final defaults = ds.Defaults.of(context);
            final plist = internal.Playlist.of(context)!;
            final current = plist.known;
            final compact = defaults.isCompact;

            return FocusScope(
              node: _selffocus,
              child: Stack(
                fit: StackFit.passthrough,
                children: [
                  debug.Lifecycle(
                    message: "video controls theme",
                    MaterialDesktopVideoControlsTheme(
                      normal: MaterialDesktopVideoControlsThemeData(
                        modifyVolumeOnScroll: false,
                        bottomButtonBar: _controls,
                        keyboardShortcuts: {},
                      ),
                      fullscreen: MaterialDesktopVideoControlsThemeData(
                        modifyVolumeOnScroll: false,
                        bottomButtonBar: _controls,
                        keyboardShortcuts: {},
                      ),
                      child: ValueListenableBuilder<Widget>(
                        valueListenable: errors,
                        builder: (context, error, child) => ds.Loading(cause: error, child!),
                        child: Video(focusNode: widget.focus, controller: controller),
                      ),
                    ),
                  ),
                  ValueListenableBuilder<bool>(
                    valueListenable: widget.overlay,
                    builder: (context, overlay, child) => Visibility(
                      maintainState: true,
                      maintainFocusability: true,
                      visible: overlay,
                      child: child!,
                    ),
                    child: Column(
                      mainAxisSize: MainAxisSize.max,
                      mainAxisAlignment: MainAxisAlignment.start,
                      verticalDirection: compact ? VerticalDirection.up : VerticalDirection.down,
                      children: [
                        Expanded(
                          child: DecoratedBox(
                            decoration: BoxDecoration(
                              color: theme.scaffoldBackgroundColor.withValues(
                                // TODO: make opacity configurable because people have
                                // personal preferences....
                                // alpha: defaults.opaque.a,
                                alpha: 1.0,
                              ),
                            ),
                            child: widget.child,
                          ),
                        ),
                        Visibility(
                          visible: !uuidx.isMin(uuidx.fromString(current.id)),
                          child: PlayerControlResume(widget.player, current, widget.overlay),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            );
          },
        ),
      );
  }
}
