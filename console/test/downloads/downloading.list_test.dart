import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/downloads/downloading.list.dart';
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/timex.dart' as timex;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

Future<media.DownloadSearchResponse> _mockSearchEmpty(
  media.DownloadSearchRequest req, {
  List<httpx.Option> options = const [],
}) async {
  return media.discoveredsearch.response(
    next: media.discoveredsearch.request(limit: 32),
  );
}

final _resolutions = Resolutions.variant();

void main() {
  group('DownloadingListDisplay', () {
    testWidgets('takes up zero height when empty', (WidgetTester tester) async {
      await tester.pumpApp(
        MediaQuery(
          data: const MediaQueryData(
            padding: EdgeInsets.only(top: 24, bottom: 34),
          ),
          child: DownloadingListDisplay(search: _mockSearchEmpty),
        ),
      );
      await tester.pumpAndSettle();
      final size = tester.getSize(find.byType(DownloadingListDisplay));
      expect(size.height, equals(0.0));
    });

    testWidgets('refreshes the list when a watched download completes', (WidgetTester tester) async {
      final StreamController<media.Download> stream = StreamController<media.Download>();
      final media.Download inprogress = media.Download(
        completedAt: timex.formatISO8601(timex.inf),
        media: media.Media(
          id: 'a',
          description: 'First Download',
          mimetype: 'video/mp4',
          createdAt: '2025-01-01T00:00:00Z',
          archiveId: uuidx.min(),
          torrentId: uuidx.min(),
          knownMediaId: uuidx.min(),
        ),
      );
      int searches = 0;

      await tester.pumpApp(
        DownloadingListDisplay(
          search: (media.DownloadSearchRequest req, {List<httpx.Option> options = const []}) async {
            searches++;
            return media.discoveredsearch.response(next: req)..items.add(inprogress);
          },
          watch: (String id, {List<httpx.Option> options = const []}) async => stream.stream,
        ),
      );
      await tester.pumpAndSettle();
      expect(searches, equals(1));

      stream.add(inprogress.deepCopy()..completedAt = '2025-06-01T12:00:00Z');
      await tester.pumpAndSettle();

      expect(searches, equals(2));
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        DownloadingListDisplay(search: _mockSearchEmpty),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);
  });
}
