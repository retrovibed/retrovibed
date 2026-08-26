import 'dart:async';
import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/authz.dart' as authz;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/media/play.queue.dart' as playqueue;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/remote.dart';
import 'package:retrovibed/retrovibed.dart' as retro;
import 'api.dart' as remote;

// Mounted as an ancestor of media.Playlist's UI-consuming subtree, drives the
// same public Playlist surface the on-screen transport buttons use
// (playlist.next()/.previous(), playlist.player, playlist.queue) rather than
// reaching into Playlist's private state directly.
class RemoteControlListener extends StatefulWidget {
  final Widget child;
  final Future<remote.RemoteControlSocket> Function({List<httpx.Option> options}) connect;
  final meta.Daemon Function() localDevice;
  const RemoteControlListener(
    this.child, {
    super.key,
    this.connect = remote.remotecontrol.listen,
    this.localDevice = retro.local_device,
  });

  @override
  State<RemoteControlListener> createState() => _State();
}

class _State extends State<RemoteControlListener> {
  String _sid = uuidx.min();
  remote.RemoteControlSocket _socket = remote.RemoteControlSocket.noop;
  StreamSubscription<remote.Stream>? _rcSubscription;
  StreamSubscription<bool>? _playingSubscription;
  StreamSubscription<double>? _volumeSubscription;
  ValueNotifier<meta.Daemon> _library = ValueNotifier(meta.Daemon());
  ValueNotifier<authz.Bearer<meta.Token>> _authz = ValueNotifier(authz.Bearer(meta.Token(), ""));
  playqueue.PlayQueue _queue = playqueue.PlayQueue();
  // the underlying/desired volume + mute state, independent of what the
  // player's actually outputting - what Sync.volume/Sync.muted report,
  // since muting doesn't lose the level the user actually wants.
  remote.Sync _previous = remote.Sync(volume: 0.0, muted: false);
  // suppresses _volumeSubscription reacting to our own mute-driven
  // setVolume(0.0)/restore calls, so they don't clobber _previous.volume -
  // same pattern as playlist.dart's _transitioning guard.
  bool _muting = false;

  // proactively (not just in reply to a request) reports the listener's
  // current library, playback queue, and a token valid against this device -
  // triggered by track/library/queue changes and by didChangeDependencies
  // whenever the token cache actually refreshes.
  Future<void> _echoSync() async {
    final queue = _queue;
    final library = _library.value;
    library..hostname = meta.daemons.isLocalDevice(library) ? widget.localDevice().hostname : library.hostname;
    final cached = authn.AuthzCache.meta(context);
    // forces a refresh if expired, so token is never blank
    return cached.auto().then((bearer) {
      if (!mounted) return;
      final sync = remote.messages.sync(
        library: library,
        token: httpx.bearer(bearer.bearer),
        expiration: bearer.token.expires,
        capacity: queue.capacity,
        current: queue.current.value?.current,
        queue: queue.queued.map((m) => m.current).toList(),
        volume: _previous.volume,
        muted: _previous.muted,
        fullscreen: ds.Full.nochrome(context),
      );
      _socket.send(
        sync,
      );
      setState(() {
        _sid = sync.sid;
      });
    });
  }

  void _rcReconnect() {
    if (!mounted) return;
    Future.delayed(const Duration(seconds: 2), _rcConnect);
  }

  // listen is process-local and never depends on the profile being logged
  // in, so it connects unconditionally and reconnects forever on close.
  void _rcConnect() {
    widget
        .connect(options: [httpx.Request.bearer(() => Future.value(retro.remote_control_listen_token()))])
        .then((socket) {
          final c = Completer();

          _rcSubscription = socket.messages.listen(
            (msg) => _applyRemoteCommand(msg).whenComplete(_echoSync),
            cancelOnError: true,
            onError: c.completeError,
            onDone: c.complete,
          );

          setState(() {
            _socket = socket;
          });

          return c.future;
        })
        .catchError((cause) {
          debugPrint("remote control listen socket failed: ${cause}");
        })
        .whenComplete(() {
          _socket = RemoteControlSocket.noop;
          _rcSubscription = null;
          _rcReconnect();
        });
  }

  Future<void> _applyRemoteCommand(remote.Stream msg) {
    print("executing remote command ${msg}");

    switch (msg.whichCommand()) {
      case remote.Stream_Command.queue:
        return _applyQueue(msg);
      case remote.Stream_Command.dequeue:
        return _applyDequeue(msg);
      case remote.Stream_Command.playpause:
        return _applyPlayPause(msg);
      case remote.Stream_Command.seek:
        return _applySeek(msg);
      case remote.Stream_Command.volume:
        return _applyVolume(msg);
      case remote.Stream_Command.fullscreen:
        return _applyFullscreen(msg);
      case remote.Stream_Command.mute:
        return _applyMute(msg);
      case remote.Stream_Command.sync:
      case remote.Stream_Command.notSet:
        return Future.value();
    }
  }

  Future<void> _applyQueue(remote.Stream msg) async {
    final playlist = media.Playlist.of(context);
    if (playlist == null) return;
    playlist.maybeNext(playqueue.PlayableMedia(msg.queue.media));
  }

  Future<void> _applyDequeue(remote.Stream msg) async {
    final playlist = media.Playlist.of(context);
    if (playlist == null) return;
    playlist.queue.remove(msg.dequeue.id);
  }

  Future<void> _applyPlayPause(remote.Stream msg) async {
    if (_sid.compareTo(msg.sid) >= 0) return;
    final playlist = media.Playlist.of(context);
    if (playlist == null) return;
    setState(() => _sid = msg.sid);
    msg.playpause.paused ? playlist.player.pause() : playlist.player.play();
  }

  Future<void> _applySeek(remote.Stream msg) async {
    if (_sid.compareTo(msg.sid) >= 0) return;
    final playlist = media.Playlist.of(context);
    if (playlist == null) return;
    setState(() => _sid = msg.sid);
    final offset = msg.seek.offset;
    if (offset == remote.SeekOffset.next) {
      playlist.next();
    } else if (offset == remote.SeekOffset.previous) {
      playlist.previous();
    } else {
      playlist.player.seek(playlist.player.state.position + Duration(milliseconds: offset));
    }
  }

  // relative adjustment against the remembered (not necessarily currently
  // audible, if muted) level - touching volume is assumed to mean "I want
  // sound," so this always unmutes.
  Future<void> _applyVolume(remote.Stream msg) async {
    if (_sid.compareTo(msg.sid) >= 0) return;
    final playlist = media.Playlist.of(context);
    if (playlist == null) return;
    setState(() {
      _sid = msg.sid;
      _previous = remote.Sync(
        volume: (_previous.volume + msg.volume.offset).clamp(0.0, 100.0).toDouble(),
        muted: false,
      );
    });
    await playlist.player.setVolume(_previous.volume);
  }

  // toggles mute without losing the level to restore to - unlike the
  // ctrl+M keyboard shortcut (playlist.dart), which just flips between 0
  // and 100 with no memory of what it was before.
  Future<void> _applyMute(remote.Stream msg) async {
    if (_sid.compareTo(msg.sid) >= 0) return;
    final playlist = media.Playlist.of(context);
    if (playlist == null) return;
    setState(() {
      _sid = msg.sid;
      _previous = remote.Sync(volume: _previous.volume, muted: !_previous.muted);
    });
    _muting = true;
    await playlist.player.setVolume(_previous.muted ? 0.0 : _previous.volume);
    _muting = false;
  }

  Future<void> _applyFullscreen(remote.Stream msg) async {
    if (_sid.compareTo(msg.sid) >= 0) return;
    setState(() => _sid = msg.sid);
    ds.Full.toggle(context);
  }

  @override
  void initState() {
    super.initState();

    final _player = media.Playlist.of(context)?.player;

    _previous = remote.Sync(volume: _player?.state.volume ?? 0.0, muted: false);
    // echo local state back over the listen socket so any /rc/connect
    // observers can see what this device is doing.
    _playingSubscription = _player?.stream.playing.listen((playing) {
      _socket.send(remote.messages.playpause(!playing));
    });

    // local (non-remote, e.g. keyboard-shortcut) volume changes still need
    // tracking into _previous.volume and echoing - except our own
    // mute-driven writes, which _muting suppresses so they don't clobber
    // the level being remembered for restore-on-unmute.
    _volumeSubscription = _player?.stream.volume.listen((volume) {
      if (!_muting) {
        _previous = remote.Sync(volume: volume, muted: _previous.muted);
      }
      _echoSync();
    });

    _queue = media.Playlist.of(context)?.queue ?? _queue;
    _queue.current.addListener(_echoSync);
    _queue.revision.addListener(_echoSync);

    _library = meta.EndpointAuto.of(context)?.changed ?? _library;
    _library.addListener(_echoSync);

    _authz = authn.AuthzCache.of(context).changed;
    _authz.addListener(_echoSync);

    _rcConnect();
  }

  @override
  void dispose() {
    super.dispose();
    _queue.current.removeListener(_echoSync);
    _queue.revision.removeListener(_echoSync);
    _library.removeListener(_echoSync);
    _authz.removeListener(_echoSync);
    _playingSubscription?.cancel();
    _volumeSubscription?.cancel();
    _rcSubscription?.cancel();
    _socket.close();
  }

  @override
  Widget build(BuildContext context) => widget.child;
}
