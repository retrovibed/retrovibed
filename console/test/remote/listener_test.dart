import 'dart:async';

import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:media_kit/media_kit.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/authz.dart' as authz;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/media/play.queue.dart' as playqueue;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/remote/api.dart' as remote;
import 'package:retrovibed/remote/listener.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

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

void expectFullyPopulatedSync(remote.Stream msg, {required meta.Daemon library, int? queueLength}) {
  expect(msg.whichCommand(), remote.Stream_Command.sync);
  expect(msg.sync.token, isNotEmpty);
  expect(msg.sync.token, "bearer test-bearer");
  expect(msg.sync.expiration, greaterThan(fixnum.Int64(DateTime.now().millisecondsSinceEpoch ~/ 1000)));
  expect(msg.sync.library.hostname, library.hostname);
  expect(msg.sync.capacity, greaterThan(0));
  if (queueLength != null) expect(msg.sync.queue.length, queueLength);
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

  testWidgets('sends a fully populated sync on mount', (tester) async {
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
            const SizedBox(),
          ),
        ),
      ),
    );
    await _settle(tester);

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
            const SizedBox(),
          ),
        ),
      ),
    );
    await _settle(tester);
    final before = fakeSocket.sent.length;

    fakeSocket.emit(remote.Stream(sid: "req-1", sync: remote.Sync()));
    await tester.pump();

    expect(fakeSocket.sent.length, greaterThan(before));
    expectFullyPopulatedSync(fakeSocket.sent.last, library: meta.Daemon(hostname: "localhost:9998"));
  });

  testWidgets('changing the library sends a fully populated sync', (tester) async {
    final fakeSocket = _FakeRemoteControlSocket();

    late BuildContext capturedContext;

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
                capturedContext = context;
                return const SizedBox();
              },
            ),
          ),
        ),
      ),
    );
    await _settle(tester);

    final otherDaemon = meta.Daemon(hostname: "otherhost:1234");
    meta.EndpointAuto.of(capturedContext)!.changed.value = otherDaemon;
    await tester.pump();

    expectFullyPopulatedSync(fakeSocket.sent.last, library: otherDaemon);
  });

  testWidgets('refreshing the auth token sends a fully populated sync', (tester) async {
    final fakeSocket = _FakeRemoteControlSocket();

    late BuildContext capturedContext;

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
                capturedContext = context;
                return const SizedBox();
              },
            ),
          ),
        ),
      ),
    );
    await _settle(tester);
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
    final fakeSocket = _FakeRemoteControlSocket();

    late BuildContext capturedContext;

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
                capturedContext = context;
                return const SizedBox();
              },
            ),
          ),
        ),
      ),
    );
    await _settle(tester);

    media.Playlist.of(capturedContext)!.queue.current.value = playqueue.PlayableMedia(
      media.Media(id: "now-playing"),
    );
    await tester.pump();

    expectFullyPopulatedSync(fakeSocket.sent.last, library: meta.Daemon(hostname: "localhost:9998"));
    expect(fakeSocket.sent.last.sync.current.id, "now-playing");
  });
}
