import 'dart:async';

import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/media/play.queue.dart' as playqueue;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/retrovibed.dart' as retro;
import 'api.dart' as remote;

// Mounted as an ancestor of media.Playlist's UI-consuming subtree, drives the
// same public Playlist surface the on-screen transport buttons use
// (playlist.next()/.previous(), playlist.player, playlist.queue) rather than
// reaching into Playlist's private state directly.
class RemoteControlListener extends StatefulWidget {
  final Widget child;
  const RemoteControlListener(this.child, {super.key});

  @override
  State<RemoteControlListener> createState() => _State();
}

class _State extends State<RemoteControlListener> {
  remote.RemoteControlSocket? _rc;
  StreamSubscription<remote.Stream>? _rcSubscription;
  StreamSubscription<bool>? _playingSubscription;
  // Cached in initState because dispose() runs after the element tree is
  // deactivated, when context.findAncestorStateOfType is unsafe to call.
  ValueNotifier<playqueue.PlayableMedia?>? _current;
  ValueNotifier<meta.Daemon>? _library;
  ValueNotifier<int>? _queueRevision;

  // now that Sync carries the full state, redirect the old ad-hoc track-change
  // echo through it instead of sending a bare Queue frame.
  void _echoCurrent() => _echoSync();

  // proactively (not just in reply to a request) reports the listener's
  // current library, playback queue, and a token valid against this device -
  // triggered by track/library/queue changes and by didChangeDependencies
  // whenever the token cache actually refreshes.
  void _echoSync() {
    final queue = media.Playlist.of(context)?.queue ?? playqueue.PlayQueue();
    final library = meta.EndpointAuto.of(context)?.changed.value;
    final cached = authn.AuthzCache.meta(context).current; // authz.Bearer<meta.Token>
    _rc?.send(remote.messages.syncrsp(
      library: library,
      token: httpx.bearer(cached.bearer),
      expiration: cached.token.expires,
      capacity: queue.capacity,
      current: queue.current.value?.current,
      queue: queue.queued.map((m) => m.current).toList(),
    ));
  }

  void _rcReconnect() {
    if (!mounted) return;
    Future.delayed(const Duration(seconds: 2), _rcConnect);
  }

  // listen is process-local and never depends on the profile being logged
  // in, so it connects unconditionally and reconnects forever on close.
  void _rcConnect() {
    remote.remotecontrol
        .listen(options: [httpx.Request.bearer(() => Future.value(retro.remote_control_listen_token()))])
        .then((socket) {
          _rc = socket;
          final c = Completer();

          _rcSubscription = socket.messages.listen(
            _applyRemoteCommand,
            cancelOnError: true,
            onError: c.completeError,
            onDone: c.complete,
          );

          return c.future;
        })
        .catchError((cause) {
          debugPrint("remote control listen socket failed: ${cause}");
        })
        .whenComplete(() {
          _rc = null;
          _rcSubscription = null;
          _rcReconnect();
        });
  }

  void _applyRemoteCommand(remote.Stream msg) {
    print("received remote command ${msg}");

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
      case remote.Stream_Command.sync:
      case remote.Stream_Command.notSet:
        break;
    }
  }

  @override
  void initState() {
    super.initState();

    // echo local state back over the listen socket so any /rc/connect
    // observers can see what this device is doing.
    _playingSubscription = media.Playlist.of(context)?.player.stream.playing.listen((playing) {
      _rc?.send(remote.messages.playpause(!playing));
    });
    _current = media.Playlist.of(context)?.queue.current;
    _current?.addListener(_echoCurrent);
    _library = meta.EndpointAuto.of(context)?.changed;
    _library?.addListener(_echoSync);
    _queueRevision = media.Playlist.of(context)?.queue.revision;
    _queueRevision?.addListener(_echoSync);

    _rcConnect();
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    authn.AuthzCache.meta(context); // registers this State's dependency on AuthzTokenData
    _echoSync();
  }

  @override
  void dispose() {
    super.dispose();
    _current?.removeListener(_echoCurrent);
    _library?.removeListener(_echoSync);
    _queueRevision?.removeListener(_echoSync);
    _playingSubscription?.cancel();
    _rcSubscription?.cancel();
    _rc?.close();
  }

  @override
  Widget build(BuildContext context) => widget.child;
}
