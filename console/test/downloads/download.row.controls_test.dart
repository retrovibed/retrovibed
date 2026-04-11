import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/media/download.row.controls.dart';
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/timex.dart' as timex;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

media.Download _download({required String completedAt}) {
  return media.Download(
    completedAt: completedAt,
    media: media.Media(
      id: 'test-id',
      description: 'Test Media',
      mimetype: 'video/mp4',
      createdAt: '2025-01-01T00:00:00Z',
      archiveId: uuidx.min(),
      torrentId: uuidx.min(),
      knownMediaId: uuidx.min(),
    ),
  );
}

void main() {
  group('DownloadRowControls completedAt', () {
    testWidgets('shows check icon when completedAt is a past timestamp', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        DownloadRowControls(current: _download(completedAt: '2025-06-01T12:00:00Z')),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
      expect(find.byIcon(Icons.check), findsOneWidget);
      expect(find.byIcon(Icons.pause_circle_outline), findsNothing);
    });

    testWidgets('shows pause icon when completedAt is timex.inf', (
      WidgetTester tester,
    ) async {
      final infString = timex.formatISO8601(timex.inf);
      await tester.pumpApp(
        DownloadRowControls(current: _download(completedAt: infString)),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
      expect(find.byIcon(Icons.pause_circle_outline), findsOneWidget);
      expect(find.byIcon(Icons.check), findsNothing);
    });

    testWidgets('renders without overflow when completed', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        DownloadRowControls(current: _download(completedAt: '2025-06-01T12:00:00Z')),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without overflow when in progress', (
      WidgetTester tester,
    ) async {
      final infString = timex.formatISO8601(timex.inf);
      await tester.pumpApp(
        DownloadRowControls(current: _download(completedAt: infString)),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    });
  });
}
