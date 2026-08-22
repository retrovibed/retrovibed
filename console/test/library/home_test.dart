import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/discovery.dart' as disc;
import 'package:retrovibed/library/dropdown.upload.dart';
import 'package:retrovibed/library/home.dart';
import 'package:retrovibed/library/known.media.display.dart';
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

IconData _searchMenuItemIcon(WidgetTester tester) {
  final icon = find.descendant(
    of: find.ancestor(of: find.text('Discover'), matching: find.byType(ListTile)),
    matching: find.byType(Icon),
  );
  return tester.widget<Icon>(icon).icon!;
}

final _resolutions = Resolutions.variant();

void main() {
  group('Home', () {
    testWidgets('hides library items once discover mode is activated', (
      WidgetTester tester,
    ) async {
      final entry = _resolutions.currentValue!;
      final search = ValueNotifier<media.MediaSearchState>(
        media.MediaSearchState(next: media.media.request(limit: 32)),
      );
      addTearDown(search.dispose);
      await tester.pumpApp(
        Home(
          apisearch: (req, {host, options = const []}) async {
            return media.MediaSearchResponse(
              items: [
                media.Media(
                  id: uuidx.withSuffix(1),
                  description: 'Test Media One',
                  mimetype: 'video/mp4',
                  createdAt: '2025-01-01T00:00:00Z',
                  archiveId: uuidx.min(),
                  torrentId: uuidx.min(),
                  knownMediaId: uuidx.min(),
                ),
              ],
              next: req,
            );
          },
          search: search,
          highlighted: '',
        ),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(find.byType(KnownMediaDisplay), findsOneWidget);

      await tester.tap(find.byType(DropdownUpload));
      await tester.pumpAndSettle();
      await tester.ensureVisible(find.text('Discover'));
      await tester.tap(find.text('Discover'));
      await tester.pumpAndSettle();

      expect(find.byType(KnownMediaDisplay), findsNothing);
    }, variant: _resolutions);

    testWidgets('keeps the discover-mode check icon after switching mimetype filters', (
      WidgetTester tester,
    ) async {
      final entry = _resolutions.currentValue!;
      final search = ValueNotifier<media.MediaSearchState>(
        media.MediaSearchState(next: media.media.request(limit: 32)),
      );
      addTearDown(search.dispose);
      await tester.pumpApp(
        Home(
          apisearch: (req, {host, options = const []}) async {
            return media.MediaSearchResponse(items: [], next: req);
          },
          search: search,
          highlighted: '',
        ),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byType(DropdownUpload));
      await tester.pumpAndSettle();
      await tester.ensureVisible(find.text('Discover'));
      await tester.tap(find.text('Discover'));
      await tester.pumpAndSettle();

      // Selecting the mode item closes the menu; reopen it to confirm the
      // check icon now reflects discovery mode.
      await tester.tap(find.byType(DropdownUpload));
      await tester.pumpAndSettle();
      await tester.ensureVisible(find.text('Discover'));
      expect(_searchMenuItemIcon(tester), equals(Icons.check));

      // Switching mimetype filters doesn't close the menu; the Discover item's
      // icon should still reflect discovery mode within the same open session.
      await tester.ensureVisible(find.text('Movies'));
      await tester.tap(find.text('Movies'));
      await tester.pumpAndSettle();
      expect(_searchMenuItemIcon(tester), equals(Icons.check));

      // Selecting the mode item again switches back to library mode and
      // closes the menu; reopen it to confirm the icon reverted.
      await tester.tap(find.text('Discover'));
      await tester.pumpAndSettle();

      await tester.tap(find.byType(DropdownUpload));
      await tester.pumpAndSettle();
      await tester.ensureVisible(find.text('Discover'));
      expect(_searchMenuItemIcon(tester), equals(Icons.travel_explore));
    }, variant: _resolutions);

    testWidgets('switches to discovery mode when the empty-results discover button is tapped', (
      WidgetTester tester,
    ) async {
      final entry = _resolutions.currentValue!;
      final search = ValueNotifier<media.MediaSearchState>(
        media.MediaSearchState(next: media.media.request(limit: 32, query: 'something')),
      );
      addTearDown(search.dispose);
      await tester.pumpApp(
        Home(
          apisearch: (req, {host, options = const []}) async {
            return media.MediaSearchResponse(items: [], next: req);
          },
          search: search,
          highlighted: '',
        ),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();

      expect(find.widgetWithText(ElevatedButton, 'discover'), findsOneWidget);
      expect(find.byType(disc.DiscoveryGrid), findsNothing);

      await tester.tap(find.widgetWithText(ElevatedButton, 'discover'));
      await tester.pumpAndSettle();

      expect(find.byType(disc.DiscoveryGrid), findsOneWidget);

      await tester.tap(find.byType(DropdownUpload));
      await tester.pumpAndSettle();
      await tester.ensureVisible(find.text('Discover'));
      expect(_searchMenuItemIcon(tester), equals(Icons.check));

      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('shows the search tray discover button only in discovery mode', (
      WidgetTester tester,
    ) async {
      final entry = _resolutions.currentValue!;
      final search = ValueNotifier<media.MediaSearchState>(
        media.MediaSearchState(
          next: media.media.request(limit: 32, mimetypes: mimex.of(mimex.icomovie))..query = 'something',
        ),
      );
      addTearDown(search.dispose);
      await tester.pumpApp(
        Home(
          apisearch: (req, {host, options = const []}) async {
            return media.MediaSearchResponse(items: [], next: req);
          },
          search: search,
          highlighted: '',
        ),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();

      final traySearchButton = find.descendant(
        of: find.byType(ds.SearchTray),
        matching: find.byType(disc.SearchButton),
      );

      expect(traySearchButton, findsNothing);

      await tester.tap(find.byType(DropdownUpload));
      await tester.pumpAndSettle();
      await tester.ensureVisible(find.text('Discover'));
      await tester.tap(find.text('Discover'));
      await tester.pumpAndSettle();

      // On compact layouts unpinned trailing widgets collapse into the
      // search tray's overflow menu until it's opened.
      final traySearchOverflow = find.descendant(
        of: find.byType(ds.SearchTray),
        matching: find.byIcon(Icons.more_vert),
      );
      if (traySearchOverflow.evaluate().isNotEmpty) {
        await tester.tap(traySearchOverflow);
        await tester.pumpAndSettle();
      }

      expect(traySearchButton, findsOneWidget);

      await tester.tap(find.byType(DropdownUpload));
      await tester.pumpAndSettle();
      await tester.ensureVisible(find.text('Discover'));
      await tester.tap(find.text('Discover'));
      await tester.pumpAndSettle();

      expect(traySearchButton, findsNothing);
    }, variant: _resolutions);
  });
}
