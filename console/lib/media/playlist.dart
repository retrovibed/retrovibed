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
import 'play.queue.dart' as playqueue;

// The surface RemoteControlListener needs from a Playlist, extracted so it
// can depend on this instead of the real widget/native player directly -
// tests substitute a plain-Dart fake instead of mounting a real Playlist.
abstract class PlaylistControl {
  // no-op sentinel for "no Playlist ancestor" - same null-object pattern
  // remote.RemoteControlSocket.noop already establishes, so callers never
  // need to null-check.
  static final PlaylistControl zero = _ZeroPlaylistControl();

  playqueue.PlayQueue get queue;
  void maybeNext(playqueue.PlayableMedia m);
  void next();
  void previous();
  void playOrPause();
  Duration get position;
  void seek(Duration position);
  ValueNotifier<double> get volume;
  Future<void> setVolume(double volume);
  ValueNotifier<bool> get playing;
}

class _ZeroPlaylistControl implements PlaylistControl {
  @override
  final playqueue.PlayQueue queue = playqueue.PlayQueue();
  @override
  final ValueNotifier<double> volume = ValueNotifier(0.0);
  @override
  final ValueNotifier<bool> playing = ValueNotifier(false);
  @override
  Duration get position => Duration.zero;
  @override
  void maybeNext(playqueue.PlayableMedia m) {}
  @override
  void next() {}
  @override
  void previous() {}
  @override
  void playOrPause() {}
  @override
  void seek(Duration position) {}
  @override
  Future<void> setVolume(double volume) async {}
}

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
            ValueListenableBuilder<playqueue.PlayableMedia?>(
              valueListenable: s.current,
              builder: (context, _, __) => b(context, s),
            ),
            enabled: defaults.desktop,
            bindings: {
              const SingleActivator(LogicalKeyboardKey.escape): (
                const Text('play/pause'),
                () {
                  // when we're dealing with audio content opening the search overlay
                  // shouldnt pause playback.
                  if (mimex.isAudio(s.current.value?.current.mimetype ?? mimex.binary)) {
                    print("shortcut: play/pause (esc) -  audio");
                    s.overlay.value = !s.overlay.value;
                    s.player.play();
                  } else {
                    print("shortcut: play/pause (esc) - non-audio");
                    s.player.playOrPause();
                  }

                  return KeyEventResult.handled;
                },
              ),
              const SingleActivator(LogicalKeyboardKey.pause): (
                const Text('pause'),
                () {
                  print("shortcut: pause");
                  s.player.playOrPause();
                  return KeyEventResult.handled;
                },
              ),
              const SingleActivator(control: true, LogicalKeyboardKey.arrowRight): (
                const Text('seek forward 10s'),
                () {
                  final pos = s.player.state.position + const Duration(seconds: 10);
                  print("shortcut: seek forward -> ${pos}");
                  s.player.seek(pos);
                  return KeyEventResult.handled;
                },
              ),
              const SingleActivator(control: true, LogicalKeyboardKey.arrowLeft): (
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

class _PlaylistState extends State<Playlist> implements PlaylistControl {
  final playqueue.PlayQueue _queue = playqueue.PlayQueue();
  final Player player = Player(
    configuration: PlayerConfiguration(
      title: 'retrovibed',
    ),
  );
  // PlaylistControl's view of volume/playing - kept in sync with the real
  // player via the subscriptions set up in initState, exposed as
  // ValueNotifiers (read via .value, observe via .addListener) to match how
  // current/queue.revision are already exposed below.
  @override
  final ValueNotifier<double> volume = ValueNotifier(0.0);
  @override
  final ValueNotifier<bool> playing = ValueNotifier(false);
  playqueue.RangeFn autoqueue = playqueue.search;
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
  ValueNotifier<api.MediaSearchState> search = ValueNotifier(
    api.MediaSearchState(
      next: api.media.request(limit: 32, mimetypes: mimex.of(mimex.icoaudio)),
    ),
  );
  final ValueNotifier<bool> overlay = ValueNotifier<bool>(true);
  // guards the player.stream.playing listener during _advance()/_reverse():
  // both can still fall back to Player.open(), which emits a transient
  // playing:false->true blip while (re)loading a track - suppress reacting
  // to that noise and resync once from the real, settled state afterward.
  bool _transitioning = false;

  Known get known => _queue.current.value.known;
  ValueNotifier<playqueue.PlayableMedia?> get current => _queue.current;
  @override
  playqueue.PlayQueue get queue => _queue;
  @override
  Duration get position => player.state.position;
  @override
  void seek(Duration position) => player.seek(position);
  @override
  void playOrPause() => player.playOrPause();
  @override
  Future<void> setVolume(double volume) => player.setVolume(volume);

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  @override
  void initState() {
    super.initState();

    if (player.platform is NativePlayer) {
      final native = player.platform as NativePlayer;
      // Instructs libmpv to drop late frames at VO layer to keep audio in sync
      native
          .setProperty('framedrop', 'vo')
          .then((_) {
            print("framedrop property set");
          })
          .catchError((cause) {
            print("framedrop property failed to set ${cause}");
          });
    }

    volume.value = player.state.volume;
    playing.value = player.state.playing;
    player.stream.volume.listen((v) => volume.value = v);
    player.stream.playing.listen((p) => playing.value = p);

    search.addListener(() {
      setState(() {
        autoqueue = switch (mimex.category(search.value.next.mimetypes)) {
          mimex.audio => playqueue.acoustic,
          _ => playqueue.search,
        };
      });
    });

    player.stream.tracks.listen((track) {
      final current = LanguageCode.code;
      final audio = track.audio.firstWhere((t) {
        return langcodex.match(current, t.language ?? "");
      }, orElse: AudioTrack.auto);

      final subtitles = audio.id != AudioTrack.auto().id
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

    player.stream.completed.listen((completed) {
      if (!completed) return;

      debugPrint(
        "advancing through playlist ${player.state.playlist.medias.length}",
      );
      completed ? next() : player.pause();
    });

    // _swap() (play.queue.dart) removes the blip for in-queue advances at the
    // source; _transitioning (set around _advance()/_reverse()) suppresses the
    // remaining open() call sites (reverse(), and the first track of a freshly
    // reset queue) where stream.playing still emits a transient false->true
    // pair while the track is being (re)loaded.
    player.stream.playing.listen((playing) {
      if (_transitioning) return;
      _applyPlaying(playing);
    });
  }

  void _applyPlaying(bool playing) {
    final defaults = ds.Defaults.of(context);
    final focus = playing || defaults.mobile ? playerfocus : searchfocus;
    debugPrint("playlist.playing: ${playing} - ${focus}");
    focus.requestFocus();
    overlay.value = !playing;
  }

  void setPlaylist(
    api.MediaSearchRequest q,
    api.Media current,
    playqueue.RangeFn autoqueue, {
    Duration pos = const Duration(milliseconds: 0),
  }) {
    setState(() {
      this.autoqueue = autoqueue;
      _queue.reset(
        autoqueue(
          q,
          _queue,
          options: () => [authn.request(authn.AuthzCache.meta(context))],
        ),
        current,
        pos: pos,
      );
      search.value..next = q;
      controller.text = q.query;
    });
    next();
  }

  // jump start playback if its not already running, otherwise just ignore.
  void maybeNext(playqueue.PlayableMedia m) {
    print(
      "maybe next initiated: ${player.state.playing} ${_queue.previous} | ${known.description} - ${_queue.pos} | ${_queue.upcoming}",
    );

    _queue.push(m);

    if (player.state.playing) return;

    _advance().whenComplete(() {
      print(
        "next completed: ${_queue.previous} | ${known.description} | ${_queue.upcoming}",
      );
    });
  }

  void next() {
    print(
      "next initiated: ${_queue.previous} | ${known.description} - ${_queue.pos} | ${_queue.upcoming}",
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

  Future<playqueue.PlayableMedia?> _reverse() {
    _transitioning = true;
    return authn
        .bearer(authn.AuthzCache.meta(context))
        .then((auth) => _queue.reverse(auth, player))
        .then((m) {
          if (m != null) setState(() {});
          return m;
        })
        .whenComplete(() {
          _transitioning = false;
          _applyPlaying(player.state.playing);
        });
  }

  Future<playqueue.PlayableMedia?> _advance() {
    _transitioning = true;
    return authn
        .bearer(authn.AuthzCache.meta(context))
        .then((auth) => _queue.advance(auth, player))
        .then((m) {
          if (m != null) setState(() {});
          return m;
        })
        .whenComplete(() {
          _transitioning = false;
          _applyPlaying(player.state.playing);
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
      node: _selffocus,
      child: widget.child,
    );
  }
}
