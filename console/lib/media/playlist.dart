import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:retrovibed/library.dart' hide Media;
import 'package:stream_transform/stream_transform.dart';
import 'package:media_kit/media_kit.dart';
import 'package:language_code/language_code.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/langcodex.dart' as langcodex;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/debug.dart' as debug;
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart' as api;
import 'play.queue.dart';

class Playlist extends StatefulWidget {
  static void _noop(
    BuildContext ctx,
    Duration pos,
    Duration dur,
    api.MediaSearchRequest q,
    String id,
  ) {}

  final Widget child;
  final void Function(
    BuildContext ctx,
    Duration position,
    Duration duration,
    api.MediaSearchRequest query,
    String id,
  )
  tracing;

  Playlist(this.child, {Key? key, this.tracing = _noop}) : super(key: key);

  static _PlaylistState? of(BuildContext context) {
    return context.findAncestorStateOfType<_PlaylistState>();
  }

  @override
  State<Playlist> createState() => _PlaylistState();

  static Future<void> file(BuildContext context, String path) {
    print("opening ${path}");
    return of(context)?.player.open(Media(path)).catchError((
          cause,
        ) {
          print("failed ${path} - ${cause}");
          return Future.error(cause);
        }) ??
        Future.value();
  }

  static Widget wrap(
    Widget Function(BuildContext context, _PlaylistState s) b,
  ) {
    return ds.build(
      (context) {
        final _PlaylistState? current = Playlist.of(context);
        // if we don't have a playlist ancestor thats a bug.
        assert(current != null);
        final s = current!;
        final defaults = ds.Defaults.of(context);
        return debug.Lifecycle(
          message: "playlist shortcuts",
          ds.Shortcuts(
            b(context, s),
            enabled: defaults.desktop,
            bindings: {
              const SingleActivator(LogicalKeyboardKey.escape): (
                const Text('play/pause'),
                () {
                  print("shortcut: play/pause");
                  s.player.playOrPause();
                  return KeyEventResult.handled;
                },
              ),
              const SingleActivator(LogicalKeyboardKey.arrowRight): (
                const Text('seek forward 10s'),
                () {
                  final pos = s.player.state.position + const Duration(seconds: 10);
                  print("shortcut: seek forward -> ${pos}");
                  s.player.seek(pos);
                  return KeyEventResult.handled;
                },
              ),
              const SingleActivator(LogicalKeyboardKey.arrowLeft): (
                const Text('seek backward 10s'),
                () {
                  final pos = s.player.state.position - const Duration(seconds: 10);
                  final clamped = pos < Duration.zero ? Duration.zero : pos;
                  print("shortcut: seek backward -> ${clamped}");
                  s.player.seek(clamped);
                  return KeyEventResult.handled;
                },
              ),
              const SingleActivator(control: true, LogicalKeyboardKey.arrowUp): (
                const Text('volume up'),
                () {
                  final vol = (s.player.state.volume + 5.0).clamp(0.0, 100.0);
                  print("shortcut: volume up -> ${vol}");
                  s.player.setVolume(vol);
                  return KeyEventResult.handled;
                },
              ),
              const SingleActivator(
                control: true,
                LogicalKeyboardKey.arrowDown,
              ): (
                const Text('volume down'),
                () {
                  final vol = (s.player.state.volume - 5.0).clamp(0.0, 100.0);
                  print("shortcut: volume down -> ${vol}");
                  s.player.setVolume(vol);
                  return KeyEventResult.handled;
                },
              ),
              const SingleActivator(control: true, LogicalKeyboardKey.keyM): (
                const Text('mute'),
                () {
                  final vol = s.player.state.volume > 0 ? 0.0 : 100.0;
                  print("shortcut: mute -> ${vol}");
                  s.player.setVolume(vol);
                  return KeyEventResult.handled;
                },
              ),
              const SingleActivator(control: true, LogicalKeyboardKey.keyN): (
                const Text('next'),
                () {
                  print("shortcut: next");
                  s.next();
                  return KeyEventResult.handled;
                },
              ),
              const SingleActivator(control: true, LogicalKeyboardKey.keyP): (
                const Text('previous'),
                () {
                  print("shortcut: previous");
                  s.previous();
                  return KeyEventResult.handled;
                },
              ),
              const SingleActivator(control: true, LogicalKeyboardKey.keyF): (
                const Text('fullscreen'),
                () {
                  print("shortcut: fullscreen");
                  ds.Full.of(context)?.toggle();
                  return KeyEventResult.handled;
                },
              ),
            },
          ),
        );
      },
    );
  }
}

class _PlaylistState extends State<Playlist> {
  final PlayQueue _queue = PlayQueue();
  final Player player = Player();
  final TextEditingController controller = TextEditingController();
  final FocusNode playerfocus = FocusNode(
    debugLabel: 'playlist.player.focus',
    onKeyEvent: (node, event) {
      return KeyEventResult.ignored;
    },
  );
  final FocusNode searchfocus = FocusNode(debugLabel: 'playlist.search.focus');
  final FocusScopeNode _selffocus = FocusScopeNode(
    debugLabel: 'playlist.focus',
  );
  ValueNotifier<api.MediaSearchResponse> search = ValueNotifier(
    api.media.response(
      next: api.media.request(limit: 32, mimetypes: mimex.of(mimex.icoaudio)),
    ),
  );

  Known get known => _queue.current.value.known;
  ValueNotifier<PlayableMedia?> get current => _queue.current;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  @override
  void initState() {
    super.initState();
    player.stream.tracks.listen((track) {
      final current = LanguageCode.code;
      final audio = track.audio.firstWhere((t) {
        return langcodex.match(current, t.language ?? "");
      }, orElse: AudioTrack.auto);

      final subtitles =
          audio.id != AudioTrack.auto().id
              ? SubtitleTrack.no()
              : track.subtitle.firstWhere((t) {
                return langcodex.match(current, t.language ?? "");
              }, orElse: SubtitleTrack.auto);
      player.setAudioTrack(audio);
      player.setSubtitleTrack(subtitles);

      // print(
      //   "current: ${current} ${LanguageCode.locale} ${LanguageCode.code.englishName}",
      // );
      // print("audio: ${audio.id} ${audio.language} ${audio.title} -- ${audio}");
      // print(
      //   "subtitles: ${subtitles.id} ${subtitles.language} ${subtitles.title} -- ${subtitles}",
      // );
    });

    player.stream.position.throttle(const Duration(seconds: 3), trailing: true).listen((pos) {
      final id = known.id;
      if (id == uuidx.min()) return;
      if (search.value.next.query.trim().isEmpty) return;
      widget.tracing(context, pos, player.state.duration, search.value.next, id);
    });

    player.stream.error.listen((err) {
      print("error: ${err}");
    });

    player.stream.completed.listen((completed) {
      if (!completed) return;

      print(
        "advancing through playlist ${player.state.playlist.medias.length} ${player.state.playlist.medias}",
      );
      completed ? next() : player.pause();
    });

    player.stream.playing.listen((playing) {
      final defaults = ds.Defaults.of(context);
      final focus = playing || defaults.mobile ? playerfocus : searchfocus;
      print("playlist.playing: ${playing} - ${focus}");
      focus.requestFocus();
    });
  }

  void setPlaylist(api.MediaSearchRequest q, Stream<PlayableMedia> pl) {
    setState(() {
      _queue.reset(pl);
      search.value..next = q;
      controller.text = q.query;
    });
    next();
  }

  void next() {
    print(
      "next initiated: ${_queue.previous} | ${known.description} - ${_queue.currentStart} | ${_queue.upcoming}",
    );
    _advance().whenComplete(() {
      print(
        "next completed: ${_queue.previous} | ${known.description} | ${_queue.upcoming}",
      );
    });
  }

  void previous() {
    print(
      "previous initiated: ${_queue.previous} | ${known.description} | ${_queue.upcoming}",
    );
    _reverse().then((m) {
      print(
        "previous completed: ${_queue.previous} | ${known.description} | ${_queue.upcoming}",
      );
    });
  }

  Future<PlayableMedia?> _reverse() {
    return authn.bearer(authn.AuthzCache.meta(context)).then((auth) => _queue.reverse(auth, player)).then((m) {
      if (m != null) setState(() {});
      return m;
    });
  }

  Future<PlayableMedia?> _advance() {
    return authn.bearer(authn.AuthzCache.meta(context)).then((auth) => _queue.advance(auth, player)).then((m) {
      if (m != null) setState(() {});
      return m;
    });
  }

  @override
  void dispose() {
    super.dispose();
    player.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return FocusScope(
      // JAL: we believe this was added to deal with prior issues with mediakits state management nonsense.
      // key: ValueKey(_queue.current.id), // causes VideoControlsThemeDataInjector deactivated-ancestor crash in media_kit_video
      node: _selffocus,
      child: widget.child,
    );
  }
}
