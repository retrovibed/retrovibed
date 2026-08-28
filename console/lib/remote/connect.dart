import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/stateful.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/media/play.queue.dart' as playqueue;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/library/search.minimal.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'api.dart' as remote;
import 'player.control.seek.dart';
import 'player.control.fullscreen.dart';
import 'player.control.playpause.dart';
import 'player.control.sync.dart';
import 'player.control.volume.dart';
import 'playlist.current.dart';
import 'playlist.queue.dart';

// Public entrypoint: wraps _Connect in an authn.AuthedEndpoint so it can target a
// user-selected remote daemon (via its DaemonDropdown) with a matching
// scoped auth token, independent of the app-root EndpointAuto/AuthzCache.
class AutoConnect extends StatelessWidget {
  final ValueNotifier<media.MediaSearchState> search;
  final Future<Stream<meta.Daemon>> Function({List<httpx.Option> options}) daemonDiscover;
  final Future<remote.RemoteControlSocket> Function({required String host, List<httpx.Option> options}) connect;
  // endpoint factories for the remote daemon's search/random lookups,
  // defaulting to the real hardwired HTTP endpoints - lets tests seed the
  // autoqueue without a live network.
  final media.FnMediaSearch Function(String host, List<httpx.Option> options) apisearch;
  final media.FnMediaFind Function(String host, List<httpx.Option> options) apirandom;
  // long-lived autoqueue: tapping a search result seeds this from the tapped
  // item's search results, and _fillQueue tops it back up toward
  // _autoqueueTarget every time the daemon echoes a sync showing its
  // upcoming queue drained below that. Never null - starts as an empty
  // stream (mirrors PlayQueue._stream) so _fillQueue needs no null check.
  final int autoqueueTarget;

  const AutoConnect({
    super.key,
    required this.search,
    this.daemonDiscover = meta.daemons.discover,
    this.connect = remote.remotecontrol.connect,
    this.apisearch = media.media.searchendpoint,
    this.apirandom = media.media.randomendpoint,
    this.autoqueueTarget = Connect._autoqueueTargetDefault,
  });

  @override
  Widget build(BuildContext context) {
    return authn.AuthedEndpoint(
      Connect(
        search: search,
        daemonDiscover: daemonDiscover,
        connect: connect,
        apisearch: apisearch,
        apirandom: apirandom,
        autoqueueTarget: autoqueueTarget,
      ),
    );
  }
}

class Connect extends StatefulWidget {
  static const _autoqueueTargetDefault = 5;
  final ValueNotifier<media.MediaSearchState> search;
  final Future<Stream<meta.Daemon>> Function({List<httpx.Option> options}) daemonDiscover;
  final Future<remote.RemoteControlSocket> Function({required String host, List<httpx.Option> options}) connect;
  final media.FnMediaSearch Function(String host, List<httpx.Option> options) apisearch;
  final media.FnMediaFind Function(String host, List<httpx.Option> options) apirandom;
  final int autoqueueTarget;

  const Connect({
    required this.search,
    required this.daemonDiscover,
    required this.connect,
    required this.apisearch,
    required this.apirandom,
    this.autoqueueTarget = _autoqueueTargetDefault,
  });

  @override
  State<Connect> createState() => _State();
}

class _State extends State<Connect> with LoadingState {
  remote.RemoteControlSocket _socket = remote.RemoteControlSocket.noop;
  Stream<remote.Stream> _messages = Stream.empty();
  // nil-sid sentinel; unset oneof -> _latest.sync reads as a zero Sync.
  remote.Stream _latest = remote.Stream(sid: uuidx.min());
  // which widget occupies the focused slot below the transport controls,
  // defaulting to the pending-queue view. Kept in sync with the live
  // search/queue widgets at the top of build() (matched by key, since
  // search/queue are rebuilt fresh every build() - they close over live
  // _socket/_latest - so a stale reference here would stop updating).
  Widget? _focused;
  ValueNotifier<meta.Daemon> _endpoint = ValueNotifier(meta.Daemon());

  playqueue.SafeStreamIterator<playqueue.PlayableMedia> _autoqueue = playqueue.SafeStreamIterator(
    const Stream.empty(),
  );
  // serializes _fillQueue: sync echoes can arrive faster than a fill pass
  // completes.
  int _filling = 0;

  // completed iff autoplay is switched on; reset to a fresh pending
  // Completer when switched off - and doubles as the flag itself, so there's
  // no separate bool to drift out of sync. lets the already-running
  // _apirandom call (inside range()'s while(true) loop) just wait on this
  // instead of erroring/polling - the next toggle wakes it straight back up.
  Completer<void> _autoplay = Completer<void>();

  void _casautoplay(bool enabled) {
    setState(() {
      if (enabled) {
        _autoplay.complete();
      } else {
        _autoplay = Completer<void>();
      }
    });
  }

  // temporary hack fix until we fix the listener (and its deployed).
  // problem was when the device was serving its own content it sent localhost
  // as the library over the sync protocol. resulting in the wrong token used
  // against the wrong host. the longterm fix is to have golang ffi return the proper hostname.
  // see DaemonFromHost.
  String get _hostname =>
      meta.daemons.isLocalDevice(_latest.sync.library) ? _endpoint.value.hostname : _latest.sync.library.hostname;
  // resolved fresh on every access (not cached) since _latest changes over
  // the widget's lifetime and these must always target whatever
  // daemon/hostname/token is current, whether read from build() or later
  // from _onPlay/_fillQueue.
  List<httpx.Option> get _bearerOptions => [httpx.Request.bearer(() => Future.value(_latest.sync.token))];
  media.FnMediaSearch get _apisearch => widget.apisearch(_hostname, _bearerOptions);
  media.FnMediaFind get _apirandom => (req, {List<httpx.Option> options = const []}) async {
    if (!_autoplay.isCompleted) await _autoplay.future;
    return widget.apirandom(_hostname, _bearerOptions)(req, options: options);
  };

  // Mirrors media.PlayAction's shape but queues the media on the connected
  // remote daemon's playlist instead of this device's local Playlist, since
  // SearchMinimal here searches the app's own library to feed *that* daemon.
  Future<void> Function()? _onPlay(BuildContext context, media.Media current, media.MediaSearchResponse s) {
    switch (mimex.icon(current.mimetype)) {
      case mimex.icomovie:
      case mimex.icoaudio:
        return () async {
          await _autoqueue.cancel();
          final anchor = playqueue.PlayQueue()..current.value = playqueue.PlayableMedia(current);
          final queue = playqueue.SafeStreamIterator(
            playqueue.range(s.next, anchor, search: _apisearch, random: _apirandom),
          );
          setState(() {
            _autoqueue = queue;
          });
          await _fillQueue(queue);
          setState(() => _focused = null);
        };
      default:
        return null;
    }
  }

  bool _casfilling(int o) {
    final ret = _filling > 0;
    _filling += o;
    return ret;
  }

  // tops the daemon's upcoming queue back up toward _autoqueueTarget by
  // pulling more results from _autoqueue - called after every tap and after
  // every sync that shows the daemon's queue has drained, which is what
  // makes this a long-lived queue rather than a one-shot burst.
  // sync echoes can arrive faster than a fill pass completes (each queued
  // item advancing the daemon triggers its own echo); a call that comes in
  // while one's already running is just dropped - the next sync will
  // trigger another pass anyway.
  Future<void> _fillQueue(playqueue.SafeStreamIterator<playqueue.PlayableMedia> queue) async {
    final needed = widget.autoqueueTarget - _latest.sync.queue.length;
    if (needed <= 0) return;
    try {
      if (_casfilling(1)) return;
      if (await queue.moveNext()) {
        final mut = remote.syncmut.queue(queue.current.current);
        final update = _latest.deepCopy()..sync = mut(_latest.sync.deepCopy());
        setState(() {
          _latest = update;
          _socket.send(remote.messages.queue(queue.current.current));
        });

        ds.postframe(() => _fillQueue(_autoqueue));
      }
    } finally {
      _casfilling(-1);
    }
  }

  void _onEndpointChanged() {
    _socket.close();
    _autoqueue.cancel();
    _autoqueue = playqueue.SafeStreamIterator(const Stream.empty());
    setState(() {
      _socket = remote.RemoteControlSocket.noop;
      _latest = remote.Stream(sid: uuidx.min());
      _focused = PlaylistQueue(_latest.sync, remote.RemoteControlSocket.noop, key: const ValueKey("queue"));
    });

    _connect();
  }

  void _volumeAdjust(double delta) {
    _socket.send(remote.messages.volume(delta.round()));
  }

  void _volumeMute() {
    _socket.send(remote.messages.mute());
  }

  void _reconnect() {
    if (!mounted) return;
    print("reconnecting");
    setState(() => _socket = remote.RemoteControlSocket.noop);
    Future.delayed(const Duration(seconds: 2), _connect);
  }

  void _connect() {
    if (!mounted) return;
    if (_endpoint.value.hostname.isEmpty) return;
    if (meta.daemons.isLocalDevice(_endpoint.value)) {
      setState(() {
        loading = false;
        _socket = remote.RemoteControlSocket.disabled;
        cause = ds.Error.text(
          "you do not have permission to remotely control this device",
          decoration: ds.ErrorDecorations.info,
        );
      });
      return;
    }

    setState(() => loading = true);

    widget
        .connect(host: _endpoint.value.hostname, options: [authn.request(authn.AuthedEndpoint.token(context))])
        .then((socket) {
          final c = Completer();
          setState(() {
            loading = false;
            cause = ds.Error.zero;
            _socket = socket;
            _messages = socket.messages.asBroadcastStream();
            _messages.listen(
              (msg) {
                // print("sync ${msg.sid}");
                if (msg.whichCommand() != remote.Stream_Command.sync) return;
                print("sync received ${msg.sid}");
                // vid is a monotonic sequence number, unlike sid (a uuidv7)
                // whose ordering isn't guaranteed for two ids minted within
                // the same millisecond.
                if (msg.vid <= _latest.vid) return;
                print("sync accepted ${msg.sid} ${msg.sync.queue.length}");
                setState(() {
                  _latest = msg;
                });
                _fillQueue(_autoqueue);
              },
              cancelOnError: true,
              onError: c.completeError,
              onDone: c.complete,
            );
          });

          socket.send(remote.messages.sync());

          return c.future;
        })
        .then((_) => _reconnect())
        .catchError((error) {
          setState(() {
            loading = false;
            _socket = remote.RemoteControlSocket.noop;
            cause = ds.Error.unauthorized(
              error,
              message: const Text("remote control is disabled on this device"),
              decoration: ds.ErrorDecorations.info,
            );
          });
        }, test: httpx.ErrorsTest.forbidden)
        .catchError((error) {
          setState(() {
            loading = false;
            _socket = remote.RemoteControlSocket.noop;
            cause = ds.Error.unauthorized(
              error,
              message: const Text("you do not have permission to remotely control this device"),
              decoration: ds.ErrorDecorations.info,
            );
          });
        }, test: httpx.ErrorsTest.unauthorized)
        .catchError((error) {
          setState(() {
            loading = false;
            _socket = remote.RemoteControlSocket.noop;
            cause = ds.Errors.httpauto(error, onTap: resetCause);
          });
          _reconnect();
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((error) {
          setState(() {
            loading = false;
            _socket = remote.RemoteControlSocket.noop;
            cause = ds.Error.unknown(error, onTap: resetCause);
          });
          _reconnect();
        });
  }

  @override
  void initState() {
    super.initState();
    _autoplay.complete(); // enable autoplay by default.
    _endpoint = authn.AuthedEndpoint.daemon(context);
    _endpoint.addListener(_onEndpointChanged);
    ds.postframe(_connect);
  }

  @override
  void dispose() {
    super.dispose();
    _endpoint.removeListener(_onEndpointChanged);
    _socket.close();
    _autoqueue.cancel();
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final search = SearchMinimal(
      key: const ValueKey("search"),
      empty: ds.Empty,
      onPlay: _onPlay,
      apisearch: _apisearch,
    );
    final queue = PlaylistQueue(
      _latest.sync,
      _socket,
      key: const ValueKey("queue"),
      onChange: (mut) {
        final upd = _latest.deepCopy()..sync = mut(_latest.sync.deepCopy());
        setState(() {
          _latest = upd;
        });
      },
    );
    return ds.Shortcuts(
      bindings: {
        const SingleActivator(LogicalKeyboardKey.audioVolumeUp): (
          const Text("increase volume on the remote device"),
          () {
            _volumeAdjust(1);
            return KeyEventResult.handled;
          },
        ),
        const SingleActivator(LogicalKeyboardKey.audioVolumeDown): (
          const Text("decrease volume on the remote device"),
          () {
            _volumeAdjust(-1);
            return KeyEventResult.handled;
          },
        ),
        const SingleActivator(LogicalKeyboardKey.audioVolumeMute): (
          const Text("mute the remote device"),
          () {
            _volumeMute();
            return KeyEventResult.handled;
          },
        ),
      },
      ds.Container(
        Column(
          verticalDirection: defaults.isCompact ? VerticalDirection.up : VerticalDirection.down,
          children: [
            meta.DaemonDropdown(
              library: _endpoint,
              discover: widget.daemonDiscover,
              onSelect: meta.DaemonDropdown.local,
              remoteonly: true,
              readonly: true,
              trailing: [
                _focused == null
                    ? ds.Empty
                    : ds.LoadingIconButton.close(
                        onPressed: () {
                          setState(() => _focused = null);
                          return Future.value(null);
                        },
                      ),
              ],
            ),

            Expanded(
              child: ds.Container(
                padding: defaults.padding,
                ds.Loading(
                  cause: cause,
                  loading: _socket == remote.RemoteControlSocket.noop,
                  _focused ??
                      Column(
                        verticalDirection: defaults.isCompact ? VerticalDirection.up : VerticalDirection.down,
                        children: [
                          Wrap(
                            alignment: WrapAlignment.center,
                            spacing: defaults.spacing,
                            children: [
                              ds.LoadingIconButton(
                                icon: Icon(Icons.graphic_eq),
                                toggled: _autoplay.isCompleted,
                                onPressed: ds.LoadingIconButton.convert(() {
                                  _casautoplay(!_autoplay.isCompleted);
                                }),
                                tooltip: "enable autoqueue playback",
                                help: ds.Hint(
                                  const Text(
                                    "when the user has not queued any results, it will automatically queue up random content",
                                  ),
                                ),
                              ),
                              PlayerControlSeek.prev(socket: _socket),
                              PlayerControlSeek.backward(socket: _socket),
                              PlayerControlPlayPause(socket: _socket, paused: _latest.sync.paused),
                              PlayerControlSeek.forward(socket: _socket),
                              PlayerControlSeek.next(socket: _socket),
                              PlayerControlFullscreen(socket: _socket, current: _latest.sync.fullscreen),
                              ds.LoadingIconButton.search(
                                toggled: _focused?.key == search.key,
                                onPressed: ds.LoadingIconButton.convert(() {
                                  setState(() => _focused = _focused?.key == search.key ? ds.Empty : search);
                                }),
                                tooltip: "search the remote device's library",
                                help: ds.Hint(const Text("search the remote device's library to queue media on it")),
                              ),
                              if (authn.developer(context).debug) PlayerControlSync(socket: _socket),
                            ],
                          ),
                          PlayerControlVolume(socket: _socket, current: _latest.sync),
                          PlaylistCurrent(_latest.sync.current),
                          Expanded(child: queue),
                        ],
                      ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
