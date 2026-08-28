import 'dart:async';

import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/remote/api.dart' as remote;
import 'package:retrovibed/remote/connect.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;

Future<Stream<meta.Daemon>> _noopDaemonDiscover({List<httpx.Option> options = const []}) async {
  return const Stream<meta.Daemon>.empty();
}

// never emits/closes on its own, so the _connect() future it's awaited
// under stays pending and _socket never gets reset back to noop by
// _reconnect(). echoes a sync (reflecting everything queued so far) right
// back after every queue send, like a real daemon acking one item at a
// time instead of batching - this is what races _fillQueue's own pass.
class _FakeRemoteControlSocket implements remote.RemoteControlSocket {
  final StreamController<remote.Stream> _incoming = StreamController();
  final List<remote.Stream> sent = [];
  final List<media.Media> _acked = [];
  // vid is a plain monotonic counter, unlike sid (a uuidv7) whose ordering
  // isn't guaranteed for two ids minted within the same millisecond - which
  // is exactly what this fake hits, firing echoes far faster than any real
  // daemon would. no artificial delay needed once ordering is by counter.
  int _vid = 0;

  @override
  Stream<remote.Stream> get messages => _incoming.stream;

  @override
  void send(remote.Stream msg) {
    print("sending ${msg.sid} ${msg.whichCommand()}");
    sent.add(msg);
    if (msg.whichCommand() == remote.Stream_Command.queue) {
      _acked.add(msg.queue.media);
      print("sync response");
      _incoming.add(
        remote.Stream(
          sid: uuidx.v7(),
          vid: fixnum.Int64(++_vid),
          sync: remote.Sync(queue: List.of(_acked)),
        ),
      );
    }
  }

  @override
  Future<void> close() async {
    await _incoming.close();
  }
}

void main() {
  testWidgets('unmounting Connect does not throw during dispose', (tester) async {
    bool visible = true;
    late StateSetter setLocalState;

    await tester.pumpApp(
      StatefulBuilder(
        builder: (context, setState) {
          setLocalState = setState;
          return visible
              ? AutoConnect(
                  search: ValueNotifier(media.MediaSearchState(next: media.MediaSearchRequest())),
                )
              : const SizedBox();
        },
      ),
    );
    await tester.pump();

    setLocalState(() => visible = false);
    await tester.pump();

    expect(tester.takeException(), isNull);
  });

  testWidgets('selecting the local device is refused without calling connect', (tester) async {
    bool connectCalled = false;

    await tester.pumpApp(
      AutoConnect(
        search: ValueNotifier(media.MediaSearchState(next: media.MediaSearchRequest())),
        daemonDiscover: _noopDaemonDiscover,
        connect: ({required String host, List<httpx.Option> options = const []}) async {
          connectCalled = true;
          throw StateError('connect should not be called for the local device');
        },
      ),
    );
    await tester.pumpN(5);

    authn.AuthedEndpoint.daemon(tester.element(find.byType(meta.DaemonDropdown))).value = meta.Daemon(
      hostname: "localhost:9998",
    );
    await tester.pumpN(5);

    expect(find.text("you do not have permission to remotely control this device"), findsOneWidget);
    expect(connectCalled, isFalse);
    expect(tester.takeException(), isNull);
  });

  testWidgets('a forbidden response from connect disables remote control', (tester) async {
    await tester.pumpApp(
      AutoConnect(
        search: ValueNotifier(media.MediaSearchState(next: media.MediaSearchRequest())),
        daemonDiscover: _noopDaemonDiscover,
        connect: ({required String host, List<httpx.Option> options = const []}) async {
          throw http.Response('', 403);
        },
      ),
    );
    await tester.pumpN(5);

    authn.AuthedEndpoint.daemon(tester.element(find.byType(meta.DaemonDropdown))).value = meta.Daemon(
      hostname: "example.remote:1234",
    );
    await tester.pumpN(5);

    expect(find.text("remote control is disabled on this device"), findsOneWidget);
    expect(tester.takeException(), isNull);
  });

  testWidgets('autoqueue toggle starts off and uses a distinct icon from search', (tester) async {
    final socket = _FakeRemoteControlSocket();

    final daemon = ValueNotifier(meta.Daemon());

    await tester.pumpApp(
      authn.Endpoint(
        Connect(
          search: ValueNotifier(media.MediaSearchState(next: media.MediaSearchRequest())),
          daemonDiscover: _noopDaemonDiscover,
          connect: ({required String host, List<httpx.Option> options = const []}) async => socket,
          apisearch: media.media.searchendpoint,
          apirandom: media.media.randomendpoint,
          autoqueueTarget: 5,
        ),
        daemon: daemon,
      ),
    );
    await tester.pumpN(5);

    daemon.value = meta.Daemon(hostname: "example.remote:1234");
    await tester.pumpN(5);

    // _onEndpointChanged() defaults _focused to the queue view - dismiss it
    // to reveal the transport-controls row with the toggle buttons.
    await tester.tap(find.byIcon(Icons.close));
    await tester.pumpN(2);

    final queueButton = tester.widget<ds.LoadingIconButton>(
      find.byWidgetPredicate((w) => w is ds.LoadingIconButton && w.tooltip == "enable autoqueue playback"),
    );
    final searchButton = tester.widget<ds.LoadingIconButton>(
      find.byWidgetPredicate((w) => w is ds.LoadingIconButton && w.tooltip == "search the remote device's library"),
    );

    expect(queueButton.toggled, isFalse);
    expect((queueButton.icon as Icon).icon, isNot((searchButton.icon as Icon).icon));
    expect(tester.takeException(), isNull);
  });

  testWidgets('tapping the autoqueue button toggles its state', (tester) async {
    final socket = _FakeRemoteControlSocket();
    final daemon = ValueNotifier(meta.Daemon());

    await tester.pumpApp(
      authn.Endpoint(
        Connect(
          search: ValueNotifier(media.MediaSearchState(next: media.MediaSearchRequest())),
          daemonDiscover: _noopDaemonDiscover,
          connect: ({required String host, List<httpx.Option> options = const []}) async => socket,
          apisearch: media.media.searchendpoint,
          apirandom: media.media.randomendpoint,
          autoqueueTarget: 5,
        ),
        daemon: daemon,
      ),
    );
    await tester.pumpN(5);

    daemon.value = meta.Daemon(hostname: "example.remote:1234");
    await tester.pumpN(5);

    // _onEndpointChanged() defaults _focused to the queue view - dismiss it
    // to reveal the transport-controls row with the toggle buttons.
    await tester.tap(find.byIcon(Icons.close));
    await tester.pumpN(2);

    final queueButtonFinder = find.byWidgetPredicate(
      (w) => w is ds.LoadingIconButton && w.tooltip == "enable autoqueue playback",
    );

    expect(tester.widget<ds.LoadingIconButton>(queueButtonFinder).toggled, isFalse);

    await tester.tap(queueButtonFinder);
    await tester.pumpN(2);
    expect(tester.widget<ds.LoadingIconButton>(queueButtonFinder).toggled, isTrue);

    await tester.tap(queueButtonFinder);
    await tester.pumpN(2);
    expect(tester.widget<ds.LoadingIconButton>(queueButtonFinder).toggled, isFalse);

    expect(tester.takeException(), isNull);
  });

  testWidgets('unmounting while autoqueue is enabled does not throw', (tester) async {
    final socket = _FakeRemoteControlSocket();
    final daemon = ValueNotifier(meta.Daemon());
    bool visible = true;
    late StateSetter setLocalState;

    await tester.pumpApp(
      authn.Endpoint(
        StatefulBuilder(
          builder: (context, setState) {
            setLocalState = setState;
            return visible
                ? Connect(
                    search: ValueNotifier(media.MediaSearchState(next: media.MediaSearchRequest())),
                    daemonDiscover: _noopDaemonDiscover,
                    connect: ({required String host, List<httpx.Option> options = const []}) async => socket,
                    apisearch: media.media.searchendpoint,
                    apirandom: media.media.randomendpoint,
                    autoqueueTarget: 5,
                  )
                : const SizedBox();
          },
        ),
        daemon: daemon,
      ),
    );
    await tester.pumpN(5);

    daemon.value = meta.Daemon(hostname: "example.remote:1234");
    await tester.pumpN(5);

    // _onEndpointChanged() defaults _focused to the queue view - dismiss it
    // to reveal the transport-controls row with the toggle buttons.
    await tester.tap(find.byIcon(Icons.close));
    await tester.pumpN(2);

    await tester.tap(
      find.byWidgetPredicate((w) => w is ds.LoadingIconButton && w.tooltip == "enable autoqueue playback"),
    );
    await tester.pumpN(2);

    setLocalState(() => visible = false);
    await tester.pump();

    expect(tester.takeException(), isNull);
  });

  testWidgets('_fillQueue overshoots the target when sync echoes race the fill loop', (tester) async {
    final socket = _FakeRemoteControlSocket();

    Future<media.MediaSearchResponse> fakeSearch(
      media.MediaSearchRequest req, {
      String? host,
      List<httpx.Option> options = const [],
    }) async {
      print("search invoked");
      return media.MediaSearchResponse(
        items: List.generate(
          10,
          (i) => media.Media(
            id: uuidx.withSuffix(i + 1),
            description: 'Track $i',
            mimetype: 'audio/mp3',
            createdAt: '2025-01-01T00:00:00Z',
            archiveId: uuidx.min(),
            torrentId: uuidx.min(),
            knownMediaId: uuidx.min(),
          ),
        ),
        next: media.media.request(limit: 32),
      );
    }

    final daemon = ValueNotifier(meta.Daemon());

    await tester.pumpApp(
      authn.Endpoint(
        Connect(
          search: ValueNotifier(media.MediaSearchState(next: media.MediaSearchRequest())),
          daemonDiscover: _noopDaemonDiscover,
          connect: ({required String host, List<httpx.Option> options = const []}) async => socket,
          apisearch: (host, options) => fakeSearch,
          apirandom: media.media.randomendpoint,
          autoqueueTarget: 5,
        ),
        daemon: daemon,
      ),
    );
    await tester.pumpN(5);

    daemon.value = meta.Daemon(hostname: "example.remote:1234");
    await tester.pumpN(5);

    // _onEndpointChanged() defaults _focused to the queue view - dismiss it,
    // then open the search view (SearchMinimal only renders once _focused
    // points at it) to reach the tappable rows.
    await tester.tap(find.byIcon(Icons.close));
    await tester.pumpN(2);
    await tester.tap(
      find.byWidgetPredicate((w) => w is ds.LoadingIconButton && w.tooltip == "search the remote device's library"),
    );
    await tester.pumpN(5);

    // tester.tap() only dispatches the gesture - it doesn't await the
    // returned onTap future, so a fixed pumpN afterward can't guarantee the
    // _onPlay -> _fillQueue -> moveNext() -> range() chain actually finishes
    // before the test body moves on. Grab RowDisplay's onTap directly and
    // await it to completion instead.
    final row = tester.widget<media.RowDisplay>(
      find.ancestor(of: find.text('Track 0').first, matching: find.byType(media.RowDisplay)),
    );

    await row.onTap!();
    await tester.pumpN(30);

    final queued = socket.sent.where((m) => m.whichCommand() == remote.Stream_Command.queue).length;

    // the autoqueue targets 5 upcoming items - each sync echo the fake
    // daemon sends back after a queue send re-triggers _fillQueue until the
    // target is reached.
    expect(queued, equals(5));
    expect(tester.takeException(), isNull);
  });
}
