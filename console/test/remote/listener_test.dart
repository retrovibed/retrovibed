import 'dart:async';

import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:media_kit/media_kit.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/authz.dart' as authz;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/media/play.queue.dart' as playqueue;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/remote/api.dart' as remote;
import 'package:retrovibed/remote/listener.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;

class _FakeRemoteControlSocket implements remote.RemoteControlSocket {
  final StreamController<remote.Stream> _incoming = StreamController();
  final List<remote.Stream> sent = [];

  @override
  Stream<remote.Stream> get messages => _incoming.stream;

  @override
  void send(remote.Stream msg) => sent.add(msg);

  @override
  Future<void> close() async {
    await _incoming.close();
  }

  void emit(remote.Stream msg) => _incoming.add(msg);
}

// media.Playlist mounts a real media_kit Player, whose native streams keep
// scheduling frames outside the fake-async test clock, so pumpAndSettle()
// never settles here (same reason routes_test.dart avoids it too).
Future<void> _settle(WidgetTester tester) async {
  for (var i = 0; i < 10; i++) {
    await tester.pump(const Duration(milliseconds: 50));
  }
}

void expectFullyPopulatedSync(
  remote.Stream msg, {
  required meta.Daemon library,
  int? queueLength,
  bool? fullscreen,
  bool? muted,
}) {
  expect(msg.whichCommand(), remote.Stream_Command.sync);
  expect(msg.sync.token, isNotEmpty);
  expect(msg.sync.token, "bearer test-bearer");
  expect(msg.sync.expiration, greaterThan(fixnum.Int64(DateTime.now().millisecondsSinceEpoch ~/ 1000)));
  expect(msg.sync.library.hostname, library.hostname);
  expect(msg.sync.capacity, greaterThan(0));
  if (queueLength != null) expect(msg.sync.queue.length, queueLength);
  if (fullscreen != null) expect(msg.sync.fullscreen, fullscreen);
  if (muted != null) expect(msg.sync.muted, muted);
}

void main() {
  setUpAll(() {
    MediaKit.ensureInitialized();
  });

  // Expiration must stay in the future: AuthzCache treats an already-expired
  // token as always-stale, and RemoteControlListener re-echoes a sync every
  // time the token cache notifies of a "refresh" - a fixed past timestamp
  // makes that feedback loop spin forever.
  Future<meta.AuthzResponse> fixedAuth({String? host}) async {
    return meta.AuthzResponse(
      bearer: "test-bearer",
      token: meta.Token(expires: fixnum.Int64((DateTime.now().millisecondsSinceEpoch ~/ 1000) + 3600)),
    );
  }

  Future<meta.DaemonLookupResponse> mockLatest() async {
    return meta.DaemonLookupResponse(daemon: meta.Daemon(hostname: "localhost:9998"));
  }

  Future<meta.Daemon> mockConnectable(meta.Daemon d) async => d;

  // Mounts RemoteControlListener with the standard EndpointAuto/Playlist
  // scaffold every test needs, optionally without a Playlist ancestor (to
  // exercise the "no playlist" guard) or with a Full ancestor (to exercise
  // the fullscreen command's actual toggle rather than its safe no-op).
  Future<(_FakeRemoteControlSocket, BuildContext)> mount(
    WidgetTester tester, {
    bool withPlaylist = true,
    bool withFull = false,
  }) async {
    final fakeSocket = _FakeRemoteControlSocket();
    late BuildContext capturedContext;

    Widget child = RemoteControlListener(
      connect: ({List<httpx.Option> options = const []}) async => fakeSocket,
      localDevice: () => meta.Daemon(hostname: "localhost:9998"),
      Builder(
        builder: (context) {
          capturedContext = context;
          return const SizedBox();
        },
      ),
    );

    if (withPlaylist) child = media.Playlist(child);

    Widget tree = meta.EndpointAuto(
      latest: mockLatest,
      connectable: mockConnectable,
      backoff: httpx.Backoff.constant(Duration.zero),
      child,
    );

    if (withFull) tree = ds.Full(tree);

    await tester.pumpApp(authzCurrent: fixedAuth, tree);
    await _settle(tester);

    return (fakeSocket, capturedContext);
  }

  // Shared shape of the _sid ordering guard, identical across playpause /
  // seek / volume / mute / fullscreen: accept once (newer sid), attempt a
  // stale sid (rejected), attempt the exact same sid again (rejected, >=
  // not >), then recover with a genuinely newer sid (accepted). accept/
  // reject/recover must each produce a distinguishable probe() outcome so a
  // wrongly-applied "reject" or a wrongly-ignored "recover" is caught.
  Future<void> expectSidGuarded(
    WidgetTester tester,
    _FakeRemoteControlSocket socket,
    remote.Stream Function(String sid) accept,
    remote.Stream Function(String sid) reject,
    remote.Stream Function(String sid) recover,
    Object? Function() probe,
  ) async {
    final now = DateTime.now();
    final acceptedSid = uuidx.v7(at: now);

    socket.emit(accept(acceptedSid));
    await _settle(tester);
    final afterAccept = probe();

    socket.emit(reject(uuidx.v7(at: now.subtract(const Duration(minutes: 1)))));
    await _settle(tester);
    expect(probe(), afterAccept, reason: "a stale sid (generated earlier, arriving late) must not apply");

    socket.emit(reject(acceptedSid));
    await _settle(tester);
    expect(probe(), afterAccept, reason: "an equal sid must not apply (guard is >=, not >)");

    socket.emit(recover(uuidx.v7(at: now.add(const Duration(seconds: 1)))));
    await _settle(tester);
    expect(probe(), isNot(afterAccept), reason: "a genuinely newer sid must apply - guard isn't latched");
  }

  group('sync', () {
    testWidgets('sends a fully populated sync on mount', (tester) async {
      final (fakeSocket, _) = await mount(tester);

      expect(fakeSocket.sent, isNotEmpty);
      expectFullyPopulatedSync(fakeSocket.sent.last, library: meta.Daemon(hostname: "localhost:9998"));
    });

    testWidgets('pushing to the queue sends a fully populated sync', (tester) async {
      final fakeSocket = _FakeRemoteControlSocket();

      await tester.pumpApp(
        authzCurrent: fixedAuth,
        meta.EndpointAuto(
          latest: mockLatest,
          connectable: mockConnectable,
          backoff: httpx.Backoff.constant(Duration.zero),
          media.Playlist(
            RemoteControlListener(
              connect: ({List<httpx.Option> options = const []}) async => fakeSocket,
              localDevice: () => meta.Daemon(hostname: "localhost:9998"),
              Builder(
                builder: (context) {
                  return ElevatedButton(
                    onPressed: () {
                      media.Playlist.of(context)!.queue.push(playqueue.PlayableMedia(media.Media(id: "m1")));
                    },
                    child: const SizedBox(),
                  );
                },
              ),
            ),
          ),
        ),
      );
      await _settle(tester);

      await tester.tap(find.byType(ElevatedButton));
      await tester.pump();

      expectFullyPopulatedSync(fakeSocket.sent.last, library: meta.Daemon(hostname: "localhost:9998"), queueLength: 1);
      expect(fakeSocket.sent.last.sync.queue.single.id, "m1");
    });

    testWidgets('an incoming sync request triggers a fully populated response', (tester) async {
      final (fakeSocket, _) = await mount(tester);
      final before = fakeSocket.sent.length;

      fakeSocket.emit(remote.Stream(sid: "req-1", sync: remote.Sync()));
      await tester.pump();

      expect(fakeSocket.sent.length, greaterThan(before));
      expectFullyPopulatedSync(fakeSocket.sent.last, library: meta.Daemon(hostname: "localhost:9998"));
    });

    testWidgets('an incoming sync request with an old sid still gets a reply - sync has no ordering guard', (
      tester,
    ) async {
      final (fakeSocket, _) = await mount(tester);
      final before = fakeSocket.sent.length;

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(at: DateTime(2000)), sync: remote.Sync()));
      await tester.pump();

      expect(fakeSocket.sent.length, greaterThan(before));
      expectFullyPopulatedSync(fakeSocket.sent.last, library: meta.Daemon(hostname: "localhost:9998"));
    });

    testWidgets('changing the library sends a fully populated sync', (tester) async {
      final (fakeSocket, capturedContext) = await mount(tester);

      final otherDaemon = meta.Daemon(hostname: "otherhost:1234");
      meta.EndpointAuto.of(capturedContext)!.changed.value = otherDaemon;
      await tester.pump();

      expectFullyPopulatedSync(fakeSocket.sent.last, library: otherDaemon);
    });

    testWidgets('refreshing the auth token sends a fully populated sync', (tester) async {
      final (fakeSocket, capturedContext) = await mount(tester);
      final before = fakeSocket.sent.length;

      authn.AuthzCache.of(capturedContext).changed.value = authz.Bearer(
        meta.Token(expires: fixnum.Int64(9876543210)),
        "refreshed-bearer",
      );
      await tester.pump();

      expect(fakeSocket.sent.length, greaterThan(before));
      expectFullyPopulatedSync(fakeSocket.sent.last, library: meta.Daemon(hostname: "localhost:9998"));
    });

    testWidgets('changing the current media sends a fully populated sync', (tester) async {
      final (fakeSocket, capturedContext) = await mount(tester);

      media.Playlist.of(capturedContext)!.queue.current.value = playqueue.PlayableMedia(
        media.Media(id: "now-playing"),
      );
      await tester.pump();

      expectFullyPopulatedSync(fakeSocket.sent.last, library: meta.Daemon(hostname: "localhost:9998"));
      expect(fakeSocket.sent.last.sync.current.id, "now-playing");
    });
  });

  group('queue', () {
    // unguarded: no _sid ordering protection at all, unlike playpause/seek/volume/mute/fullscreen.
    testWidgets('an inbound queue command pushes onto the playlist queue', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), queue: remote.Queue(media: media.Media(id: "m1"))));
      await _settle(tester);

      expect(playlist.queue.queued.map((m) => m.current.id), contains("m1"));
    });

    testWidgets('two inbound queue commands both land, in delivery order', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), queue: remote.Queue(media: media.Media(id: "m1"))));
      await _settle(tester);
      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), queue: remote.Queue(media: media.Media(id: "m2"))));
      await _settle(tester);

      expect(playlist.queue.queued.map((m) => m.current.id).toList(), ["m1", "m2"]);
    });

    testWidgets('a queue command with a sid older than already seen still applies', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;
      final now = DateTime.now();

      // advance the listener's known sid past "now" via an accepted (guarded) mute first.
      fakeSocket.emit(remote.Stream(sid: uuidx.v7(at: now), mute: remote.Mute()));
      await _settle(tester);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: now.subtract(const Duration(minutes: 1))),
          queue: remote.Queue(media: media.Media(id: "m1")),
        ),
      );
      await _settle(tester);

      expect(playlist.queue.queued.map((m) => m.current.id), contains("m1"));
    });

    testWidgets('still produces a trailing sync echo', (tester) async {
      final (fakeSocket, _) = await mount(tester);
      final before = fakeSocket.sent.length;

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), queue: remote.Queue(media: media.Media(id: "m1"))));
      await _settle(tester);

      expect(fakeSocket.sent.length, greaterThan(before));
      expect(fakeSocket.sent.last.whichCommand(), remote.Stream_Command.sync);
    });
  });

  group('dequeue', () {
    testWidgets('removes a previously queued id', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;
      playlist.queue.push(playqueue.PlayableMedia(media.Media(id: "m1")));

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), dequeue: remote.Dequeue(id: "m1")));
      await _settle(tester);

      expect(playlist.queue.queued.map((m) => m.current.id), isNot(contains("m1")));
    });

    testWidgets('dequeuing an absent id is a no-op', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;
      playlist.queue.push(playqueue.PlayableMedia(media.Media(id: "m1")));

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), dequeue: remote.Dequeue(id: "does-not-exist")));
      await _settle(tester);

      expect(playlist.queue.queued.map((m) => m.current.id).toList(), ["m1"]);
    });

    testWidgets('a dequeue arriving before its queue is a no-op, not a pending cancellation', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), dequeue: remote.Dequeue(id: "m1")));
      await _settle(tester);
      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), queue: remote.Queue(media: media.Media(id: "m1"))));
      await _settle(tester);

      expect(playlist.queue.queued.map((m) => m.current.id), contains("m1"));
    });

    testWidgets('duplicate/replayed dequeues for the same id are idempotent', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;
      playlist.queue.push(playqueue.PlayableMedia(media.Media(id: "m1")));

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), dequeue: remote.Dequeue(id: "m1")));
      await _settle(tester);
      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), dequeue: remote.Dequeue(id: "m1")));
      await _settle(tester);

      expect(playlist.queue.queued, isEmpty);
    });
  });

  group('playpause', () {
    testWidgets('paused:false plays', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), playpause: remote.PlayPause(paused: false)));
      await _settle(tester);

      expect(playlist.player.state.playing, isTrue);
    });

    testWidgets('paused:true pauses', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;
      await playlist.player.play();
      await _settle(tester);

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), playpause: remote.PlayPause(paused: true)));
      await _settle(tester);

      expect(playlist.player.state.playing, isFalse);
    });

    testWidgets('sid ordering guard', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;

      await expectSidGuarded(
        tester,
        fakeSocket,
        (sid) => remote.Stream(sid: sid, playpause: remote.PlayPause(paused: false)),
        (sid) => remote.Stream(sid: sid, playpause: remote.PlayPause(paused: true)),
        (sid) => remote.Stream(sid: sid, playpause: remote.PlayPause(paused: true)),
        () => playlist.player.state.playing,
      );
    });
  });

  group('seek', () {
    testWidgets('offset==SeekOffset.next advances via playlist.next()', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;
      playlist.queue.push(playqueue.PlayableMedia(media.Media(id: "m1")));
      playlist.queue.push(playqueue.PlayableMedia(media.Media(id: "m2")));

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), seek: remote.Seek(offset: remote.SeekOffset.next)));
      await _settle(tester);

      expect(playlist.queue.queued.length, 1);
    });

    testWidgets('offset==SeekOffset.previous reverses via playlist.previous()', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;
      playlist.queue.push(playqueue.PlayableMedia(media.Media(id: "m1")));
      playlist.queue.push(playqueue.PlayableMedia(media.Media(id: "m2")));
      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), seek: remote.Seek(offset: remote.SeekOffset.next)));
      await _settle(tester);
      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), seek: remote.Seek(offset: remote.SeekOffset.next)));
      await _settle(tester);
      final beforePrevious = playlist.queue.upcoming;

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), seek: remote.Seek(offset: remote.SeekOffset.previous)));
      await _settle(tester);

      expect(tester.takeException(), isNull);
      expect(playlist.queue.upcoming, greaterThan(beforePrevious));
    });

    testWidgets('sid ordering guard', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;
      for (var i = 0; i < 4; i++) {
        playlist.queue.push(playqueue.PlayableMedia(media.Media(id: "m$i")));
      }

      // probes via queue drain since the headless test backend doesn't
      // advance player.state.position without a loaded track.
      await expectSidGuarded(
        tester,
        fakeSocket,
        (sid) => remote.Stream(sid: sid, seek: remote.Seek(offset: remote.SeekOffset.next)),
        (sid) => remote.Stream(sid: sid, seek: remote.Seek(offset: remote.SeekOffset.next)),
        (sid) => remote.Stream(sid: sid, seek: remote.Seek(offset: remote.SeekOffset.next)),
        () => playlist.queue.queued.length,
      );
    });
  });

  group('volume', () {
    testWidgets('a relative adjustment adds to the current volume', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;
      await playlist.player.setVolume(50.0);
      await _settle(tester);

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), volume: remote.Seek(offset: 10)));
      await _settle(tester);

      expect(playlist.player.state.volume, 60.0);
    });

    testWidgets('clamps at the top', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;
      await playlist.player.setVolume(90.0);
      await _settle(tester);

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), volume: remote.Seek(offset: 50)));
      await _settle(tester);

      expect(playlist.player.state.volume, 100.0);
    });

    testWidgets('clamps at the bottom', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;
      await playlist.player.setVolume(10.0);
      await _settle(tester);

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), volume: remote.Seek(offset: -50)));
      await _settle(tester);

      expect(playlist.player.state.volume, 0.0);
    });

    testWidgets('a negative offset decreases volume', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;
      await playlist.player.setVolume(50.0);
      await _settle(tester);

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), volume: remote.Seek(offset: -15)));
      await _settle(tester);

      expect(playlist.player.state.volume, 35.0);
    });

    testWidgets('adjusting volume while muted unmutes and applies against the remembered level', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;
      await playlist.player.setVolume(60.0);
      await _settle(tester);

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), mute: remote.Mute()));
      await _settle(tester);
      expect(playlist.player.state.volume, 0.0);

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), volume: remote.Seek(offset: 10)));
      await _settle(tester);

      expect(playlist.player.state.volume, 70.0);
      expectFullyPopulatedSync(fakeSocket.sent.last, library: meta.Daemon(hostname: "localhost:9998"), muted: false);
    });

    testWidgets('the trailing sync echo reports the resulting volume', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;
      await playlist.player.setVolume(20.0);
      await _settle(tester);

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), volume: remote.Seek(offset: 30)));
      await _settle(tester);

      expect(fakeSocket.sent.last.whichCommand(), remote.Stream_Command.sync);
      expect(fakeSocket.sent.last.sync.volume, 50.0);
    });

    testWidgets('sid ordering guard', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;
      await playlist.player.setVolume(50.0);
      await _settle(tester);

      await expectSidGuarded(
        tester,
        fakeSocket,
        (sid) => remote.Stream(sid: sid, volume: remote.Seek(offset: 20)),
        (sid) => remote.Stream(sid: sid, volume: remote.Seek(offset: 5)),
        (sid) => remote.Stream(sid: sid, volume: remote.Seek(offset: 7)),
        () => playlist.player.state.volume,
      );
    });
  });

  group('mute', () {
    testWidgets('mutes without losing the remembered level', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;
      await playlist.player.setVolume(60.0);
      await _settle(tester);

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), mute: remote.Mute()));
      await _settle(tester);

      expect(playlist.player.state.volume, 0.0);
      expectFullyPopulatedSync(fakeSocket.sent.last, library: meta.Daemon(hostname: "localhost:9998"), muted: true);
      expect(fakeSocket.sent.last.sync.volume, 60.0);
    });

    testWidgets('unmuting restores the exact pre-mute level, not a flat max', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;
      await playlist.player.setVolume(60.0);
      await _settle(tester);
      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), mute: remote.Mute()));
      await _settle(tester);

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), mute: remote.Mute()));
      await _settle(tester);

      expect(playlist.player.state.volume, 60.0);
      expectFullyPopulatedSync(fakeSocket.sent.last, library: meta.Daemon(hostname: "localhost:9998"), muted: false);
    });

    testWidgets('two mutes back to back end where they started', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;
      await playlist.player.setVolume(42.0);
      await _settle(tester);

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), mute: remote.Mute()));
      await _settle(tester);
      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), mute: remote.Mute()));
      await _settle(tester);

      expect(playlist.player.state.volume, 42.0);
    });

    testWidgets('a local (non-remote) volume change updates the remembered level', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;

      // simulates a keyboard-shortcut-driven change (playlist.dart's ctrl+arrow),
      // which calls player.setVolume directly, bypassing the remote-command path.
      await playlist.player.setVolume(77.0);
      await _settle(tester);

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), mute: remote.Mute()));
      await _settle(tester);
      expect(playlist.player.state.volume, 0.0);

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), mute: remote.Mute()));
      await _settle(tester);

      expect(playlist.player.state.volume, 77.0);
    });

    testWidgets('sid ordering guard', (tester) async {
      final (fakeSocket, context) = await mount(tester);
      final playlist = media.Playlist.of(context)!;
      await playlist.player.setVolume(50.0);
      await _settle(tester);

      await expectSidGuarded(
        tester,
        fakeSocket,
        (sid) => remote.Stream(sid: sid, mute: remote.Mute()),
        (sid) => remote.Stream(sid: sid, mute: remote.Mute()),
        (sid) => remote.Stream(sid: sid, mute: remote.Mute()),
        () => playlist.player.state.volume,
      );
    });
  });

  group('fullscreen', () {
    testWidgets('with no Full ancestor, does not throw and reports false', (tester) async {
      final (fakeSocket, _) = await mount(tester, withFull: false);

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), fullscreen: remote.Fullscreen()));
      await _settle(tester);

      expect(tester.takeException(), isNull);
      expectFullyPopulatedSync(
        fakeSocket.sent.last,
        library: meta.Daemon(hostname: "localhost:9998"),
        fullscreen: false,
      );
    });

    testWidgets('with a Full ancestor mounted, toggles and reports the flip', (tester) async {
      final calls = <String>[];
      TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger.setMockMethodCallHandler(
        const MethodChannel('window_manager'),
        (call) async {
          calls.add(call.method);
          return null;
        },
      );
      addTearDown(() {
        TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger.setMockMethodCallHandler(
          const MethodChannel('window_manager'),
          null,
        );
      });

      final (fakeSocket, _) = await mount(tester, withFull: true);

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), fullscreen: remote.Fullscreen()));
      await _settle(tester);

      expect(calls, contains('setFullScreen'));
      expectFullyPopulatedSync(
        fakeSocket.sent.last,
        library: meta.Daemon(hostname: "localhost:9998"),
        fullscreen: true,
      );

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: DateTime.now().add(const Duration(seconds: 1))),
          fullscreen: remote.Fullscreen(),
        ),
      );
      await _settle(tester);

      expectFullyPopulatedSync(
        fakeSocket.sent.last,
        library: meta.Daemon(hostname: "localhost:9998"),
        fullscreen: false,
      );
    });

    testWidgets('sid ordering guard', (tester) async {
      final (fakeSocket, context) = await mount(tester);

      await expectSidGuarded(
        tester,
        fakeSocket,
        (sid) => remote.Stream(sid: sid, fullscreen: remote.Fullscreen()),
        (sid) => remote.Stream(sid: sid, fullscreen: remote.Fullscreen()),
        (sid) => remote.Stream(sid: sid, fullscreen: remote.Fullscreen()),
        () => ds.Full.nochrome(context),
      );
    });

    testWidgets('works even without a Playlist ancestor - unlike the other five guarded/unguarded commands', (
      tester,
    ) async {
      TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger.setMockMethodCallHandler(
        const MethodChannel('window_manager'),
        (call) async => null,
      );
      addTearDown(() {
        TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger.setMockMethodCallHandler(
          const MethodChannel('window_manager'),
          null,
        );
      });

      final (fakeSocket, context) = await mount(tester, withPlaylist: false, withFull: true);

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), fullscreen: remote.Fullscreen()));
      await _settle(tester);

      expect(tester.takeException(), isNull);
      expect(ds.Full.nochrome(context), isTrue);
    });
  });

  group('playlist ancestor absent', () {
    testWidgets('queue/dequeue/playpause/seek/volume/mute are silently dropped, no throw', (tester) async {
      final (fakeSocket, context) = await mount(tester, withPlaylist: false);

      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), queue: remote.Queue(media: media.Media(id: "m1"))));
      await _settle(tester);
      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), dequeue: remote.Dequeue(id: "m1")));
      await _settle(tester);
      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), playpause: remote.PlayPause(paused: false)));
      await _settle(tester);
      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), seek: remote.Seek(offset: remote.SeekOffset.next)));
      await _settle(tester);
      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), volume: remote.Seek(offset: 10)));
      await _settle(tester);
      fakeSocket.emit(remote.Stream(sid: uuidx.v7(), mute: remote.Mute()));
      await _settle(tester);

      expect(tester.takeException(), isNull);
      expect(media.Playlist.of(context), isNull);
    });
  });

  group('notSet / malformed', () {
    testWidgets('an empty command does not throw and still produces a trailing sync echo', (tester) async {
      final (fakeSocket, _) = await mount(tester);
      final before = fakeSocket.sent.length;

      fakeSocket.emit(remote.Stream(sid: "x"));
      await _settle(tester);

      expect(tester.takeException(), isNull);
      expect(fakeSocket.sent.length, greaterThan(before));
      expect(fakeSocket.sent.last.whichCommand(), remote.Stream_Command.sync);
    });
  });
}
