import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/library/list.display.dart';
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:http/http.dart' as http;

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

Future<media.MediaUploadResponse> _mockUpload(
  http.MultipartRequest Function(http.MultipartRequest req) mkreq,
) async {
  return media.MediaUploadResponse(
    media: media.Media(
      id: uuidx.withSuffix(99),
      description: 'Uploaded Media',
      mimetype: 'video/mp4',
      createdAt: '2025-01-01T00:00:00Z',
      archiveId: uuidx.min(),
      torrentId: uuidx.min(),
      knownMediaId: uuidx.min(),
    ),
  );
}

final _resolutions = Resolutions.variant();

void main() {
  group('AvailableListDisplay', () {
    testWidgets('renders empty without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        AvailableListDisplay(search: _mockSearchEmpty, upload: _mockUpload),
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
        AvailableListDisplay(search: _mockSearchWithItems, upload: _mockUpload),
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
        AvailableListDisplay(
          search: _mockSearchWithLongNames,
          upload: _mockUpload,
        ),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders with custom row builder without overflow', (
      WidgetTester tester,
    ) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        AvailableListDisplay(
          search: _mockSearchWithItems,
          upload: _mockUpload,
          row: (v) => Text(v.description),
        ),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders with controller without overflow', (
      WidgetTester tester,
    ) async {
      final entry = _resolutions.currentValue!;
      final controller = TextEditingController(text: 'initial query');
      addTearDown(controller.dispose);
      await tester.pumpApp(
        AvailableListDisplay(
          search: _mockSearchWithItems,
          upload: _mockUpload,
          controller: controller,
        ),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders empty without overflow at narrow 300x600', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        AvailableListDisplay(search: _mockSearchEmpty, upload: _mockUpload),
        physicalSize: const Size(300, 600),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders with items without overflow at narrow 300x600', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        AvailableListDisplay(search: _mockSearchWithItems, upload: _mockUpload),
        physicalSize: const Size(300, 600),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    });
  });
}
