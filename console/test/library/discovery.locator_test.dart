import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/library/discovery.locator.dart';
import 'package:retrovibed/library/known.media.locator.dart';
import 'package:retrovibed/library/known.media.card.dart';
import 'package:retrovibed/library/api.dart' as api;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('DiscoveryLocator', () {
    testWidgets('renders nothing when the query is blank', (tester) async {
      await tester.pumpApp(
        DiscoveryLocator(query: '  ', mimetype: 'video', onFound: (located) async => const SizedBox.shrink()),
      );
      await tester.pumpAndSettle();
      expect(find.byType(ds.Card), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without overflow', (tester) async {
      await tester.pumpApp(
        DiscoveryLocator(query: 'ubuntu', mimetype: 'video', onFound: (located) async => const SizedBox.shrink()),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    });

    testWidgets('declining the p2p prompt never calls locate', (tester) async {
      bool locateCalled = false;
      await tester.pumpApp(
        DiscoveryLocator(
          query: 'ubuntu',
          mimetype: 'video',
          onFound: (located) async => const SizedBox.shrink(),
          ensureP2P: (context, {options = const []}) async => false,
          locate: (req, {options = const []}) async {
            locateCalled = true;
            return api.LocateCreateResponse(locate: req);
          },
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(ds.Card));
      await tester.pumpAndSettle();

      expect(locateCalled, isFalse);
      expect(tester.takeException(), isNull);
      expect(find.byIcon(Icons.travel_explore_rounded), findsOneWidget);
    });

    testWidgets('tapping locates with the query and mimetype and switches to pending', (tester) async {
      api.Locate? requested;
      await tester.pumpApp(
        DiscoveryLocator(
          query: 'ubuntu',
          mimetype: 'video',
          onFound: (located) async => const SizedBox.shrink(),
          ensureP2P: (context, {options = const []}) async => true,
          locate: (req, {options = const []}) async {
            requested = req;
            return api.LocateCreateResponse(locate: (req..id = 'locate-1'));
          },
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(ds.Card));
      await tester.pumpAndSettle();

      expect(requested?.query, equals('ubuntu'));
      expect(requested?.mimetype, equals('video'));
      expect(tester.takeException(), isNull);
      expect(find.byIcon(Icons.query_builder_rounded), findsOneWidget);
      expect(find.byType(ds.Card), findsOneWidget);
    });

    testWidgets('a failed locate is handled internally', (tester) async {
      await tester.pumpApp(
        DiscoveryLocator(
          query: 'ubuntu',
          mimetype: 'video',
          onFound: (located) async => const SizedBox.shrink(),
          ensureP2P: (context, {options = const []}) async => true,
          locate: (req, {options = const []}) => Future.error('boom'),
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(ds.Card));
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('polls lookup every 10 seconds while pending', (tester) async {
      int lookupCalls = 0;
      await tester.pumpApp(
        DiscoveryLocator(
          query: 'ubuntu',
          mimetype: 'video',
          onFound: (located) async => const SizedBox.shrink(),
          ensureP2P: (context, {options = const []}) async => true,
          locate: (req, {options = const []}) async => api.LocateCreateResponse(locate: (req..id = 'locate-1')),
          lookup: (id, {options = const []}) async {
            lookupCalls++;
            return api.LocateLookupResponse(locate: api.Locate(id: id));
          },
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(ds.Card));
      await tester.pumpAndSettle();

      expect(lookupCalls, equals(0));

      await tester.pump(const Duration(seconds: 10));
      await tester.pump();
      expect(lookupCalls, equals(1));

      await tester.pump(const Duration(seconds: 10));
      await tester.pump();
      expect(lookupCalls, equals(2));

      expect(tester.takeException(), isNull);
    });

    testWidgets('stops polling and shows the found media with a download action once located', (tester) async {
      int lookupCalls = 0;
      api.Locate? foundWith;
      final found = api.Known(id: 'rec-1', uid: 'torrent-1', description: 'Ubuntu', summary: 'summary');
      await tester.pumpApp(
        DiscoveryLocator(
          query: 'ubuntu',
          mimetype: 'video',
          onFound: (located) async {
            foundWith = located;
            return ConstrainedBox(
              constraints: BoxConstraints.tightForFinite(height: 256),
              child: KnownMediaLocator.future(Future.value(found)),
            );
          },
          ensureP2P: (context, {options = const []}) async => true,
          locate: (req, {options = const []}) async => api.LocateCreateResponse(locate: (req..id = 'locate-1')),
          lookup: (id, {options = const []}) async {
            lookupCalls++;
            return api.LocateLookupResponse(
              locate: api.Locate(id: id, locatedTorrentId: lookupCalls >= 2 ? 'torrent-1' : ''),
            );
          },
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(ds.Card));
      await tester.pumpAndSettle();

      await tester.pump(const Duration(seconds: 10));
      await tester.pump();
      expect(lookupCalls, equals(1));
      expect(find.byIcon(Icons.query_builder_rounded), findsOneWidget);

      await tester.pump(const Duration(seconds: 10));
      await tester.pump();
      await tester.pumpAndSettle();
      expect(lookupCalls, equals(2));
      expect(foundWith?.locatedTorrentId, equals('torrent-1'));
      expect(find.byType(KnownMediaLocator), findsOneWidget);
      expect(find.byType(KnownMediaCard), findsOneWidget);
      expect(find.text('found — added to recommendations'), findsNothing);
      expect(find.textContaining('recommendations'), findsNothing);

      await tester.pump(const Duration(seconds: 10));
      await tester.pump();
      expect(lookupCalls, equals(2), reason: 'polling should stop once found');

      expect(tester.takeException(), isNull);
    });

    testWidgets('polling stops once found even if onFound renders an error state', (tester) async {
      int lookupCalls = 0;
      await tester.pumpApp(
        DiscoveryLocator(
          query: 'ubuntu',
          mimetype: 'video',
          onFound: (located) async => Text('failed to load ${located.locatedTorrentId}'),
          ensureP2P: (context, {options = const []}) async => true,
          locate: (req, {options = const []}) async => api.LocateCreateResponse(locate: (req..id = 'locate-1')),
          lookup: (id, {options = const []}) async {
            lookupCalls++;
            return api.LocateLookupResponse(locate: api.Locate(id: id, locatedTorrentId: 'torrent-1'));
          },
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(ds.Card));
      await tester.pumpAndSettle();

      await tester.pump(const Duration(seconds: 10));
      await tester.pump();
      expect(lookupCalls, equals(1));
      expect(find.text('failed to load torrent-1'), findsOneWidget);

      await tester.pump(const Duration(seconds: 10));
      await tester.pump();
      expect(lookupCalls, equals(1), reason: 'polling should stop once found, regardless of onFound content');

      expect(tester.takeException(), isNull);
    });

    testWidgets('poll lookup errors are swallowed and polling retries', (tester) async {
      int lookupCalls = 0;
      await tester.pumpApp(
        DiscoveryLocator(
          query: 'ubuntu',
          mimetype: 'video',
          onFound: (located) async => const SizedBox.shrink(),
          ensureP2P: (context, {options = const []}) async => true,
          locate: (req, {options = const []}) async => api.LocateCreateResponse(locate: (req..id = 'locate-1')),
          lookup: (id, {options = const []}) {
            lookupCalls++;
            return Future.error('boom');
          },
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(ds.Card));
      await tester.pumpAndSettle();

      await tester.pump(const Duration(seconds: 10));
      await tester.pump();
      expect(lookupCalls, equals(1));

      await tester.pump(const Duration(seconds: 10));
      await tester.pump();
      expect(lookupCalls, equals(2));

      expect(tester.takeException(), isNull);
      expect(find.byType(ds.Card), findsOneWidget);
    });
  });
}
