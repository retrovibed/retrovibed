import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/discovery.dart' as disc;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/library.dart' as lib;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/remote/api.dart' as remote;
import 'package:retrovibed/remote/empty.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;

void main() {
  testWidgets('renders recent items from apirecentlatest when the queue is empty', (tester) async {
    await tester.pumpApp(
      Empty(
        remote.Sync(token: 'bearer test'),
        search: ValueNotifier(media.MediaSearchState(next: media.MediaSearchRequest())),
        socket: remote.RemoteControlSocket.noop,
        sessionID: uuidx.v7(),
        onPlay: (context, m, resp) => null,
        apisearch: (req, {String? host, List<httpx.Option> options = const []}) async {
          return media.MediaSearchResponse(next: req);
        },
        apirecentlatest: (req, {String? host, List<httpx.Option> options = const []}) async {
          return media.RecentSearchResponse(
            items: [
              media.RecentRecordRequest(
                id: '1',
                media: media.Media(id: 'm1', description: 'Track A', mimetype: 'video/mp4', knownMediaId: uuidx.min()),
              ),
              media.RecentRecordRequest(
                id: '2',
                media: media.Media(id: 'm2', description: 'Track B', mimetype: 'video/mp4', knownMediaId: uuidx.min()),
              ),
            ],
          );
        },
      ),
    );
    await tester.pumpN(5);

    expect(find.text('Continue Watching'), findsOneWidget);
    expect(find.byType(lib.KnownMediaRowDisplay), findsNWidgets(2));
    expect(tester.takeException(), isNull);
  });

  testWidgets('tapping a recent item plays it under its recorded query', (tester) async {
    final played = <(media.Media, media.MediaSearchResponse)>[];

    await tester.pumpApp(
      Empty(
        remote.Sync(token: 'bearer test'),
        search: ValueNotifier(media.MediaSearchState(next: media.MediaSearchRequest())),
        socket: remote.RemoteControlSocket.noop,
        sessionID: uuidx.v7(),
        onPlay: (context, m, resp) {
          played.add((m, resp));
          return null;
        },
        apisearch: (req, {String? host, List<httpx.Option> options = const []}) async {
          return media.MediaSearchResponse(next: req);
        },
        apirecentlatest: (req, {String? host, List<httpx.Option> options = const []}) async {
          return media.RecentSearchResponse(
            items: [
              media.RecentRecordRequest(
                id: '1',
                media: media.Media(id: 'm1', description: 'Track A', mimetype: 'video/mp4', knownMediaId: uuidx.min()),
                query: media.MediaSearchRequest(query: 'my query'),
              ),
            ],
          );
        },
      ),
    );
    await tester.pumpN(5);

    await tester.widget<lib.KnownMediaRowDisplay>(find.byType(lib.KnownMediaRowDisplay)).onTap!();
    await tester.pumpN(2);

    expect(played, hasLength(1));
    expect(played.single.$1.id, 'm1');
    expect(played.single.$2.next.query, 'my query');
    expect(tester.takeException(), isNull);
  });

  testWidgets('shows loading instead of recent items while the sync has no token', (tester) async {
    await tester.pumpApp(
      Empty(
        remote.Sync(),
        search: ValueNotifier(media.MediaSearchState(next: media.MediaSearchRequest())),
        socket: remote.RemoteControlSocket.noop,
        sessionID: uuidx.v7(),
        onPlay: (context, m, resp) => null,
        apisearch: (req, {String? host, List<httpx.Option> options = const []}) async {
          return media.MediaSearchResponse(next: req);
        },
        apirecentlatest: (req, {String? host, List<httpx.Option> options = const []}) async {
          return media.RecentSearchResponse();
        },
      ),
    );
    await tester.pumpN(5);

    expect(find.byType(disc.RecentList), findsNothing);
    expect(tester.takeException(), isNull);
  });

  testWidgets('shows library search when there is no recent history', (tester) async {
    await tester.pumpApp(
      Empty(
        remote.Sync(token: 'bearer test'),
        search: ValueNotifier(media.MediaSearchState(next: media.MediaSearchRequest())),
        socket: remote.RemoteControlSocket.noop,
        sessionID: uuidx.v7(),
        onPlay: (context, m, resp) => null,
        apisearch: (req, {String? host, List<httpx.Option> options = const []}) async {
          return media.MediaSearchResponse(next: req);
        },
        apirecentlatest: (req, {String? host, List<httpx.Option> options = const []}) async {
          return media.RecentSearchResponse();
        },
      ),
    );
    await tester.pumpN(5);

    expect(find.byType(lib.SearchMinimal), findsOneWidget);
    expect(find.text('Media you watch will appear here'), findsNothing);
    expect(tester.takeException(), isNull);
  });
}
