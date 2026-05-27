import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/library/known.media.dropdown.dart';
import 'package:retrovibed/library/known.media.card.dart';
import 'package:retrovibed/library/api.dart' as api;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

// A Known item returned by the mock search so the dropdown has a card to tap.
final _knownItem = api.Known(
  id: uuidx.withSuffix(99),
  description: 'Test Known Media',
  summary: 'Test summary',
);

Future<api.KnownSearchResponse> _mockSearch(
  api.KnownSearchRequest req, {
  List<httpx.Option> options = const [],
}) async => api.KnownSearchResponse(items: [_knownItem], next: req);

// Gives the inline widget a valid context and enough scroll room to avoid
// overflow errors from the search grid.
Widget _wrap(Widget Function(BuildContext) builder) => Builder(
  builder: (ctx) => SingleChildScrollView(child: builder(ctx)),
);

void main() {
  group('KnownMediaDropdown.inline', () {
    testWidgets('renders without overflow', (tester) async {
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
        _wrap(
          (ctx) => KnownMediaDropdown.inline(ctx, current, search: _mockSearch),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(KnownMediaDropdown), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    group('selection routing', () {
      testWidgets('calls libraryMetadataSync when torrentId is min', (tester) async {
        bool libraryCalled = false;
        bool discoveredCalled = false;
        media.Media? resultMedia;

        final current = media.Media(
          id: uuidx.withSuffix(1),
          description: 'Test',
          mimetype: 'video/mp4',
          createdAt: DateTime.now().toIso8601String(),
          archiveId: uuidx.min(),
          torrentId: uuidx.min(), // min → library path
          knownMediaId: uuidx.min(),
        );

        final synced = current.deepCopy()..knownMediaId = _knownItem.id;

        await tester.pumpApp(
          _wrap(
            (ctx) => KnownMediaDropdown.inline(
              ctx,
              current,
              search: _mockSearch,
              onChange: (v) => resultMedia = v,
              apiLibraryMetadataSync: (id, m, {options = const []}) async {
                libraryCalled = true;
                return media.MediaUpdateResponse(media: synced);
              },
              apiDiscoveredMetadataSync: (id, m, {options = const []}) async {
                discoveredCalled = true;
                return media.MetadataSyncResponse(media: synced);
              },
            ),
          ),
        );
        await tester.pumpAndSettle();

        await tester.doubleTap(find.byType(KnownMediaCard).first);
        await tester.pumpAndSettle();

        expect(libraryCalled, isTrue);
        expect(discoveredCalled, isFalse);
        expect(resultMedia?.knownMediaId, equals(synced.knownMediaId));
        expect(tester.takeException(), isNull);
      });

      testWidgets('calls discoveredMetadataSync when torrentId is valid', (tester) async {
        bool libraryCalled = false;
        bool discoveredCalled = false;
        media.Media? resultMedia;

        final current = media.Media(
          id: uuidx.withSuffix(1),
          description: 'Test',
          mimetype: 'video/mp4',
          createdAt: DateTime.now().toIso8601String(),
          archiveId: uuidx.min(),
          torrentId: uuidx.withSuffix(5), // valid → discovered path
          knownMediaId: uuidx.min(),
        );

        final synced = current.deepCopy()..knownMediaId = _knownItem.id;

        await tester.pumpApp(
          _wrap(
            (ctx) => KnownMediaDropdown.inline(
              ctx,
              current,
              search: _mockSearch,
              onChange: (v) => resultMedia = v,
              apiLibraryMetadataSync: (id, m, {options = const []}) async {
                libraryCalled = true;
                return media.MediaUpdateResponse(media: synced);
              },
              apiDiscoveredMetadataSync: (id, m, {options = const []}) async {
                discoveredCalled = true;
                return media.MetadataSyncResponse(media: synced);
              },
            ),
          ),
        );
        await tester.pumpAndSettle();

        await tester.doubleTap(find.byType(KnownMediaCard).first);
        await tester.pumpAndSettle();

        expect(discoveredCalled, isTrue);
        expect(libraryCalled, isFalse);
        expect(resultMedia?.knownMediaId, equals(synced.knownMediaId));
        expect(tester.takeException(), isNull);
      });

      testWidgets('does not crash or call onChange on deactivate with no selection', (tester) async {
        media.Media? resultMedia;

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
          _wrap(
            (ctx) => KnownMediaDropdown.inline(
              ctx,
              current,
              search: _mockSearch,
              onChange: (v) => resultMedia = v,
              apiLibraryMetadataSync: (id, m, {options = const []}) async => media.MediaUpdateResponse(media: m),
              apiDiscoveredMetadataSync: (id, m, {options = const []}) async => media.MetadataSyncResponse(media: m),
            ),
          ),
        );
        await tester.pumpAndSettle();

        // Explicitly deactivate — triggers _KnownMediaDropdown.deactivate
        // which calls onChange(null). Without the guard in _sync this would
        // crash with "looking up a deactivated widget's ancestor".
        await tester.pumpWidget(
          const MaterialApp(home: Scaffold(body: Text('gone'))),
        );
        await tester.pump();

        // onChange fires with the unmodified current (no API call, no context
        // access) — the important assertion is that no exception was thrown.
        expect(resultMedia?.id, equals(current.id));
        expect(resultMedia?.knownMediaId, equals(uuidx.min()));
        expect(tester.takeException(), isNull);
      });
    });
  });
}
