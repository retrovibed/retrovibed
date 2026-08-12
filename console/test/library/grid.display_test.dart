import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/library/grid.display.dart';
import 'package:retrovibed/library/known.media.display.dart';
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

media.Media _testMedia({required String id, required String description}) {
  return media.Media(
    id: id,
    description: description,
    mimetype: 'video/mp4',
    createdAt: '2025-01-01T00:00:00Z',
    archiveId: uuidx.min(),
    torrentId: uuidx.min(),
    knownMediaId: uuidx.min(),
  );
}

final _resolutions = Resolutions.variant();

void main() {
  group('Grid', () {
    testWidgets('renders without overflow when empty', (
      WidgetTester tester,
    ) async {
      final entry = _resolutions.currentValue!;
      final search = ValueNotifier<media.MediaSearchState>(
        media.MediaSearchState(next: media.media.request(limit: 32)),
      );
      addTearDown(search.dispose);
      await tester.pumpApp(
        Grid(
          apisearch: (req, {host, options = const []}) async {
            return media.MediaSearchResponse(items: [], next: req);
          },
          search: search,
          empty: const SizedBox.shrink(),
          highlighted: '',
        ),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders with items without overflow', (
      WidgetTester tester,
    ) async {
      final entry = _resolutions.currentValue!;
      final search = ValueNotifier<media.MediaSearchState>(
        media.MediaSearchState(next: media.media.request(limit: 32)),
      );
      addTearDown(search.dispose);
      await tester.pumpApp(
        Grid(
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
          empty: const SizedBox.shrink(),
          highlighted: '',
        ),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders with long names without overflow', (
      WidgetTester tester,
    ) async {
      final entry = _resolutions.currentValue!;
      final search = ValueNotifier<media.MediaSearchState>(
        media.MediaSearchState(next: media.media.request(limit: 32)),
      );
      addTearDown(search.dispose);
      await tester.pumpApp(
        Grid(
          apisearch: (req, {host, options = const []}) async {
            return media.MediaSearchResponse(
              items: [
                media.Media(
                  id: uuidx.withSuffix(1),
                  description:
                      'An Extremely Long Title That Could Potentially Overflow The Display Widget And Cause Layout Issues In The Grid',
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
          empty: const SizedBox.shrink(),
          highlighted: '',
        ),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders with highlighted item without overflow', (
      WidgetTester tester,
    ) async {
      final entry = _resolutions.currentValue!;
      final search = ValueNotifier<media.MediaSearchState>(
        media.MediaSearchState(next: media.media.request(limit: 32)),
      );
      addTearDown(search.dispose);
      await tester.pumpApp(
        Grid(
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
          empty: const SizedBox.shrink(),
          highlighted: uuidx.withSuffix(1),
        ),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('onChange replaces the matching item in place without touching others', (
      WidgetTester tester,
    ) async {
      final search = ValueNotifier<media.MediaSearchState>(
        media.MediaSearchState(next: media.media.request(limit: 32)),
      );
      addTearDown(search.dispose);
      await tester.pumpApp(
        Grid(
          apisearch: (req, {host, options = const []}) async {
            return media.MediaSearchResponse(
              items: [
                _testMedia(id: 'test-id-1', description: 'Media One'),
                _testMedia(id: 'test-id-2', description: 'Media Two'),
              ],
              next: req,
            );
          },
          search: search,
          empty: const SizedBox.shrink(),
          highlighted: '',
        ),
      );
      await tester.pumpAndSettle();

      final before = tester.widgetList<KnownMediaDisplay>(find.byType(KnownMediaDisplay)).toList();
      expect(before.length, 2);
      final target = before.firstWhere((d) => d.media.id == 'test-id-1');

      target.onChange!(_testMedia(id: 'test-id-1', description: 'Media One Updated'));
      await tester.pump();

      final after = tester.widgetList<KnownMediaDisplay>(find.byType(KnownMediaDisplay)).toList();
      expect(after.length, 2);
      expect(after.firstWhere((d) => d.media.id == 'test-id-1').media.description, 'Media One Updated');
      expect(after.firstWhere((d) => d.media.id == 'test-id-2').media.description, 'Media Two');
      expect(tester.takeException(), isNull);
    });
  });
}
