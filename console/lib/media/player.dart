import 'dart:async';
import 'package:flutter/material.dart';
import 'package:media_kit/media_kit.dart'; // Provides [Player], [Media], [Playlist] etc.
import 'package:media_kit_video/media_kit_video.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/mimex.dart' as mimex;
import './playlist.dart' as internal;
import './player.settings.dart';

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
        final plist = internal.Playlist.of(context);
        final title = plist?.current.description ?? "";
        final compact = defaults.isCompact;
        final controls = [
          ds.build(
            (context) => IconButton(
              onPressed: () {
                internal.Playlist.of(context)?.previous();
              },
              icon: Icon(Icons.skip_previous_rounded),
            ),
          ),
          MaterialPlayOrPauseButton(),
          Builder(
            builder:
                (context) => IconButton(
                  onPressed: () {
                    internal.Playlist.of(context)?.next();
                  },
                  icon: Icon(Icons.skip_next_rounded),
                ),
          ),
          MaterialDesktopVolumeButton(),
          SizedBox.square(dimension: defaults.spacing),
          Expanded(
            child: StreamBuilder(
              stream: widget.player.stream.track,
              builder: (context, snapshot) {
                final plist = internal.Playlist.of(context);
                final title = plist?.current.description ?? "";
                return Text(
                  title,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                );
              },
            ),
          ),
          SizedBox.square(dimension: defaults.spacing),
          MaterialPositionIndicator(),
          SizedBox.square(dimension: defaults.spacing),
          ds.Help(
            ds.build(
              (context) => IconButton(
                onPressed: () {
                  ds.modals.push(
                    context,
                    PlayerSettings(
                      constraints: BoxConstraints(maxWidth: 256),
                      padding: defaults.padding,
                      current: widget.player,
                    ),
                  );
                },
                icon: Icon(Icons.tune),
              ),
            ),
            ds.Hint(const Text("open settings to select audio, video, or subtitle tracks")),
          ),
          SizedBox.square(dimension: defaults.spacing),
          ds.build(
            (context) {
              // havent figured out how to both filter and allow folder navigation on linux GTK.
              const mimetypes = <String>[
                // ...mimex.folders,
                ...mimex.audios,
                ...mimex.videos,
                "application/x-iso9660-image", // sometimes bluray/dvd are iso images.
              ];
              final playfile = (ds.FilesEvent evt, {ValueNotifier<int>? progress}) {
                print("play file checkpoint 0 ${evt.files.length}");
                if (evt.files.isEmpty) return Future.value(ds.NullWidget);
                final file = evt.files.firstWhere((v) {
                  print("play file checkpoint 1 ${v.mimeType}");
                  return mimetypes.any((x) => x == v.mimeType);
                }, orElse: () => evt.files.first);

                print("play file checkpoint 2 ${file.name}");
                return internal.Playlist.file(context, "file://${file.path}").then((v) => ds.NullWidget);
              };

              return ds.FileDropWell.icon(
                playfile,
                mimetypes: mimetypes,
                icon: Icons.video_collection_rounded,
                tooltip: "play a local media file",
                help: ds.Hint(const Text("play a local media file")),
              );
            },
          ),
          SizedBox.square(dimension: defaults.spacing),
          IconButton(
            onPressed: () => ds.Full.toggle(context),
            icon: Icon(
              ds.Full.nochrome(context) ? Icons.fit_screen : Icons.crop_free,
            ),
            tooltip: ds.Full.nochrome(context) ? 'Exit Full Screen' : 'Full Screen',
          ),
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
                      visible: !uuidx.isMin(uuidx.fromString(plist?.current.id ?? uuidx.min())),
                      child: DecoratedBox(
                        decoration: BoxDecoration(
                          color: theme.scaffoldBackgroundColor,
                        ),
                        child: IconButton(
                          onPressed: () {
                            widget.player.play();
                          },
                          icon: Row(
                            spacing: defaults.spacing,
                            children: [
                              Icon(Icons.play_circle_outline_rounded),
                              Text("Resume ${title}"),
                            ],
                          ),
                        ),
                      ),
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
