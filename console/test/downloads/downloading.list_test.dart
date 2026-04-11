import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/downloads/downloading.list.dart';
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/httpx.dart' as httpx;
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
