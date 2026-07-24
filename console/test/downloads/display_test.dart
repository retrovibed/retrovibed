import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/downloads/display.dart';
import 'package:retrovibed/downloads/grid.settings.dart';
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

// Mirrors the real backend (shallows/media/http.discovered.go), which caps
// its SQL `.Limit(...)` clause at the caller's requested limit (max 100) and
// echoes the decoded request back as `next`. A mock that ignores the request
// and always returns 100 items can't happen against the real server.
Future<media.DownloadSearchResponse> _mockSearchHonoringLimit(
  media.DownloadSearchRequest req, {
  List<httpx.Option> options = const [],
}) async {
  final limit = req.limit.toInt().clamp(0, 100);
  return media.DownloadSearchResponse(
    items: List.generate(
      limit,
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
    next: req,
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
        fit: FlexFit.tight,
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
        fit: FlexFit.tight,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets(
      'renders without overflow with a full page of results from both searches',
      (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          Display(
            downloadingSearch: _mockSearchHonoringLimit,
            availableSearch: _mockSearchHonoringLimit,
            downloadWatch: _mockWatch,
          ),
          physicalSize: entry.value,
          fit: FlexFit.tight,
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
        fit: FlexFit.tight,
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
        fit: FlexFit.tight,
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
          fit: FlexFit.tight,
        );
        await tester.pumpAndSettle();
        await openTuning(tester);
        expect(find.byType(GridSettings), findsOneWidget);
      }, variant: _resolutions);

      testWidgets('opens at narrow 300x600', (WidgetTester tester) async {
        await tester.pumpApp(
          Display(downloadingSearch: _mockSearchWithItems, availableSearch: _mockSearchWithItems, downloadWatch: _mockWatch),
          physicalSize: const Size(300, 600),
          fit: FlexFit.tight,
        );
        await tester.pumpAndSettle();
        await openTuning(tester);
        expect(find.byType(GridSettings), findsOneWidget);
      });

      testWidgets('opens and closes at narrow 300x600', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Display(downloadingSearch: _mockSearchWithItems, availableSearch: _mockSearchWithItems, downloadWatch: _mockWatch),
          physicalSize: const Size(300, 600),
          fit: FlexFit.tight,
        );
        await tester.pumpAndSettle();
        await openTuning(tester);
        expect(find.byType(GridSettings), findsOneWidget);

        await tester.tap(find.byIcon(Icons.tune).first);
        await tester.pumpAndSettle();
        expect(find.byType(GridSettings), findsNothing);
      });
    });
  });
}
