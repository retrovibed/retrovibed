import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/library/known.media.dropdown.dart';
import 'package:retrovibed/library/known.media.card.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/library/api.dart' as api;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

// A Known item returned by the mock search so the dropdown has a card to tap.
final _knownItem = api.Known(
  id: uuidx.withSuffix(98),
  uid: uuidx.withSuffix(99),
  description: 'Test Known Media',
  summary: 'Test summary',
);

Future<api.KnownSearchResponse> _mockSearch(
  api.KnownSearchRequest req, {
  List<httpx.Option> options = const [],
}) async => api.KnownSearchResponse(items: [_knownItem], next: req);

void main() {
  group('KnownMediaDropdown', () {
    testWidgets('renders without overflow', (tester) async {
      await tester.pumpApp(
        SingleChildScrollView(child: KnownMediaDropdown(search: _mockSearch)),
      );
      await tester.pumpAndSettle();

      expect(find.byType(KnownMediaDropdown), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    group('selection', () {
      testWidgets('calls onChange with the tapped known media', (tester) async {
        api.Known? selected;

        await tester.pumpApp(
          SingleChildScrollView(
            child: KnownMediaDropdown(
              search: _mockSearch,
              onChange: (k) async => selected = k,
            ),
          ),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.byType(KnownMediaCard).first);
        await tester.pumpAndSettle();

        expect(selected?.uid, equals(_knownItem.uid));
        expect(tester.takeException(), isNull);
      });

      testWidgets('does not call onChange on deactivate with no selection', (tester) async {
        bool called = false;
        api.Known? selected = _knownItem;

        await tester.pumpApp(
          SingleChildScrollView(
            child: KnownMediaDropdown(
              search: _mockSearch,
              onChange: (k) async {
                called = true;
                selected = k;
                return k;
              },
            ),
          ),
        );
        await tester.pumpAndSettle();

        // Explicitly deactivate — triggers _KnownMediaDropdown.deactivate
        // which has no loaded known media to clear.
        await tester.pumpWidget(
          const MaterialApp(home: Scaffold(body: Text('gone'))),
        );
        await tester.pump();

        expect(called, isFalse);
        expect(selected, equals(_knownItem));
        expect(tester.takeException(), isNull);
      });
    });

    group('modal', () {
      testWidgets('syncs the selected known media through the library endpoint', (tester) async {
        String? syncedId;
        String? syncedKnownId;
        media.Media? resultMedia;
        Future<void>? pending;

        final current = media.Media(
          id: uuidx.withSuffix(1),
          description: 'Test',
          mimetype: 'video/mp4',
          createdAt: DateTime.now().toIso8601String(),
          archiveId: uuidx.min(),
          torrentId: uuidx.withSuffix(5), // library media imported from a torrent
          knownMediaId: uuidx.min(),
        );

        await tester.pumpApp(
          ds.Node(
            Builder(
              builder: (ctx) => TextButton(
                onPressed: () => pending = KnownMediaDropdown.modal(
                  ctx,
                  current,
                  search: _mockSearch,
                  onChange: (v) => resultMedia = v,
                  libraryMetadataSync: (id, m, {options = const []}) async {
                    syncedId = id;
                    syncedKnownId = m.knownMediaId;
                    return media.MediaUpdateResponse(media: m.deepCopy());
                  },
                )(),
                child: const Text('identify'),
              ),
            ),
          ),
        );

        await tester.tap(find.text('identify'));
        await tester.pumpAndSettle();

        await tester.tap(find.byType(KnownMediaCard).first);
        await tester.pumpAndSettle();
        await pending;

        expect(syncedId, equals(current.id));
        expect(syncedKnownId, equals(_knownItem.uid));
        expect(resultMedia?.id, equals(current.id));
        expect(resultMedia?.knownMediaId, equals(_knownItem.uid));
        expect(find.byType(KnownMediaDropdown), findsNothing);
        expect(tester.takeException(), isNull);
      });

      testWidgets('dismissing with nothing selected skips the sync', (tester) async {
        bool synced = false;
        media.Media? resultMedia;
        Future<void>? pending;

        final current = media.Media(
          id: uuidx.withSuffix(1),
          description: 'Test',
          mimetype: 'video/mp4',
          createdAt: DateTime.now().toIso8601String(),
          archiveId: uuidx.min(),
          torrentId: uuidx.min(),
          knownMediaId: uuidx.min(),
        );

        await tester.pumpApp(
          ds.Node(
            Builder(
              builder: (ctx) => TextButton(
                onPressed: () => pending = KnownMediaDropdown.modal(
                  ctx,
                  current,
                  search: _mockSearch,
                  onChange: (v) => resultMedia = v,
                  libraryMetadataSync: (id, m, {options = const []}) async {
                    synced = true;
                    return media.MediaUpdateResponse(media: m);
                  },
                )(),
                child: const Text('identify'),
              ),
            ),
          ),
        );

        await tester.tap(find.text('identify'));
        await tester.pumpAndSettle();
        expect(find.byType(KnownMediaDropdown), findsOneWidget);

        tester.state<ds.NodeState>(find.byType(ds.Node)).reset();
        await tester.pumpAndSettle();
        await pending;

        expect(synced, isFalse);
        expect(resultMedia, isNull);
        expect(tester.takeException(), isNull);
      });
    });

    group('search failure handling', () {
      testWidgets('search error surfaces as an unknown-error state instead of crashing', (tester) async {
        api.KnownSearchRequest? capturedReq;

        await tester.pumpApp(
          SingleChildScrollView(
            child: KnownMediaDropdown(
              search: (req, {options = const []}) async {
                capturedReq = req;
                throw Exception('boom');
              },
            ),
          ),
        );
        await tester.pumpAndSettle();

        // Reproduces the reported crash report: the user starts a `released`
        // date-range filter but leaves it unfinished ("@released:"), then hits
        // the search button before completing it. SearchTray's search button
        // forwards this raw, unfinished text straight to KnownMediaDropdown's
        // onSubmitted, which sends it on as the search query.
        await tester.enterText(find.byType(TextField), '@released:');
        await tester.pumpAndSettle();

        await tester.tap(find.byIcon(Icons.search_rounded));
        await tester.pumpAndSettle();

        expect(capturedReq?.query, equals('@released:'));
        expect(find.text('an unexpected problem has occurred'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });
  });
}
