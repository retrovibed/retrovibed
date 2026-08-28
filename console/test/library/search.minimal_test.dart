import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/library/search.minimal.dart';
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

Future<media.MediaSearchResponse> _mockSearchEmpty(
  media.MediaSearchRequest req, {
  String? host,
  List<httpx.Option> options = const [],
}) async {
  return media.media.response(
    next: media.media.request(limit: 32),
  );
}

Future<media.MediaSearchResponse> _mockSearchWithItems(
  media.MediaSearchRequest req, {
  String? host,
  List<httpx.Option> options = const [],
}) async {
  return media.MediaSearchResponse(
    items: [
      media.Media(
        id: uuidx.withSuffix(1),
        description: 'Test Media One with a fairly long description',
        mimetype: 'video/mp4',
        createdAt: '2025-01-01T00:00:00Z',
        archiveId: uuidx.min(),
        torrentId: uuidx.min(),
        knownMediaId: uuidx.min(),
      ),
      media.Media(
        id: uuidx.withSuffix(2),
        description: 'Test Media Two',
        mimetype: 'audio/mp3',
        createdAt: '2025-01-02T00:00:00Z',
        archiveId: uuidx.min(),
        torrentId: uuidx.min(),
        knownMediaId: uuidx.min(),
      ),
    ],
    next: media.media.request(limit: 32),
  );
}

Future<media.MediaSearchResponse> _mockSearchWithLongNames(
  media.MediaSearchRequest req, {
  String? host,
  List<httpx.Option> options = const [],
}) async {
  return media.MediaSearchResponse(
    items: [
      media.Media(
        id: uuidx.withSuffix(1),
        description:
            'An Extremely Long Title That Could Potentially Overflow The Display Widget And Cause Layout Issues In The List',
        mimetype: 'video/mp4',
        createdAt: '2025-01-01T00:00:00Z',
        archiveId: uuidx.min(),
        torrentId: uuidx.min(),
        knownMediaId: uuidx.min(),
      ),
    ],
    next: media.media.request(limit: 32),
  );
}

final _resolutions = Resolutions.variant();

ValueNotifier<media.MediaSearchState> _search() => ValueNotifier(
  media.MediaSearchState(next: media.media.request(limit: 32)),
);

void main() {
  group('SearchMinimal', () {
    testWidgets('renders empty without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        SearchMinimal(apisearch: _mockSearchEmpty, search: _search()),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders with items without overflow', (
      WidgetTester tester,
    ) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        SearchMinimal(apisearch: _mockSearchWithItems, search: _search()),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders with long names without overflow', (
      WidgetTester tester,
    ) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        SearchMinimal(apisearch: _mockSearchWithLongNames, search: _search()),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders empty without overflow at narrow 300x600', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SearchMinimal(apisearch: _mockSearchEmpty, search: _search()),
        physicalSize: const Size(300, 600),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders with items without overflow at narrow 300x600', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SearchMinimal(apisearch: _mockSearchWithItems, search: _search()),
        physicalSize: const Size(300, 600),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    });
  });
}
