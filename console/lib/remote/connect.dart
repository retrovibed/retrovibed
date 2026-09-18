import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:stream_transform/stream_transform.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/stateful.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/media/play.queue.dart' as playqueue;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/library.dart' as lib;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/discovery.dart' as disc;
import 'api.dart' as remote;
import 'player.control.playback.dart';
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
  final Future<meta.DaemonSearchResponse> Function(meta.DaemonSearchRequest) daemonSearch;
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
    this.daemonSearch = meta.daemons.search,
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
        daemonSearch: daemonSearch,
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
  final Future<meta.DaemonSearchResponse> Function(meta.DaemonSearchRequest) daemonSearch;
  final Future<remote.RemoteControlSocket> Function({required String host, List<httpx.Option> options}) connect;
  final media.FnMediaSearch Function(String host, List<httpx.Option> options) apisearch;
  final media.FnMediaFind Function(String host, List<httpx.Option> options) apirandom;
  final lib.FnRecent apirecentlatest;
  final lib.FnRecentTombstone apirecenttombstone;
  final lib.FnRecentRecord apirecentrecord;
  final lib.FnRecentHistory apirecenthistory;
  // test seam: the window a test passes to tester.pump() to deterministically
  // control when the throttled recent-watch recording (see _State._record)
  // fires, mirroring RemoteControlListener's own position-echo throttle test.
  final Duration recentRecordThrottle;
  final int autoqueueTarget;

  const Connect({
    required this.search,
    required this.daemonDiscover,
    required this.daemonSearch,
    required this.connect,
    required this.apisearch,
    required this.apirandom,
    this.apirecentlatest = lib.recent.latest,
    this.apirecenttombstone = lib.recent.delete,
    this.apirecentrecord = lib.recent.record,
    this.apirecenthistory = lib.recent.history,
    this.recentRecordThrottle = const Duration(seconds: 3),
    this.autoqueueTarget = _autoqueueTargetDefault,
  });

  @override
  State<Connect> createState() => _State();
}

class _State extends State<Connect> with LoadingState {
  ValueNotifier<media.MediaSearchState> _search = ValueNotifier(
    media.MediaSearchState(
      next: media.media.request(limit: 32, mimetypes: mimex.of(mimex.icoaudio)),
    ),
  );
  remote.RemoteControlSocket _socket = remote.RemoteControlSocket.noop;
  Stream<remote.Stream> _messages = Stream.empty();
  // nil-sid sentinel; unset oneof -> _latest.sync reads as a zero Sync.
  remote.Stream _latest = remote.Stream(sid: uuidx.min());
  // which widget occupies the focused slot below the transport controls,
  // defaulting to the pending-queue view. Kept in sync with the live
  // search/queue widgets at the top of build()
  Widget? _focused;
  ValueNotifier<meta.Daemon> _endpoint = ValueNotifier(meta.Daemon());
  // identifies the current "pick session" - stable across reconnects and
  // repeated taps under the same query, but reminted (_onPlay) whenever the
  // active search query changes, so a Sync.current whose session_id matches
  // can be attributed back to _sessionQuery for both PlaylistCurrent's
  // highlight and _record's recent-watch write.
  String _sessionID = uuidx.v7();
  media.MediaSearchRequest _sessionQuery = media.MediaSearchRequest();

  // watch history heartbeat state, shared with media/playlist.dart via
  // lib.WatchHistoryTracker. Unlike the local player, Playback messages are
  // only sent while something is actually playing (see Playback's doc
  // comment), so ticks are never marked playing: false here.
  final lib.WatchHistoryTracker _history = lib.WatchHistoryTracker();

  playqueue.SafeStreamIterator<playqueue.PlayableMedia> _autoqueue = playqueue.SafeStreamIterator(
    const Stream.empty(),
  );
  // serializes _fillQueue: sync echoes can arrive faster than a fill pass
  // completes.
  int _filling = 0;

  // media dispatched via _fillQueue for the current autoqueue session that
  // no real daemon echo has confirmed yet - removed once a confirmed sync's
  // queue shows that media's id. Needed because _latest gets wholesale
  // overwritten by any newer-vid sync (see the message listener in
  // _connect), which would otherwise silently erase _fillQueue's own
  // optimistic bump to _latest.sync.queue before the daemon's real echo
  // catches up, reopening the "needed" gate and causing repeated over-fills.
  final Set<media.Media> _unconfirmed = {};

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
  media.FnMediaSearch get _apisearch =>
      widget.apisearch(_hostname, [httpx.Request.bearer(() => Future.value(_latest.sync.token))]);
  media.FnMediaFind get _apirandom => (req, {List<httpx.Option> options = const []}) async {
    if (!_autoplay.isCompleted) await _autoplay.future;
    return widget.apirandom(_hostname, [httpx.Request.bearer(() => Future.value(_latest.sync.token))])(
      req,
      options: options,
    );
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
          final anchor = playqueue.PlayQueue()
            ..current.value = playqueue.PlayableMedia(current, profileId: uuidx.min(), sessionId: uuidx.min());
          final queue = playqueue.SafeStreamIterator(
            playqueue.range(s.next, anchor, search: _apisearch, random: _apirandom),
          );
          setState(() {
            if (s.next.query != _sessionQuery.query) {
              _sessionID = uuidx.v7();
              _sessionQuery = s.next;
            }
            _autoqueue = queue;
            _unconfirmed.clear();
          });
          await _fillQueue(queue);
          setState(() => _focused = null);
        };
      default:
        return null;
    }
  }

  Future<void> _onRecentTap(BuildContext context, media.RecentRecordRequest item) async {
    return (_onPlay(context, item.media, media.MediaSearchResponse(next: item.query)) ?? () async {})();
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
    final needed = widget.autoqueueTarget - _latest.sync.queue.length - _unconfirmed.length;
    if (needed <= 0) return;
    try {
      if (_casfilling(1)) return;
      if (await queue.moveNext()) {
        final next = queue.current.current;
        final mut = remote.syncmut.queue(next, sessionId: _sessionID);
        final update = _latest.deepCopy()..sync = mut(_latest.sync.deepCopy());
        setState(() {
          _latest = update;
          _unconfirmed.add(next);
          _socket.send(remote.messages.queue(next, sessionId: _sessionID));
        });

        if (needed > 1) {
          ds.postframe(() => _socket.send(remote.messages.sync(sessionId: _sessionID)));
        }
      }
    } finally {
      _casfilling(-1);
    }
  }

  // records a recent-watch event for the daemon's current pick, but only
  // when this session is the one that enqueued it (msg.sessionId ==
  // _sessionID - the same comparison PlaylistCurrent uses to highlight it)
  // and a non-empty query is known for it. Called off a throttled stream
  // (see _connect), so this itself fires on every qualifying tick - no
  // dedupe/staleness bookkeeping needed here.
  void _record(remote.Stream msg) {
    if (!mounted) return;
    if (!msg.playback.media.hasId()) return;
    if (msg.sessionId != _sessionID) return;
    if (_sessionQuery.query.trim().isEmpty) return;
    final playback = msg.playback;
    widget
        .apirecentrecord(
          lib.RecentRecordRequest(
            media: playback.media,
            query: _sessionQuery,
            mimetype: mimex.category(_sessionQuery.mimetypes),
            position: playback.position,
            duration: playback.duration,
          ),
          host: _endpoint.value.hostname,
          options: [authn.request(authn.AuthedEndpoint.token(context))],
        )
        .then((v) {})
        .catchError((cause) {
          print("failed to record remote watch event - $cause");
        })
        .ignore();

    final record = _history.tick(playback.media.id);

    widget
        .apirecenthistory(
          lib.WatchHistoryRecordRequest(record: record),
          host: _endpoint.value.hostname,
          options: [authn.request(authn.AuthedEndpoint.token(context))],
        )
        .then((v) {
          print("recorded remote watch history ${_history.id}/${_history.mediaId}/${_history.watched}");
        })
        .catchError((cause) {
          print("failed to record remote watch history - $cause");
        })
        .ignore();
  }

  void _onEndpointChanged() {
    _socket.close();
    _autoqueue.cancel();
    _autoqueue = playqueue.SafeStreamIterator(const Stream.empty());
    setState(() {
      _socket = remote.RemoteControlSocket.noop;
      _latest = remote.Stream(sid: uuidx.min());
      _unconfirmed.clear();
      _focused = PlaylistQueue(
        _latest.sync,
        remote.RemoteControlSocket.noop,
        sessionId: _sessionID,
        key: const ValueKey("queue"),
      );
    });

    _connect();
  }

  void _volumeAdjust(double delta) {
    _socket.send(remote.messages.volume(delta.round(), sessionId: _sessionID));
  }

  void _volumeMute() {
    _socket.send(remote.messages.mute(sessionId: _sessionID));
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
        cause = SizedBox.expand(
          child: ds.Container(
            meta.DaemonList(
              search: widget.daemonSearch,
              remoteonly: true,
              decoration: const InputDecoration(hintText: "search for device to control"),
              onSelect: (v) => meta.DaemonDropdown.local(context, v).then((v) {
                _endpoint.value = v;
                return v;
              }),
            ),
          ),
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
                switch (msg.whichCommand()) {
                  case remote.Stream_Command.sync:
                    // print("sync ${msg.sid}");
                    print("sync received ${msg.sid} ${_latest.sync.queue.length}");
                    // vid is a monotonic sequence number, unlike sid (a uuidv7)
                    // whose ordering isn't guaranteed for two ids minted within
                    // the same millisecond.
                    if (msg.vid <= _latest.vid) return;
                    print("sync accepted ${msg.sid} ${msg.sync.queue.length}");
                    setState(() {
                      _unconfirmed.removeWhere((m) => msg.sync.queue.any((s) => s.asMedia.id == m.id));
                      _latest = msg;
                    });
                    _fillQueue(_autoqueue);
                    break;
                  case remote.Stream_Command.playback:
                    // a standalone playback frame reports live
                    // position/duration for _latest.sync.current - pair it
                    // onto _latest.sync.playback in place rather than
                    // waiting on the next full sync echo.
                    setState(() {
                      _latest = _latest.deepCopy()..sync = (_latest.sync.deepCopy()..playback = msg.playback);
                    });
                    break;
                  default:
                    break;
                }
              },
              cancelOnError: true,
              onError: c.completeError,
              onDone: c.complete,
            );
            _messages
                .where((msg) => msg.whichCommand() == remote.Stream_Command.playback)
                .throttle(widget.recentRecordThrottle, trailing: true)
                .listen(_record);
          });

          socket.send(remote.messages.sync(sessionId: _sessionID));

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
            cause = ds.Errors.httpauto(error, onTap: reseterr);
          });
          _reconnect();
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((error) {
          setState(() {
            loading = false;
            _socket = remote.RemoteControlSocket.noop;
            cause = ds.Error.unknown(error, onTap: reseterr);
          });
          _reconnect();
        });
  }

  @override
  void initState() {
    super.initState();
    _search = widget.search;
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
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);
    final search = lib.SearchMinimal(
      key: const ValueKey("search"),
      empty: ds.Empty,
      onPlay: _onPlay,
      apisearch: _apisearch,
      search: _search,
    );

    final queue = ValueListenableBuilder<media.MediaSearchState>(
      valueListenable: _search,
      builder: (context, state, _) {
        final category = mimex.category(state.next.mimetypes);
        return PlaylistQueue(
          _latest.sync,
          _socket,
          key: const ValueKey("queue"),
          sessionId: _sessionID,
          onChange: (mut) {
            final upd = _latest.deepCopy()..sync = mut(_latest.sync.deepCopy());
            setState(() {
              _latest = upd;
            });
          },
          empty: ds.Loading(
            loading: _latest.sync.token.isEmpty,
            maintainState: false,
            maintainAnimation: false,
            maintainSize: false,
            Column(
              verticalDirection: defaults.isCompact ? VerticalDirection.up : VerticalDirection.down,
              children: [
                ds.Heading(Text("Continue"), padding: defaults.padding / 2),
                Expanded(
                  child: disc.RecentList(
                    margin: defaults.margin.copyWith(left: 0, right: 0),
                    padding: defaults.padding / 2,
                    category,
                    latest: widget.apirecentlatest,
                    host: _latest.sync.library.hostname,
                    authz: httpx.Request.bearer(() => Future.value(_latest.sync.token)),
                    onTap: _onRecentTap,
                  ),
                ),
              ],
            ),
          ),
        );
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
        padding: defaults.padding / 4,
        Column(
          spacing: defaults.spacing / 2,
          verticalDirection: defaults.isCompact ? VerticalDirection.up : VerticalDirection.down,
          children: [
            meta.DaemonDropdown(
              padding: EdgeInsets.zero,
              library: _endpoint,
              discover: widget.daemonDiscover,
              onSelect: meta.DaemonDropdown.local,
              remoteonly: true,
              readonly: true,
              leading: [
                ds.CompactingMenu.pinned(
                  ValueListenableBuilder<media.MediaSearchState>(
                    valueListenable: _search,
                    builder: (context, state, _) => lib.DropdownUpload(
                      icon: lib.SearchMimetypeDropdown.icon(mimex.checksum(state.next.mimetypes)),
                      help: ds.HelpScope.None,
                      items: lib.SearchMimetypeDropdown.menuItems(_search.value, (upd) {
                        _search.value = upd;
                      }),
                    ),
                  ),
                ),
              ],
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
              child: ds.Loading(
                cause: cause,
                loading: _socket == remote.RemoteControlSocket.noop,
                _focused ??
                    ds.Container(
                      padding: defaults.padding / 2,
                      decoration: BoxDecoration(color: theme.colorScheme.surfaceContainerLow),
                      Column(
                        verticalDirection: defaults.isCompact ? VerticalDirection.up : VerticalDirection.down,
                        spacing: defaults.spacing / 2,
                        children: [
                          PlayerControlPlayback(socket: _socket, sessionId: _sessionID, current: _latest.sync),
                          ds.Heading(
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
                                PlayerControlVolume(socket: _socket, sessionId: _sessionID, current: _latest.sync),
                                PlayerControlSeek.prev(socket: _socket, sessionId: _sessionID),
                                PlayerControlSeek.backward(socket: _socket, sessionId: _sessionID),
                                PlayerControlPlayPause(
                                  socket: _socket,
                                  sessionId: _sessionID,
                                  paused: _latest.sync.paused,
                                ),
                                PlayerControlSeek.forward(socket: _socket, sessionId: _sessionID),
                                PlayerControlSeek.next(socket: _socket, sessionId: _sessionID),
                                PlayerControlFullscreen(
                                  socket: _socket,
                                  sessionId: _sessionID,
                                  current: _latest.sync.fullscreen,
                                ),
                                ds.LoadingIconButton.search(
                                  toggled: _focused?.key == search.key,
                                  onPressed: ds.LoadingIconButton.convert(() {
                                    setState(() => _focused = _focused?.key == search.key ? ds.Empty : search);
                                  }),
                                  tooltip: "search the remote device's library",
                                  help: ds.Hint(
                                    const Text("search the remote device's library to queue media on it"),
                                  ),
                                ),
                                if (authn.developer(context).debug)
                                  PlayerControlSync(socket: _socket, sessionId: _sessionID),
                              ],
                            ),
                          ),
                          if (_latest.sync.current.asMedia.hasId())
                            PlaylistCurrent(_latest.sync.current, sessionId: _sessionID),
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
