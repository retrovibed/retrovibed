import 'dart:async';
import 'package:fixnum/fixnum.dart' as fixnum;
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

media.PlaylistControl _resolvePlaylist(BuildContext context) =>
    media.Playlist.of(context) ?? media.PlaylistControl.zero;

// Mounted as an ancestor of media.Playlist's UI-consuming subtree, drives the
// same public Playlist surface the on-screen transport buttons use
// (via the injectable playlist accessor's PlaylistControl) rather than
// reaching into Playlist's private state directly.
class RemoteControlListener extends StatefulWidget {
  final Widget child;
  final Future<remote.RemoteControlSocket> Function({List<httpx.Option> options}) connect;
  final meta.Daemon Function() localDevice;
  final media.PlaylistControl Function(BuildContext) playlist;
  const RemoteControlListener(
    this.child, {
    super.key,
    this.connect = remote.remotecontrol.listen,
    this.localDevice = retro.local_device,
    this.playlist = _resolvePlaylist,
  });

  @override
  State<RemoteControlListener> createState() => _State();
}

class _State extends State<RemoteControlListener> {
  String _sid = uuidx.min();
  // monotonic counter this listener stamps on every sync it sends - unlike
  // sid (a uuidv7), which only guarantees chronological ordering at
  // millisecond granularity and can sort either way for two ids minted
  // within the same millisecond.
  int _vid = 0;
  remote.RemoteControlSocket _socket = remote.RemoteControlSocket.noop;
  StreamSubscription<remote.Stream>? _rcSubscription;
  ValueNotifier<meta.Daemon> _library = ValueNotifier(meta.Daemon());
  ValueNotifier<authz.Bearer<meta.Token>> _authz = ValueNotifier(authz.Bearer(meta.Token(), ""));
  playqueue.PlayQueue _queue = playqueue.PlayQueue();
  late media.PlaylistControl _playlistControl;
  bool _muted = false;
  // captured only at the instant a mute directive arrives, not tracked
  // continuously - the level to restore to on unmute.
  double _savedVolume = 0.0;

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
        volume: _playlistControl.volume.value,
        muted: _muted,
        paused: !_playlistControl.playing.value,
        fullscreen: ds.Full.nochrome(context),
        vid: fixnum.Int64(++_vid),
      );
      _socket.send(
        sync,
      );
      setState(() {
        _sid = sync.sid;
      });
    });
  }

  // echoes local (non-remote, e.g. keyboard-shortcut) playing/volume changes
  // back over the listen socket so any /rc/connect observers can see what
  // this device is doing - via a full sync, same as every remote command.
  void _echo() {
    _echoSync();
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
      case remote.Stream_Command.pause:
        return _applyPause(msg);
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
    _playlistControl.maybeNext(playqueue.PlayableMedia(msg.queue.media));
  }

  Future<void> _applyDequeue(remote.Stream msg) async {
    _playlistControl.queue.remove(msg.dequeue.id);
  }

  // no payload - toggles between playing and paused, same shape as mute/fullscreen.
  Future<void> _applyPause(remote.Stream msg) async {
    if (_sid.compareTo(msg.sid) >= 0) return;
    setState(() => _sid = msg.sid);
    _playlistControl.playOrPause();
  }

  Future<void> _applySeek(remote.Stream msg) async {
    if (_sid.compareTo(msg.sid) >= 0) return;
    setState(() => _sid = msg.sid);
    final offset = msg.seek.offset;
    if (offset == remote.SeekOffset.next) {
      _playlistControl.next();
    } else if (offset == remote.SeekOffset.previous) {
      _playlistControl.previous();
    } else {
      _playlistControl.seek(_playlistControl.position + Duration(milliseconds: offset));
    }
  }

  // relative adjustment against the remembered (not necessarily currently
  // audible, if muted) level - touching volume is assumed to mean "I want
  // sound," so this always unmutes.
  Future<void> _applyVolume(remote.Stream msg) async {
    if (_sid.compareTo(msg.sid) >= 0) return;
    setState(() => _sid = msg.sid);
    final base = _muted ? _savedVolume : _playlistControl.volume.value;
    final next = (base + msg.volume.offset).clamp(0.0, 100.0).toDouble();
    _muted = false;
    await _playlistControl.setVolume(next);
  }

  // toggles mute without losing the level to restore to - unlike the
  // ctrl+M keyboard shortcut (playlist.dart), which just flips between 0
  // and 100 with no memory of what it was before. the level to restore is
  // captured only at the instant of muting, not tracked continuously.
  Future<void> _applyMute(remote.Stream msg) async {
    if (_sid.compareTo(msg.sid) >= 0) return;
    setState(() => _sid = msg.sid);
    if (_muted) {
      _muted = false;
      await _playlistControl.setVolume(_savedVolume);
    } else {
      _savedVolume = _playlistControl.volume.value;
      _muted = true;
      await _playlistControl.setVolume(0.0);
    }
  }

  Future<void> _applyFullscreen(remote.Stream msg) async {
    if (_sid.compareTo(msg.sid) >= 0) return;
    setState(() => _sid = msg.sid);
    ds.Full.toggle(context);
  }

  @override
  void initState() {
    super.initState();

    _playlistControl = widget.playlist(context);

    _playlistControl.playing.addListener(_echo);
    _playlistControl.volume.addListener(_echo);

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
    _playlistControl.playing.removeListener(_echo);
    _playlistControl.volume.removeListener(_echo);
    _rcSubscription?.cancel();
    _socket.close();
  }

  @override
  Widget build(BuildContext context) => widget.child;
}
