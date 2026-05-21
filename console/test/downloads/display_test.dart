import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/downloads/display.dart';
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

Future<Stream<media.Download>> _mockWatch(
  String id, {
  List<httpx.Option> options = const [], // ignore: avoid_unused_parameters
}) async {
  return StreamController<media.Download>().stream;
}

Future<media.DownloadSearchResponse> _mockSearchWithItems(
  media.DownloadSearchRequest req, {
  List<httpx.Option> options = const [],
}) async {
  return media.DownloadSearchResponse(
    items: [
      media.Download(
        media: media.Media(
          id: 'test-id-1',
          description: 'Test Media One with a fairly long description',
          mimetype: 'video/mp4',
          createdAt: '2025-01-01T00:00:00Z',
          archiveId: uuidx.min(),
          torrentId: uuidx.min(),
          knownMediaId: uuidx.min(),
        ),
      ),
      media.Download(
        media: media.Media(
          id: 'test-id-2',
          description: 'Test Media Two',
          mimetype: 'audio/mp3',
          createdAt: '2025-01-02T00:00:00Z',
          archiveId: uuidx.min(),
          torrentId: uuidx.min(),
          knownMediaId: uuidx.min(),
        ),
      ),
    ],
    next: media.discoveredsearch.request(limit: 32),
  );
}

Future<media.DownloadSearchResponse> _mockSearchWithLongNames(
  media.DownloadSearchRequest req, {
  List<httpx.Option> options = const [],
}) async {
  return media.DownloadSearchResponse(
    items: [
      media.Download(
        media: media.Media(
          id: 'very-long-identifier-that-could-potentially-cause-overflow',
          description:
              'An Extremely Long Title That Could Potentially Overflow The Display Widget And Cause Layout Issues',
          mimetype: 'video/mp4',
          createdAt: '2025-01-01T00:00:00Z',
          archiveId: uuidx.min(),
          torrentId: uuidx.min(),
          knownMediaId: uuidx.min(),
        ),
      ),
    ],
    next: media.discoveredsearch.request(limit: 32),
  );
}

Future<media.DownloadSearchResponse> _mockSearchWith100Items(
  media.DownloadSearchRequest req, {
  List<httpx.Option> options = const [],
}) async {
  return media.DownloadSearchResponse(
    items: List.generate(
      100,
      (i) => media.Download(
        media: media.Media(
          id: 'test-id-$i',
          description: 'Test Media $i',
          mimetype: 'video/mp4',
          createdAt: '2025-01-01T00:00:00Z',
          archiveId: uuidx.min(),
          torrentId: uuidx.min(),
          knownMediaId: uuidx.min(),
        ),
      ),
    ),
    next: media.discoveredsearch.request(limit: 32),
  );
}

final _resolutions = Resolutions.variant();

void main() {
  group('Display', () {
    testWidgets('renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        Display(downloadingSearch: _mockSearchWithItems, availableSearch: _mockSearchWithItems, downloadWatch: _mockWatch),
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
        Display(downloadingSearch: _mockSearchWithLongNames, availableSearch: _mockSearchWithLongNames, downloadWatch: _mockWatch),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets(
      'renders without overflow with 100 results from both searches',
      (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          Display(
            downloadingSearch: _mockSearchWith100Items,
            availableSearch: _mockSearchWith100Items,
            downloadWatch: _mockWatch,
          ),
          physicalSize: entry.value,
        );
        await tester.pumpAndSettle();
        expect(tester.takeException(), isNull);
      },
      variant: _resolutions,
    );

    testWidgets('renders without overflow at narrow 300x600', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Display(downloadingSearch: _mockSearchWithItems, availableSearch: _mockSearchWithItems, downloadWatch: _mockWatch),
        physicalSize: const Size(300, 600),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders with long names without overflow at narrow 300x600', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Display(downloadingSearch: _mockSearchWithLongNames, availableSearch: _mockSearchWithLongNames, downloadWatch: _mockWatch),
        physicalSize: const Size(300, 600),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    });

    group('tuning panel', () {
      Future<void> openTuning(WidgetTester tester) async {
        final moreVert = find.byIcon(Icons.more_vert);
        if (moreVert.evaluate().isNotEmpty) {
          await tester.tap(moreVert.first);
          await tester.pumpAndSettle();
        }
        await tester.tap(find.byIcon(Icons.tune).first);
        await tester.pumpAndSettle();
      }

      testWidgets('opens', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          Display(downloadingSearch: _mockSearchWithItems, availableSearch: _mockSearchWithItems, downloadWatch: _mockWatch),
          physicalSize: entry.value,
        );
        await tester.pumpAndSettle();
        await openTuning(tester);
      }, variant: _resolutions);

      testWidgets('opens at narrow 300x600', (WidgetTester tester) async {
        await tester.pumpApp(
          Display(downloadingSearch: _mockSearchWithItems, availableSearch: _mockSearchWithItems, downloadWatch: _mockWatch),
          physicalSize: const Size(300, 600),
        );
        await tester.pumpAndSettle();
        await openTuning(tester);
      });

      testWidgets('opens and closes at narrow 300x600', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Display(downloadingSearch: _mockSearchWithItems, availableSearch: _mockSearchWithItems, downloadWatch: _mockWatch),
          physicalSize: const Size(300, 600),
        );
        await tester.pumpAndSettle();
        await openTuning(tester);
        final tuneButton = find.byIcon(Icons.tune).first;
        await tester.ensureVisible(tuneButton);
        await tester.pumpAndSettle();
        await tester.tap(tuneButton);
        await tester.pumpAndSettle();
      });
    });
  });
}
