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
  final StreamController<remote.Stream> _incoming = StreamController.broadcast();
  final StreamController<remote.Stream> _outgoing = StreamController.broadcast();

  @override
  Stream<remote.Stream> get messages => _incoming.stream;

  Stream<remote.Stream> get sent => _outgoing.stream;

  @override
  void send(remote.Stream msg) => _outgoing.add(msg);

  @override
  Future<void> close() async {
    await _incoming.close();
    await _outgoing.close();
  }

  void emit(remote.Stream msg) => _incoming.add(msg);
}

void expectFullyPopulatedSync(remote.Stream msg, {required meta.Daemon library, int? queueLength}) {
  expect(msg.whichCommand(), remote.Stream_Command.sync);
  expect(msg.sync.token, isNotEmpty);
  expect(msg.sync.token, "bearer test-bearer");
  expect(msg.sync.expiration, fixnum.Int64(1234567890));
  expect(msg.sync.library.hostname, library.hostname);
  expect(msg.sync.capacity, greaterThan(0));
  if (queueLength != null) expect(msg.sync.queue.length, queueLength);
}

void main() {
  setUpAll(() {
    MediaKit.ensureInitialized();
  });

  Future<meta.AuthzResponse> fixedAuth({String? host}) async {
    return meta.AuthzResponse(bearer: "test-bearer", token: meta.Token(expires: fixnum.Int64(1234567890)));
  }

  Future<meta.DaemonLookupResponse> mockLatest() async {
    return meta.DaemonLookupResponse(daemon: meta.Daemon(hostname: "localhost:9998"));
  }

  Future<meta.Daemon> mockConnectable(meta.Daemon d) async => d;

  testWidgets('sends a fully populated sync on mount', (tester) async {
    final fakeSocket = _FakeRemoteControlSocket();
    final sent = <remote.Stream>[];
    fakeSocket.sent.listen(sent.add);

    await tester.pumpApp(
      authzCurrent: fixedAuth,
      meta.EndpointAuto(
        latest: mockLatest,
        connectable: mockConnectable,
        backoff: httpx.Backoff.constant(Duration.zero),
        media.Playlist(
          RemoteControlListener(
            connect: ({List<httpx.Option> options = const []}) async => fakeSocket,
            const SizedBox(),
          ),
        ),
      ),
    );
    await tester.pumpAndSettle();

    expect(sent, isNotEmpty);
    expectFullyPopulatedSync(sent.last, library: meta.Daemon(hostname: "localhost:9998"));
  });

  testWidgets('pushing to the queue sends a fully populated sync', (tester) async {
    final fakeSocket = _FakeRemoteControlSocket();
    final sent = <remote.Stream>[];
    fakeSocket.sent.listen(sent.add);

    await tester.pumpApp(
      authzCurrent: fixedAuth,
      meta.EndpointAuto(
        latest: mockLatest,
        connectable: mockConnectable,
        backoff: httpx.Backoff.constant(Duration.zero),
        media.Playlist(
          RemoteControlListener(
            connect: ({List<httpx.Option> options = const []}) async => fakeSocket,
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
    await tester.pumpAndSettle();

    await tester.tap(find.byType(ElevatedButton));
    await tester.pump();

    expectFullyPopulatedSync(sent.last, library: meta.Daemon(hostname: "localhost:9998"), queueLength: 1);
    expect(sent.last.sync.queue.single.id, "m1");
  });

  testWidgets('an incoming sync request triggers a fully populated response', (tester) async {
    final fakeSocket = _FakeRemoteControlSocket();
    final sent = <remote.Stream>[];
    fakeSocket.sent.listen(sent.add);

    await tester.pumpApp(
      authzCurrent: fixedAuth,
      meta.EndpointAuto(
        latest: mockLatest,
        connectable: mockConnectable,
        backoff: httpx.Backoff.constant(Duration.zero),
        media.Playlist(
          RemoteControlListener(
            connect: ({List<httpx.Option> options = const []}) async => fakeSocket,
            const SizedBox(),
          ),
        ),
      ),
    );
    await tester.pumpAndSettle();
    final before = sent.length;

    fakeSocket.emit(remote.Stream(sid: "req-1", sync: remote.Sync()));
    await tester.pump();

    expect(sent.length, greaterThan(before));
    expectFullyPopulatedSync(sent.last, library: meta.Daemon(hostname: "localhost:9998"));
  });

  testWidgets('changing the library sends a fully populated sync', (tester) async {
    final fakeSocket = _FakeRemoteControlSocket();
    final sent = <remote.Stream>[];
    fakeSocket.sent.listen(sent.add);

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
    await tester.pumpAndSettle();

    final otherDaemon = meta.Daemon(hostname: "otherhost:1234");
    meta.EndpointAuto.of(capturedContext)!.changed.value = otherDaemon;
    await tester.pump();

    expectFullyPopulatedSync(sent.last, library: otherDaemon);
  });

  testWidgets('refreshing the auth token sends a fully populated sync', (tester) async {
    final fakeSocket = _FakeRemoteControlSocket();
    final sent = <remote.Stream>[];
    fakeSocket.sent.listen(sent.add);

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
    await tester.pumpAndSettle();
    final before = sent.length;

    authn.AuthzCache.of(capturedContext).changed.value = authz.Bearer(
      meta.Token(expires: fixnum.Int64(9876543210)),
      "refreshed-bearer",
    );
    await tester.pump();

    expect(sent.length, greaterThan(before));
    expectFullyPopulatedSync(sent.last, library: meta.Daemon(hostname: "localhost:9998"));
  });

  testWidgets('changing the current media sends a fully populated sync', (tester) async {
    final fakeSocket = _FakeRemoteControlSocket();
    final sent = <remote.Stream>[];
    fakeSocket.sent.listen(sent.add);

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
    await tester.pumpAndSettle();

    media.Playlist.of(capturedContext)!.queue.current.value = playqueue.PlayableMedia(
      media.Media(id: "now-playing"),
    );
    await tester.pump();

    expectFullyPopulatedSync(sent.last, library: meta.Daemon(hostname: "localhost:9998"));
    expect(sent.last.sync.current.id, "now-playing");
  });
}
