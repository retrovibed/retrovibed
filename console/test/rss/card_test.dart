import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/rss/list.searchable.dart';
import 'package:retrovibed/rss/card.dart' as rss;
import 'package:retrovibed/rss/api.dart' as api;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

Future<api.FeedSearchResponse> _mockSearch(
  api.FeedSearchRequest req, {
  List<httpx.Option> options = const [],
}) async {
  return api.FeedSearchResponse(
    next: api.FeedSearchRequest(
      query: '',
      offset: ds.Int64(0),
      limit: ds.Int64(10),
    ),
    items: [
      api.Feed(
        id: 'feed-1',
        description: 'Test Feed 1',
        url: 'https://example.com/feed1',
      ),
      api.Feed(
        id: 'feed-2',
        description: 'Test Feed 2',
        url: 'https://example.com/feed2',
      ),
    ],
  );
}

void main() {
  group('rss.Card', () {
    testWidgets('renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        rss.Card(onPressed: (_) {}),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);
  });

  group('rss.Card full overlay', () {
    testWidgets('renders without overflow when viewport exceeds minHeight', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        physicalSize: const Size(800, 300),
        ds.Masked(
          ds.Container(
            // alignment: Alignment.topCenter,
            // constraints: BoxConstraints(minHeight: 512),
            // margin: EdgeInsets.zero,
            ListSearchable(search: _mockSearch),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });
  });
}
