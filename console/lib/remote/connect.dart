import 'dart:async';

import 'package:flutter/material.dart';
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
import 'player.control.playpause.dart';
import 'player.control.sync.dart';
import 'playlist.current.dart';
import 'playlist.queue.dart';

// Public entrypoint: wraps _Connect in an authn.AuthedEndpoint so it can target a
// user-selected remote daemon (via its DaemonDropdown) with a matching
// scoped auth token, independent of the app-root EndpointAuto/AuthzCache.
class Connect extends StatelessWidget {
  final ValueNotifier<media.MediaSearchState> search;
  final Future<Stream<meta.Daemon>> Function({List<httpx.Option> options}) daemonDiscover;
  final Future<remote.RemoteControlSocket> Function({required String host, List<httpx.Option> options}) connect;

  const Connect({
    super.key,
    required this.search,
    this.daemonDiscover = meta.daemons.discover,
    this.connect = remote.remotecontrol.connect,
  });

  @override
  Widget build(BuildContext context) {
    return authn.AuthedEndpoint(
      _Connect(search: search, daemonDiscover: daemonDiscover, connect: connect),
    );
  }
}

class _Connect extends StatefulWidget {
  final ValueNotifier<media.MediaSearchState> search;
  final Future<Stream<meta.Daemon>> Function({List<httpx.Option> options}) daemonDiscover;
  final Future<remote.RemoteControlSocket> Function({required String host, List<httpx.Option> options}) connect;

  const _Connect({required this.search, required this.daemonDiscover, required this.connect});

  @override
  State<_Connect> createState() => _State();
}

class _State extends State<_Connect> with LoadingState {
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

  // long-lived autoqueue: tapping a search result seeds this from the tapped
  // item's search results, and _fillQueue tops it back up toward
  // _autoqueueTarget every time the daemon echoes a sync showing its
  // upcoming queue drained below that. Never null - starts as an empty
  // stream (mirrors PlayQueue._stream) so _fillQueue needs no null check.
  static const int _autoqueueTarget = 5;
  StreamIterator<playqueue.PlayableMedia> _autoqueue = StreamIterator(const Stream.empty());
  // serializes _fillQueue: sync echoes can arrive faster than a fill pass
  // completes.
  bool _filling = false;

  // temporary hack fix until we fix the listener.
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
  media.FnMediaSearch get _apisearch => media.media.searchendpoint(_hostname, _bearerOptions);
  media.FnMediaFind get _apirandom => media.media.randomendpoint(_hostname, _bearerOptions);

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
          setState(() {
            _autoqueue = StreamIterator(
              playqueue.range(s.next, anchor, search: _apisearch, random: _apirandom),
            );
          });
          await _fillQueue();
          setState(() => _focused = null);
        };
      default:
        return null;
    }
  }

  bool _casfilling(bool o) {
    if (o) return o;
    _filling = true;
    return false;
  }

  // tops the daemon's upcoming queue back up toward _autoqueueTarget by
  // pulling more results from _autoqueue - called after every tap and after
  // every sync that shows the daemon's queue has drained, which is what
  // makes this a long-lived queue rather than a one-shot burst.
  // sync echoes can arrive faster than a fill pass completes (each queued
  // item advancing the daemon triggers its own echo); a call that comes in
  // while one's already running is just dropped - the next sync will
  // trigger another pass anyway.
  Future<void> _fillQueue() async {
    // cache the queue while filling so a later _onPlay's new _autoqueue
    // doesn't get spliced into an in-flight pass - _socket is read fresh at
    // send time instead since it can be swapped/closed independently (e.g.
    // by _onEndpointChanged) partway through a pass.
    final needed = _autoqueueTarget - _latest.sync.queue.length;
    final queue = _autoqueue;
    if (needed <= 0) return;
    if (_casfilling(_filling)) return;
    try {
      for (var sent = 0; sent < needed && await queue.moveNext(); sent++) {
        _socket.send(remote.messages.queue(queue.current.current));
      }
    } finally {
      _filling = false;
    }
  }

  void _onEndpointChanged() {
    _socket.close();
    _autoqueue.cancel();
    _autoqueue = StreamIterator(const Stream.empty());
    setState(() {
      _socket = remote.RemoteControlSocket.noop;
      _latest = remote.Stream(sid: uuidx.min());
      _focused = PlaylistQueue(<media.Media>[], remote.RemoteControlSocket.noop, key: const ValueKey("queue"));
    });

    _connect();
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
                if (msg.whichCommand() != remote.Stream_Command.sync) return;
                if (msg.sid.compareTo(_latest.sid) <= 0) return;
                setState(() => _latest = msg);
                _fillQueue();
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
    final queue = PlaylistQueue(_latest.sync.queue, _socket, key: const ValueKey("queue"));
    return ds.Container(
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
                        Row(
                          spacing: defaults.spacing,
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            PlayerControlSeek.prev(socket: _socket),
                            PlayerControlSeek.backward(socket: _socket),
                            PlayerControlPlayPause(socket: _socket, status: _messages),
                            PlayerControlSeek.forward(socket: _socket),
                            PlayerControlSeek.next(socket: _socket),
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
                        PlaylistCurrent(_latest.sync.current),
                        Expanded(child: queue),
                      ],
                    ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
