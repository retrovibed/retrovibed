import 'dart:async';

import 'package:flutter/material.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/media/play.queue.dart' as playqueue;
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

  void _echoCurrent() {
    final cur = media.Playlist.of(context)?.queue.current.value;
    if (cur == null) return;
    _rc?.send(
      remote.Stream(
        sid: uuidx.random(),
        queue: remote.Queue(media: cur.current),
      ),
    );
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
    final playlist = media.Playlist.of(context);
    if (playlist == null) return;
    print("executing remote command ${msg}");

    switch (msg.whichCommand()) {
      case remote.Stream_Command.queue:
        playlist.queue.push(playqueue.PlayableMedia(msg.queue.media));
        break;
      case remote.Stream_Command.dequeue:
        playlist.queue.remove(msg.dequeue.id);
        break;
      case remote.Stream_Command.playpause:
        msg.playpause.paused ? playlist.player.pause() : playlist.player.play();
        break;
      case remote.Stream_Command.seek:
        final offset = msg.seek.offset;
        if (offset == 0x7FFFFFFF) {
          playlist.next();
        } else if (offset == -0x80000000) {
          playlist.previous();
        } else {
          playlist.player.seek(playlist.player.state.position + Duration(milliseconds: offset));
        }
        break;
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
      _rc?.send(
        remote.Stream(
          sid: uuidx.v7(),
          playpause: remote.PlayPause(paused: !playing),
        ),
      );
    });
    _current = media.Playlist.of(context)?.queue.current;
    _current?.addListener(_echoCurrent);

    _rcConnect();
  }

  @override
  void dispose() {
    super.dispose();
    _current?.removeListener(_echoCurrent);
    _playingSubscription?.cancel();
    _rcSubscription?.cancel();
    _rc?.close();
  }

  @override
  Widget build(BuildContext context) => widget.child;
}
