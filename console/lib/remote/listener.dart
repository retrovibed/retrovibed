import 'dart:async';
import 'dart:ffi';
import 'dart:io';

import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/authz.dart' as authz;
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
  const RemoteControlListener(this.child, {super.key, this.connect = remote.remotecontrol.listen});

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
  remote.Volume _volume = remote.Volume(muted: false, level: 0.0);

  // proactively (not just in reply to a request) reports the listener's
  // current library, playback queue, and a token valid against this device -
  // triggered by track/library/queue changes and by didChangeDependencies
  // whenever the token cache actually refreshes.
  Future<void> _echoSync() async {
    final queue = _queue;
    final library = _library.value;
    library..hostname = meta.daemons.isLocalDevice(library) ? retro.local_device().hostname : library.hostname;
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
        volume: media.Playlist.of(context)?.player.state.volume ?? 0.0,
      );
      _socket.send(
        sync,
      );
      setState(() {
        _sid = sync.sid;
      });
    });
  }

  _echoVolume() {
    final sync = remote.messages.volume(
      media.Playlist.of(context)?.player.state.volume ?? _volume.level,
      _volume.muted,
    );
    _socket.send(
      sync,
    );
    setState(() {
      _sid = sync.sid;
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
            _applyRemoteCommand,
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

  void _applyRemoteCommand(remote.Stream msg) {
    if (msg.whichCommand() == remote.Stream_Command.sync) {
      _echoSync();
      return;
    }

    final playlist = media.Playlist.of(context);
    if (playlist == null) return;
    print("executing remote command ${msg}");

    switch (msg.whichCommand()) {
      case remote.Stream_Command.queue:
        playlist.maybeNext(playqueue.PlayableMedia(msg.queue.media));
        break;
      case remote.Stream_Command.dequeue:
        playlist.queue.remove(msg.dequeue.id);
        break;
      case remote.Stream_Command.playpause:
        msg.playpause.paused ? playlist.player.pause() : playlist.player.play();
        break;
      case remote.Stream_Command.seek:
        final offset = msg.seek.offset;
        if (offset == remote.SeekOffset.next) {
          playlist.next();
        } else if (offset == remote.SeekOffset.previous) {
          playlist.previous();
        } else {
          playlist.player.seek(playlist.player.state.position + Duration(milliseconds: offset));
        }
        break;
      case remote.Stream_Command.volume:
        if (_sid.compareTo(msg.sid) >= 0) {
          _echoVolume();
          break;
        }

        setState(() {
          _sid = msg.sid;
          _volume = msg.volume;
        });

        final v = msg.volume.muted ? 0.0 : msg.volume.level;
        playlist.player.setVolume(v).then((_) {
          _echoVolume();
        });
        break;
      case remote.Stream_Command.sync:
      case remote.Stream_Command.notSet:
        break;
    }
  }

  @override
  void initState() {
    super.initState();

    final _player = media.Playlist.of(context)?.player;

    _volume = Volume(muted: false, level: _player?.state.volume ?? 0.0);
    // echo local state back over the listen socket so any /rc/connect
    // observers can see what this device is doing.
    _playingSubscription = _player?.stream.playing.listen((playing) {
      _socket.send(remote.messages.playpause(!playing));
    });

    _volumeSubscription = _player?.stream.volume.listen((volume) {
      _socket.send(remote.messages.volume(volume, _volume.muted));
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
