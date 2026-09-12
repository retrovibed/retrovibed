import 'dart:async';

import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/library.dart' as lib;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/remote/api.dart' as remote;
import 'package:retrovibed/remote/connect.dart';
import 'package:retrovibed/remote/playlist.current.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;

Future<Stream<meta.Daemon>> _noopDaemonDiscover({List<httpx.Option> options = const []}) async {
  return const Stream<meta.Daemon>.empty();
}

Future<meta.DaemonSearchResponse> _noopDaemonSearch(meta.DaemonSearchRequest req) async {
  return meta.DaemonSearchResponse(items: const [], next: req);
}

// never emits/closes on its own, so the _connect() future it's awaited
// under stays pending and _socket never gets reset back to noop by
// _reconnect(). echoes a sync (reflecting everything queued so far) right
// back after every queue send, like a real daemon acking one item at a
// time instead of batching - this is what races _fillQueue's own pass.
class _FakeRemoteControlSocket implements remote.RemoteControlSocket {
  final StreamController<remote.Stream> _incoming = StreamController();
  final List<remote.Stream> sent = [];
  // Sync.queue entries are now the whole enqueue Stream frame (carrying
  // sid/profile_id/session_id alongside the media), not bare Media.
  final List<remote.Stream> _acked = [];
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
      _acked.add(msg);
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

  // test helper: pushes an arbitrary incoming frame, as if the daemon sent it.
  void emit(remote.Stream msg) => _incoming.add(msg);
}

media.Media _remoteTrack(String id, String description) => media.Media(
  id: id,
  description: description,
  mimetype: 'audio/mp3',
  createdAt: '2025-01-01T00:00:00Z',
  archiveId: uuidx.min(),
  torrentId: uuidx.min(),
  knownMediaId: uuidx.min(),
);

// Mounts a connected Connect, opens the search view, submits `query`, and
// taps the (sole) `item` result - mirroring the tap-to-play flow the
// "every enqueued item carries the same non-empty session id" test already
// exercises, but driving a real query through SearchMinimal's SearchTray so
// _sessionQuery.query ends up set to something the recent-watch tests can
// assert against. Returns the session id the resulting queue send carried,
// for feeding back into a "mine" playback echo.
Future<(_FakeRemoteControlSocket, List<lib.RecentRecordRequest>, String)> _mountAndPlayUnderQuery(
  WidgetTester tester, {
  required String query,
  required media.Media item,
  Duration recentRecordThrottle = const Duration(seconds: 3),
}) async {
  final socket = _FakeRemoteControlSocket();
  final daemon = ValueNotifier(meta.Daemon());
  final recorded = <lib.RecentRecordRequest>[];

  // autoqueueTarget defaults to 5 - a single search result gets exhausted
  // by _fillQueue and falls back to the real (unauthenticated, in tests)
  // apirandom endpoint, retrying forever. Pad with filler tracks so the
  // autoqueue never needs to fall back, mirroring the existing
  // "_fillQueue overshoots"/"every enqueued item..." tests.
  Future<media.MediaSearchResponse> fakeSearch(
    media.MediaSearchRequest req, {
    String? host,
    List<httpx.Option> options = const [],
  }) async {
    return media.MediaSearchResponse(
      items: [item, ...List.generate(9, (i) => _remoteTrack('filler-${item.id}-$i', 'Filler ${item.id} $i'))],
      next: req,
    );
  }

  await tester.pumpApp(
    authn.Endpoint(
      Connect(
        search: ValueNotifier(media.MediaSearchState(next: media.MediaSearchRequest())),
        daemonDiscover: _noopDaemonDiscover,
        daemonSearch: _noopDaemonSearch,
        connect: ({required String host, List<httpx.Option> options = const []}) async => socket,
        apisearch: (host, options) => fakeSearch,
        apirandom: media.media.randomendpoint,
        autoqueueTarget: 5,
        apirecentrecord: (req, {String? host, List<httpx.Option> options = const []}) async {
          recorded.add(req);
          return lib.RecentRecordResponse();
        },
        recentRecordThrottle: recentRecordThrottle,
      ),
      daemon: daemon,
    ),
  );
  await tester.pumpN(5);

  daemon.value = meta.Daemon(hostname: "example.remote:1234");
  await tester.pumpN(5);

  await tester.tap(find.byIcon(Icons.close));
  await tester.pumpN(2);
  await tester.tap(
    find.byWidgetPredicate((w) => w is ds.LoadingIconButton && w.tooltip == "search the remote device's library"),
  );
  await tester.pumpN(5);

  await tester.enterText(
    find.descendant(of: find.byKey(const ValueKey("search")), matching: find.byType(TextField)),
    query,
  );
  await tester.tap(find.byIcon(Icons.search_rounded));
  await tester.pumpN(5);

  final row = tester.widget<media.RowDisplay>(
    find.ancestor(of: find.text(item.description).first, matching: find.byType(media.RowDisplay)),
  );
  await row.onTap!();
  await tester.pumpN(10);

  final sessionId = socket.sent.firstWhere((m) => m.whichCommand() == remote.Stream_Command.queue).sessionId;

  return (socket, recorded, sessionId);
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
    final other = meta.Daemon(id: uuidx.v7(), description: "Living Room", hostname: "example.remote:1234");

    await tester.pumpApp(
      AutoConnect(
        search: ValueNotifier(media.MediaSearchState(next: media.MediaSearchRequest())),
        daemonDiscover: _noopDaemonDiscover,
        daemonSearch: (req) async => meta.DaemonSearchResponse(items: [other], next: req),
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

    expect(find.text("Living Room"), findsOneWidget);
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

  testWidgets('autoqueue toggle starts on and uses a distinct icon from search', (tester) async {
    final socket = _FakeRemoteControlSocket();

    final daemon = ValueNotifier(meta.Daemon());

    await tester.pumpApp(
      authn.Endpoint(
        Connect(
          search: ValueNotifier(media.MediaSearchState(next: media.MediaSearchRequest())),
          daemonDiscover: _noopDaemonDiscover,
          daemonSearch: _noopDaemonSearch,
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

    expect(queueButton.toggled, isTrue);
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
          daemonSearch: _noopDaemonSearch,
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

    expect(tester.widget<ds.LoadingIconButton>(queueButtonFinder).toggled, isTrue);

    await tester.tap(queueButtonFinder);
    await tester.pumpN(2);
    expect(tester.widget<ds.LoadingIconButton>(queueButtonFinder).toggled, isFalse);

    await tester.tap(queueButtonFinder);
    await tester.pumpN(2);
    expect(tester.widget<ds.LoadingIconButton>(queueButtonFinder).toggled, isTrue);

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
                    daemonSearch: _noopDaemonSearch,
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

    // autoqueue is enabled by default (see initState) - no need to tap the
    // toggle here, just confirm unmounting while it's on doesn't throw.
    expect(
      tester
          .widget<ds.LoadingIconButton>(
            find.byWidgetPredicate((w) => w is ds.LoadingIconButton && w.tooltip == "enable autoqueue playback"),
          )
          .toggled,
      isTrue,
    );

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
          daemonSearch: _noopDaemonSearch,
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

  testWidgets('every enqueued item carries the same non-empty session id', (tester) async {
    final socket = _FakeRemoteControlSocket();

    Future<media.MediaSearchResponse> fakeSearch(
      media.MediaSearchRequest req, {
      String? host,
      List<httpx.Option> options = const [],
    }) async {
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
          daemonSearch: _noopDaemonSearch,
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

    await tester.tap(find.byIcon(Icons.close));
    await tester.pumpN(2);
    await tester.tap(
      find.byWidgetPredicate((w) => w is ds.LoadingIconButton && w.tooltip == "search the remote device's library"),
    );
    await tester.pumpN(5);

    final row = tester.widget<media.RowDisplay>(
      find.ancestor(of: find.text('Track 0').first, matching: find.byType(media.RowDisplay)),
    );

    await row.onTap!();
    await tester.pumpN(30);

    final queued = socket.sent.where((m) => m.whichCommand() == remote.Stream_Command.queue).toList();
    final sessionIds = queued.map((m) => m.sessionId).toSet();

    expect(sessionIds, hasLength(1), reason: 'every send from one Connect mount should share one session id');
    expect(sessionIds.single, isNotEmpty);
    expect(sessionIds.single, isNot(uuidx.min()));
    expect(tester.takeException(), isNull);
  });

  group('recent-watch recording', () {
    testWidgets(
      'records the search query, position, and duration when the daemon reports this session\'s own pick as current',
      (tester) async {
        final item = _remoteTrack('m1', 'Track A');
        final (socket, recorded, sessionId) = await _mountAndPlayUnderQuery(tester, query: 'my query', item: item);

        socket.emit(
          remote.Stream(
            sessionId: sessionId,
            playback: remote.Playback(media: item, position: fixnum.Int64(15000), duration: fixnum.Int64(180000)),
          ),
        );
        await tester.pump(const Duration(seconds: 3));

        expect(recorded, hasLength(1));
        expect(recorded.single.media.id, item.id);
        expect(recorded.single.query.query, 'my query');
        expect(recorded.single.position, fixnum.Int64(15000));
        expect(recorded.single.duration, fixnum.Int64(180000));
        expect(tester.takeException(), isNull);
      },
      timeout: const Timeout(Duration(seconds: 10)),
    );

    testWidgets('does not record when another session\'s pick becomes current', (tester) async {
      final item = _remoteTrack('m1', 'Track A');
      final (socket, recorded, _) = await _mountAndPlayUnderQuery(tester, query: 'my query', item: item);

      socket.emit(
        remote.Stream(
          sessionId: uuidx.v7(),
          playback: remote.Playback(media: item, position: fixnum.Int64(15000), duration: fixnum.Int64(180000)),
        ),
      );
      await tester.pump(const Duration(seconds: 3));

      expect(recorded, isEmpty);
      expect(tester.takeException(), isNull);
    }, timeout: const Timeout(Duration(seconds: 10)));

    testWidgets(
      'regenerates the session id per query and attributes each recording to its own query',
      (
        tester,
      ) async {
        final itemA = _remoteTrack('mA', 'Track A');
        final itemB = _remoteTrack('mB', 'Track B');
        // autoqueueTarget defaults to 5 - pad with filler tracks so the
        // autoqueue never exhausts the fake search and falls back to the
        // real (unauthenticated, in tests) apirandom endpoint, which would
        // retry forever.
        final items = [itemA, ...List.generate(9, (i) => _remoteTrack('filler-A-$i', 'Filler A $i'))];

        Future<media.MediaSearchResponse> fakeSearch(
          media.MediaSearchRequest req, {
          String? host,
          List<httpx.Option> options = const [],
        }) async {
          return media.MediaSearchResponse(items: items, next: req);
        }

        final socket = _FakeRemoteControlSocket();
        final daemon = ValueNotifier(meta.Daemon());
        final recorded = <lib.RecentRecordRequest>[];

        await tester.pumpApp(
          authn.Endpoint(
            Connect(
              search: ValueNotifier(media.MediaSearchState(next: media.MediaSearchRequest())),
              daemonDiscover: _noopDaemonDiscover,
              daemonSearch: _noopDaemonSearch,
              connect: ({required String host, List<httpx.Option> options = const []}) async => socket,
              apisearch: (host, options) => fakeSearch,
              apirandom: media.media.randomendpoint,
              autoqueueTarget: 5,
              apirecentrecord: (req, {String? host, List<httpx.Option> options = const []}) async {
                recorded.add(req);
                return lib.RecentRecordResponse();
              },
            ),
            daemon: daemon,
          ),
        );
        await tester.pumpN(5);

        daemon.value = meta.Daemon(hostname: "example.remote:1234");
        await tester.pumpN(5);

        await tester.tap(find.byIcon(Icons.close));
        await tester.pumpN(2);
        await tester.tap(
          find.byWidgetPredicate((w) => w is ds.LoadingIconButton && w.tooltip == "search the remote device's library"),
        );
        await tester.pumpN(5);

        await tester.enterText(
          find.descendant(of: find.byKey(const ValueKey("search")), matching: find.byType(TextField)),
          'A',
        );
        await tester.tap(find.byIcon(Icons.search_rounded));
        await tester.pumpN(5);

        final rowA = tester.widget<media.RowDisplay>(
          find.ancestor(of: find.text('Track A').first, matching: find.byType(media.RowDisplay)),
        );
        await rowA.onTap!();
        await tester.pumpN(10);

        // _sessionID is private state - PlaylistCurrent (rendered unconditionally
        // once _onPlay resets _focused to null) is the one place it's exposed to
        // the widget tree, via mySessionId. A second real queue send isn't
        // guaranteed here - _fillQueue only sends when the daemon's last-known
        // queue depth is below autoqueueTarget, which the first tap already fills.
        final sessionIdA = tester.widget<PlaylistCurrent>(find.byType(PlaylistCurrent)).sessionId!;

        socket.emit(
          remote.Stream(
            sessionId: sessionIdA,
            playback: remote.Playback(media: itemA, position: fixnum.Int64(1000), duration: fixnum.Int64(60000)),
          ),
        );
        await tester.pump(const Duration(seconds: 3));

        // back at the default view (_onPlay reset _focused to null on
        // completion) - reopen search and submit a different query.
        items
          ..clear()
          ..addAll([itemB, ...List.generate(9, (i) => _remoteTrack('filler-B-$i', 'Filler B $i'))]);
        await tester.tap(
          find.byWidgetPredicate((w) => w is ds.LoadingIconButton && w.tooltip == "search the remote device's library"),
        );
        await tester.pumpN(5);
        await tester.enterText(
          find.descendant(of: find.byKey(const ValueKey("search")), matching: find.byType(TextField)),
          'B',
        );
        await tester.tap(find.byIcon(Icons.search_rounded));
        await tester.pumpN(5);

        final rowB = tester.widget<media.RowDisplay>(
          find.ancestor(of: find.text('Track B').first, matching: find.byType(media.RowDisplay)),
        );
        await rowB.onTap!();
        await tester.pumpN(10);

        final sessionIdB = tester.widget<PlaylistCurrent>(find.byType(PlaylistCurrent)).sessionId!;
        expect(sessionIdA, isNot(sessionIdB));

        socket.emit(
          remote.Stream(
            sessionId: sessionIdB,
            playback: remote.Playback(media: itemB, position: fixnum.Int64(2000), duration: fixnum.Int64(90000)),
          ),
        );
        await tester.pump(const Duration(seconds: 3));

        expect(recorded, hasLength(2));
        expect(recorded[0].query.query, 'A');
        expect(recorded[0].media.id, itemA.id);
        expect(recorded[1].query.query, 'B');
        expect(recorded[1].media.id, itemB.id);
        expect(tester.takeException(), isNull);
      },
      timeout: const Timeout(Duration(seconds: 10)),
    );

    testWidgets('records again as playback advances, respecting the throttle window', (tester) async {
      final item = _remoteTrack('m1', 'Track A');
      const throttle = Duration(seconds: 3);
      final (socket, recorded, sessionId) = await _mountAndPlayUnderQuery(
        tester,
        query: 'my query',
        item: item,
        recentRecordThrottle: throttle,
      );

      socket.emit(
        remote.Stream(
          sessionId: sessionId,
          playback: remote.Playback(media: item, position: fixnum.Int64(1000), duration: fixnum.Int64(180000)),
        ),
      );
      await tester.pump(throttle);

      socket.emit(
        remote.Stream(
          sessionId: sessionId,
          playback: remote.Playback(media: item, position: fixnum.Int64(4000), duration: fixnum.Int64(180000)),
        ),
      );
      await tester.pump(throttle);

      expect(recorded, hasLength(2));
      expect(recorded[0].position, fixnum.Int64(1000));
      expect(recorded[1].position, fixnum.Int64(4000));
      expect(tester.takeException(), isNull);
    }, timeout: const Timeout(Duration(seconds: 10)));

    testWidgets('throttles rapid playback echoes into a single record per window', (tester) async {
      final item = _remoteTrack('m1', 'Track A');
      const throttle = Duration(seconds: 3);
      final (socket, recorded, sessionId) = await _mountAndPlayUnderQuery(
        tester,
        query: 'my query',
        item: item,
        recentRecordThrottle: throttle,
      );

      // throttle(duration, trailing: true) emits the *first* event of a
      // window immediately (leading edge) - this one opens the window.
      socket.emit(
        remote.Stream(
          sessionId: sessionId,
          playback: remote.Playback(media: item, position: fixnum.Int64(500), duration: fixnum.Int64(180000)),
        ),
      );

      // two rapid echoes racing within that same still-open window - only
      // the *last* of these should survive to the trailing edge; the
      // 1000 in between must be dropped, not recorded on its own.
      socket.emit(
        remote.Stream(
          sessionId: sessionId,
          playback: remote.Playback(media: item, position: fixnum.Int64(1000), duration: fixnum.Int64(180000)),
        ),
      );
      socket.emit(
        remote.Stream(
          sessionId: sessionId,
          playback: remote.Playback(media: item, position: fixnum.Int64(2000), duration: fixnum.Int64(180000)),
        ),
      );
      await tester.pump(throttle);

      expect(recorded, hasLength(2));
      expect(recorded[0].position, fixnum.Int64(500));
      expect(recorded[1].position, fixnum.Int64(2000));
      expect(tester.takeException(), isNull);
    }, timeout: const Timeout(Duration(seconds: 10)));
  });
}
