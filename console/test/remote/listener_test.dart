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

// Plain-Dart stand-in for media.PlaylistControl - no native player, no
// MediaKit dependency, fully synchronous. Used for every command group
// except sync/fullscreen/playlist-ancestor-absent, which either already
// worked fine against the real media.Playlist or don't touch it at all.
// Deliberately simpler than real Playlist semantics: maybeNext always
// pushes (no "auto-advance if idle"), next()/previous() are plain counters
// - these tests verify RemoteControlListener calls the right PlaylistControl
// method with the right arguments, not Playlist's own playback orchestration.
class _FakePlaylistControl implements media.PlaylistControl {
  @override
  final playqueue.PlayQueue queue = playqueue.PlayQueue();
  @override
  final ValueNotifier<bool> playing = ValueNotifier(false);
  @override
  final ValueNotifier<double> volume;
  Duration position;
  final List<playqueue.PlayableMedia> maybeNextCalls = [];
  int nextCalls = 0;
  int previousCalls = 0;
  final List<Duration> seekCalls = [];

  _FakePlaylistControl({double volume = 100.0, this.position = Duration.zero}) : volume = ValueNotifier(volume);

  @override
  void maybeNext(playqueue.PlayableMedia m) {
    maybeNextCalls.add(m);
    queue.push(m);
  }

  @override
  void next() => nextCalls++;

  @override
  void previous() => previousCalls++;

  @override
  void playOrPause() => playing.value = !playing.value;

  @override
  void seek(Duration p) {
    seekCalls.add(p);
    position = p;
  }

  @override
  Future<void> setVolume(double v) async => volume.value = v;
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
  bool? paused,
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
  if (paused != null) expect(msg.sync.paused, paused);
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
      token: meta.Token(exp: fixnum.Int64((DateTime.now().millisecondsSinceEpoch ~/ 1000) + 3600)),
    );
  }

  Future<meta.DaemonLookupResponse> mockLatest() async {
    return meta.DaemonLookupResponse(daemon: meta.Daemon(hostname: "localhost:9998"));
  }

  Future<meta.Daemon> mockConnectable(meta.Daemon d) async => d;

  // Mounts RemoteControlListener with the standard EndpointAuto scaffold
  // every test needs. withPlaylist mounts a real media.Playlist ancestor;
  // fakePlaylist (when given) injects a _FakePlaylistControl via
  // RemoteControlListener's playlist accessor instead - independent knobs,
  // since the fake-based command groups need neither a real Playlist nor
  // the default "no ancestor -> PlaylistControl.zero" fallback.
  Future<(_FakeRemoteControlSocket, BuildContext)> mount(
    WidgetTester tester, {
    bool withPlaylist = true,
    bool withFull = false,
    media.PlaylistControl? fakePlaylist,
  }) async {
    final fakeSocket = _FakeRemoteControlSocket();
    late BuildContext capturedContext;

    Widget child = RemoteControlListener(
      connect: ({List<httpx.Option> options = const []}) async => fakeSocket,
      localDevice: () => meta.Daemon(hostname: "localhost:9998"),
      playlist: fakePlaylist != null
          ? (context) => fakePlaylist
          : (context) => media.Playlist.of(context) ?? media.PlaylistControl.zero,
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

  // Shared shape of the _sid ordering guard, identical across pause / seek /
  // volume / mute / fullscreen: accept once (newer sid), attempt a stale
  // sid (rejected), attempt the exact same sid again (rejected, >= not >),
  // then recover with a genuinely newer sid (accepted). accept/reject/
  // recover must each produce a distinguishable probe() outcome so a
  // wrongly-applied "reject" or a wrongly-ignored "recover" is caught.
  Future<void> expectSidGuarded(
    WidgetTester tester,
    _FakeRemoteControlSocket socket,
    remote.Stream Function(String sid) accept,
    remote.Stream Function(String sid) reject,
    remote.Stream Function(String sid) recover,
    Object? Function() probe,
  ) async {
    // _sid isn't just "last accepted command's sid" - _echoSync fires after
    // EVERY inbound message (accepted or rejected, via .whenComplete) and
    // always sets _sid = sync.sid using a fresh real-time uuidx.v7(), which
    // keeps advancing _sid regardless of what any given command's own sid
    // was. So an "equal sid" reuse can't just replay whatever sid a prior
    // step sent - it has to reuse whatever _sid actually holds *now*, which
    // is exactly whatever sid the most recent outgoing echo carried (since
    // that's what set it).
    final currentSid = () => socket.sent.last.sid;

    // +10s: safely past the automatic post-mount sync echo (fires once auth
    // resolves), so "accept" here can't itself lose a same-millisecond race
    // against that echo's own sid.
    socket.emit(accept(uuidx.v7(at: DateTime.now().add(const Duration(seconds: 10)))));
    await _settle(tester);
    final afterAccept = probe();

    // always safely in the past, however far _sid has since drifted forward via echoes.
    socket.emit(reject(uuidx.v7(at: DateTime(2000))));
    await _settle(tester);
    expect(probe(), afterAccept, reason: "a stale sid (generated earlier, arriving late) must not apply");

    socket.emit(reject(currentSid()));
    await _settle(tester);
    expect(probe(), afterAccept, reason: "an equal sid must not apply (guard is >=, not >)");

    // safely ahead of whatever _sid has drifted to by now via the echoes above.
    socket.emit(recover(uuidx.v7(at: DateTime.now().add(const Duration(seconds: 20)))));
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

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: DateTime(2000)),
          sync: remote.Sync(),
        ),
      );
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
        meta.Token(exp: fixnum.Int64(9876543210)),
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
    // unguarded: no _sid ordering protection at all, unlike pause/seek/volume/mute/fullscreen.
    testWidgets('an inbound queue command pushes onto the playlist queue', (tester) async {
      final fake = _FakePlaylistControl();
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(),
          queue: remote.Queue(media: media.Media(id: "m1")),
        ),
      );
      await _settle(tester);

      expect(fake.queue.queued.map((m) => m.current.id), contains("m1"));
    });

    testWidgets('two inbound queue commands both land, in delivery order', (tester) async {
      final fake = _FakePlaylistControl();
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(),
          queue: remote.Queue(media: media.Media(id: "m1")),
        ),
      );
      await _settle(tester);
      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(),
          queue: remote.Queue(media: media.Media(id: "m2")),
        ),
      );
      await _settle(tester);

      expect(fake.queue.queued.map((m) => m.current.id).toList(), ["m1", "m2"]);
    });

    testWidgets('a queue command with a sid older than already seen still applies', (tester) async {
      final fake = _FakePlaylistControl();
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);
      final now = DateTime.now().add(const Duration(seconds: 10));

      // advance the listener's known sid past "now" via an accepted (guarded) mute first.
      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: now),
          mute: remote.Mute(),
        ),
      );
      await _settle(tester);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: now.subtract(const Duration(minutes: 1))),
          queue: remote.Queue(media: media.Media(id: "m1")),
        ),
      );
      await _settle(tester);

      expect(fake.queue.queued.map((m) => m.current.id), contains("m1"));
    });

    testWidgets('still produces a trailing sync echo', (tester) async {
      final fake = _FakePlaylistControl();
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);
      final before = fakeSocket.sent.length;

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(),
          queue: remote.Queue(media: media.Media(id: "m1")),
        ),
      );
      await _settle(tester);

      expect(fakeSocket.sent.length, greaterThan(before));
      expect(fakeSocket.sent.last.whichCommand(), remote.Stream_Command.sync);
    });
  });

  group('dequeue', () {
    testWidgets('removes a previously queued id', (tester) async {
      final fake = _FakePlaylistControl();
      fake.queue.push(playqueue.PlayableMedia(media.Media(id: "m1")));
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(),
          dequeue: remote.Dequeue(id: "m1"),
        ),
      );
      await _settle(tester);

      expect(fake.queue.queued.map((m) => m.current.id), isNot(contains("m1")));
    });

    testWidgets('dequeuing an absent id is a no-op', (tester) async {
      final fake = _FakePlaylistControl();
      fake.queue.push(playqueue.PlayableMedia(media.Media(id: "m1")));
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(),
          dequeue: remote.Dequeue(id: "does-not-exist"),
        ),
      );
      await _settle(tester);

      expect(fake.queue.queued.map((m) => m.current.id).toList(), ["m1"]);
    });

    testWidgets('a dequeue arriving before its queue is a no-op, not a pending cancellation', (tester) async {
      final fake = _FakePlaylistControl();
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(),
          dequeue: remote.Dequeue(id: "m1"),
        ),
      );
      await _settle(tester);
      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(),
          queue: remote.Queue(media: media.Media(id: "m1")),
        ),
      );
      await _settle(tester);

      expect(fake.queue.queued.map((m) => m.current.id), contains("m1"));
    });

    testWidgets('duplicate/replayed dequeues for the same id are idempotent', (tester) async {
      final fake = _FakePlaylistControl();
      fake.queue.push(playqueue.PlayableMedia(media.Media(id: "m1")));
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(),
          dequeue: remote.Dequeue(id: "m1"),
        ),
      );
      await _settle(tester);
      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(),
          dequeue: remote.Dequeue(id: "m1"),
        ),
      );
      await _settle(tester);

      expect(fake.queue.queued, isEmpty);
    });

    // _echoSync reads the queue via an async token-cache lookup
    // (AuthzCache.meta(context).auto()) whose latency is unbounded - if it
    // reads queue.queued only once that lookup resolves rather than
    // snapshotting it at the moment _echoSync was triggered, a second
    // dequeue landing while the first echo is still waiting on the token
    // makes the first echo report the queue as of the *second* dequeue
    // instead of its own. Requires a real media.Playlist ancestor: the fake
    // PlaylistControl used elsewhere in this group owns its own queue
    // instance, separate from the one RemoteControlListener reads for sync.
    testWidgets(
      'a sync blocked on a slow token fetch still reports the queue as of its own trigger, not a later one',
      (tester) async {
        final fakeSocket = _FakeRemoteControlSocket();
        final authCompleter = Completer<meta.AuthzResponse>();
        // AuthzCache hides its child behind a loading placeholder until the
        // very first current() call resolves, so that one has to succeed
        // immediately (with an already-expired token, so every later auto()
        // call is treated as stale and re-invokes current()) - every call
        // after it shares the one pending completer below.
        var authResolved = false;
        Future<meta.AuthzResponse> delayedAuth({String? host}) {
          if (!authResolved) {
            authResolved = true;
            return Future.value(
              meta.AuthzResponse(bearer: "initial-bearer", token: meta.Token(exp: fixnum.Int64(0))),
            );
          }
          return authCompleter.future;
        }

        late BuildContext capturedContext;

        await tester.pumpApp(
          authzCurrent: delayedAuth,
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
                    capturedContext = context;
                    return const SizedBox();
                  },
                ),
              ),
            ),
          ),
        );
        await _settle(tester);

        final queue = media.Playlist.of(capturedContext)!.queue;
        for (final id in ["m1", "m2", "m3", "m4", "m5"]) {
          queue.push(playqueue.PlayableMedia(media.Media(id: id)));
        }
        await _settle(tester);

        fakeSocket.emit(remote.Stream(sid: uuidx.v7(), dequeue: remote.Dequeue(id: "m1")));
        await _settle(tester);
        fakeSocket.emit(remote.Stream(sid: uuidx.v7(), dequeue: remote.Dequeue(id: "m2")));
        await _settle(tester);

        expect(fakeSocket.sent, isEmpty, reason: "every echo triggered so far is still waiting on the token fetch");

        authCompleter.complete(
          meta.AuthzResponse(
            bearer: "test-bearer",
            token: meta.Token(exp: fixnum.Int64((DateTime.now().millisecondsSinceEpoch ~/ 1000) + 3600)),
          ),
        );
        await _settle(tester);

        expect(
          fakeSocket.sent.any((m) => m.whichCommand() == remote.Stream_Command.sync && m.sync.queue.length == 4),
          isTrue,
          reason: "the echo triggered by removing m1 (5 -> 4) must report 4, not whatever the queue has "
              "dropped to (3) by the time the delayed token fetch resolves",
        );
      },
    );
  });

  group('pause', () {
    // no payload - toggles play/pause, same shape as mute/fullscreen.
    testWidgets('toggles from playing to paused', (tester) async {
      final fake = _FakePlaylistControl()..playing.value = true;
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: DateTime.now().add(const Duration(seconds: 10))),
          pause: remote.Pause(),
        ),
      );
      await _settle(tester);

      expect(fake.playing.value, isFalse);
    });

    testWidgets('toggles from paused to playing', (tester) async {
      final fake = _FakePlaylistControl();
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: DateTime.now().add(const Duration(seconds: 10))),
          pause: remote.Pause(),
        ),
      );
      await _settle(tester);

      expect(fake.playing.value, isTrue);
    });

    testWidgets('sid ordering guard', (tester) async {
      final fake = _FakePlaylistControl();
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);

      await expectSidGuarded(
        tester,
        fakeSocket,
        (sid) => remote.Stream(sid: sid, pause: remote.Pause()),
        (sid) => remote.Stream(sid: sid, pause: remote.Pause()),
        (sid) => remote.Stream(sid: sid, pause: remote.Pause()),
        () => fake.playing.value,
      );
    });
  });

  group('seek', () {
    testWidgets('a relative offset seeks from the current position', (tester) async {
      final fake = _FakePlaylistControl(position: const Duration(seconds: 10));
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: DateTime.now().add(const Duration(seconds: 10))),
          seek: remote.Seek(offset: 5000),
        ),
      );
      await _settle(tester);

      expect(fake.seekCalls.last, const Duration(seconds: 15));
    });

    testWidgets('offset==SeekOffset.next calls next()', (tester) async {
      final fake = _FakePlaylistControl();
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: DateTime.now().add(const Duration(seconds: 10))),
          seek: remote.Seek(offset: remote.SeekOffset.next),
        ),
      );
      await _settle(tester);

      expect(fake.nextCalls, 1);
      expect(fake.seekCalls, isEmpty);
    });

    testWidgets('offset==SeekOffset.previous calls previous()', (tester) async {
      final fake = _FakePlaylistControl();
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: DateTime.now().add(const Duration(seconds: 10))),
          seek: remote.Seek(offset: remote.SeekOffset.previous),
        ),
      );
      await _settle(tester);

      expect(fake.previousCalls, 1);
      expect(fake.seekCalls, isEmpty);
    });

    testWidgets('sid ordering guard', (tester) async {
      final fake = _FakePlaylistControl();
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);

      await expectSidGuarded(
        tester,
        fakeSocket,
        (sid) => remote.Stream(
          sid: sid,
          seek: remote.Seek(offset: remote.SeekOffset.next),
        ),
        (sid) => remote.Stream(
          sid: sid,
          seek: remote.Seek(offset: remote.SeekOffset.next),
        ),
        (sid) => remote.Stream(
          sid: sid,
          seek: remote.Seek(offset: remote.SeekOffset.next),
        ),
        () => fake.nextCalls,
      );
    });
  });

  group('volume', () {
    testWidgets('a relative adjustment adds to the current volume', (tester) async {
      final fake = _FakePlaylistControl(volume: 50.0);
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: DateTime.now().add(const Duration(seconds: 10))),
          volume: remote.Seek(offset: 10),
        ),
      );
      await _settle(tester);

      expect(fake.volume.value, 60.0);
    });

    testWidgets('clamps at the top', (tester) async {
      final fake = _FakePlaylistControl(volume: 90.0);
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: DateTime.now().add(const Duration(seconds: 10))),
          volume: remote.Seek(offset: 50),
        ),
      );
      await _settle(tester);

      expect(fake.volume.value, 100.0);
    });

    testWidgets('clamps at the bottom', (tester) async {
      final fake = _FakePlaylistControl(volume: 10.0);
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: DateTime.now().add(const Duration(seconds: 10))),
          volume: remote.Seek(offset: -50),
        ),
      );
      await _settle(tester);

      expect(fake.volume.value, 0.0);
    });

    testWidgets('a negative offset decreases volume', (tester) async {
      final fake = _FakePlaylistControl(volume: 50.0);
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: DateTime.now().add(const Duration(seconds: 10))),
          volume: remote.Seek(offset: -15),
        ),
      );
      await _settle(tester);

      expect(fake.volume.value, 35.0);
    });

    testWidgets('adjusting volume while muted unmutes and applies against the remembered level', (tester) async {
      final fake = _FakePlaylistControl(volume: 60.0);
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);
      final now = DateTime.now().add(const Duration(seconds: 10));

      // explicit at: timestamps, strictly increasing: two bare uuidx.v7()
      // calls in quick succession aren't guaranteed ordered (pure random
      // sub-millisecond entropy), which would make the guard flakily reject
      // the second command here.
      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: now),
          mute: remote.Mute(),
        ),
      );
      await _settle(tester);
      expect(fake.volume.value, 0.0);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: now.add(const Duration(seconds: 1))),
          volume: remote.Seek(offset: 10),
        ),
      );
      await _settle(tester);

      expect(fake.volume.value, 70.0);
      expectFullyPopulatedSync(fakeSocket.sent.last, library: meta.Daemon(hostname: "localhost:9998"), muted: false);
    });

    testWidgets('the trailing sync echo reports the resulting volume', (tester) async {
      final fake = _FakePlaylistControl(volume: 20.0);
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: DateTime.now().add(const Duration(seconds: 10))),
          volume: remote.Seek(offset: 30),
        ),
      );
      await _settle(tester);

      expect(fakeSocket.sent.last.whichCommand(), remote.Stream_Command.sync);
      expect(fakeSocket.sent.last.sync.volume, 50.0);
    });

    testWidgets('sid ordering guard', (tester) async {
      final fake = _FakePlaylistControl(volume: 50.0);
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);

      await expectSidGuarded(
        tester,
        fakeSocket,
        (sid) => remote.Stream(sid: sid, volume: remote.Seek(offset: 20)),
        (sid) => remote.Stream(sid: sid, volume: remote.Seek(offset: 5)),
        (sid) => remote.Stream(sid: sid, volume: remote.Seek(offset: 7)),
        () => fake.volume.value,
      );
    });
  });

  group('mute', () {
    testWidgets('mutes and reports the live (silent) volume', (tester) async {
      final fake = _FakePlaylistControl(volume: 60.0);
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: DateTime.now().add(const Duration(seconds: 10))),
          mute: remote.Mute(),
        ),
      );
      await _settle(tester);

      expect(fake.volume.value, 0.0);
      expectFullyPopulatedSync(fakeSocket.sent.last, library: meta.Daemon(hostname: "localhost:9998"), muted: true);
      expect(fakeSocket.sent.last.sync.volume, 0.0);
    });

    testWidgets('unmuting restores the exact pre-mute level, not a flat max', (tester) async {
      final fake = _FakePlaylistControl(volume: 60.0);
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);
      final now = DateTime.now().add(const Duration(seconds: 10));

      // explicit at: timestamps, strictly increasing: two bare uuidx.v7()
      // calls in quick succession aren't guaranteed ordered (pure random
      // sub-millisecond entropy), which would make the guard flakily reject
      // the second mute here.
      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: now),
          mute: remote.Mute(),
        ),
      );
      await _settle(tester);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: now.add(const Duration(seconds: 1))),
          mute: remote.Mute(),
        ),
      );
      await _settle(tester);

      expect(fake.volume.value, 60.0);
      expectFullyPopulatedSync(fakeSocket.sent.last, library: meta.Daemon(hostname: "localhost:9998"), muted: false);
      expect(fakeSocket.sent.last.sync.volume, 60.0);
    });

    testWidgets('two mutes back to back end where they started', (tester) async {
      final fake = _FakePlaylistControl(volume: 42.0);
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);
      final now = DateTime.now().add(const Duration(seconds: 10));

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: now),
          mute: remote.Mute(),
        ),
      );
      await _settle(tester);
      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: now.add(const Duration(seconds: 1))),
          mute: remote.Mute(),
        ),
      );
      await _settle(tester);

      expect(fake.volume.value, 42.0);
    });

    testWidgets('a local (non-remote) volume change updates the remembered level', (tester) async {
      final fake = _FakePlaylistControl(volume: 100.0);
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);
      final now = DateTime.now().add(const Duration(seconds: 10));

      // simulates a keyboard-shortcut-driven change (playlist.dart's ctrl+arrow),
      // which sets the player's volume directly, bypassing the remote-command path.
      fake.volume.value = 77.0;
      await _settle(tester);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: now),
          mute: remote.Mute(),
        ),
      );
      await _settle(tester);
      expect(fake.volume.value, 0.0);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: now.add(const Duration(seconds: 1))),
          mute: remote.Mute(),
        ),
      );
      await _settle(tester);

      expect(fake.volume.value, 77.0);
    });

    testWidgets('sid ordering guard', (tester) async {
      final fake = _FakePlaylistControl(volume: 50.0);
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, fakePlaylist: fake);

      await expectSidGuarded(
        tester,
        fakeSocket,
        (sid) => remote.Stream(sid: sid, mute: remote.Mute()),
        (sid) => remote.Stream(sid: sid, mute: remote.Mute()),
        (sid) => remote.Stream(sid: sid, mute: remote.Mute()),
        () => fake.volume.value,
      );
    });
  });

  group('fullscreen', () {
    testWidgets('with no Full ancestor, does not throw and reports false', (tester) async {
      final (fakeSocket, _) = await mount(tester, withPlaylist: false, withFull: false);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: DateTime.now().add(const Duration(seconds: 10))),
          fullscreen: remote.Fullscreen(),
        ),
      );
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

      final (fakeSocket, _) = await mount(tester, withPlaylist: false, withFull: true);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: DateTime.now().add(const Duration(seconds: 10))),
          fullscreen: remote.Fullscreen(),
        ),
      );
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

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: DateTime.now().add(const Duration(seconds: 10))),
          fullscreen: remote.Fullscreen(),
        ),
      );
      await _settle(tester);

      expect(tester.takeException(), isNull);
      expect(ds.Full.nochrome(context), isTrue);
    });
  });

  group('playlist ancestor absent', () {
    testWidgets('queue/dequeue/pause/seek/volume/mute are silently dropped, no throw', (tester) async {
      final (fakeSocket, context) = await mount(tester, withPlaylist: false);

      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(),
          queue: remote.Queue(media: media.Media(id: "m1")),
        ),
      );
      await _settle(tester);
      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(),
          dequeue: remote.Dequeue(id: "m1"),
        ),
      );
      await _settle(tester);
      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: DateTime.now().add(const Duration(seconds: 10))),
          pause: remote.Pause(),
        ),
      );
      await _settle(tester);
      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: DateTime.now().add(const Duration(seconds: 10))),
          seek: remote.Seek(offset: remote.SeekOffset.next),
        ),
      );
      await _settle(tester);
      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: DateTime.now().add(const Duration(seconds: 10))),
          volume: remote.Seek(offset: 10),
        ),
      );
      await _settle(tester);
      fakeSocket.emit(
        remote.Stream(
          sid: uuidx.v7(at: DateTime.now().add(const Duration(seconds: 10))),
          mute: remote.Mute(),
        ),
      );
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
