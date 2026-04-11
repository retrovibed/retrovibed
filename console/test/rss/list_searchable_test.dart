import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/rss/list.searchable.dart';
import 'package:retrovibed/rss/api.dart' as api;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  Future<api.FeedSearchResponse> mockSearch(
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
          id: 'test-feed-1',
          description: 'Test Feed 1',
          url: 'https://example.com/feed1',
        ),
        api.Feed(
          id: 'test-feed-2',
          description: 'Test Feed 2',
          url: 'https://example.com/feed2',
        ),
      ],
    );
  }

  group('ListSearchable Widget Tests', () {
    testWidgets('Table height should increase when row is expanded', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SingleChildScrollView(child: ListSearchable(search: mockSearch)),
      );

      // Wait for initial load
      await tester.pumpAndSettle();

      // Verify initial state - rows should be visible but not expanded
      expect(find.text('Test Feed 1'), findsOneWidget);
      expect(find.text('Test Feed 2'), findsOneWidget);

      // Verify that the expanded content is NOT initially visible
      expect(find.text('url'), findsNothing);
      expect(find.text('description'), findsNothing);

      // Get the initial height of the table
      final heightlookup = () async {
        await tester.pump();
        return tester
            .renderObject<RenderBox>(find.byType(ListSearchable))
            .size
            .height;
      };
      final initialHeight = await tester.runAsync(heightlookup) ?? -1;

      // Find the first feed row and tap it to expand
      final firstFeedRow = find.byType(ds.TableRow).first;
      await tester.tap(firstFeedRow);
      await tester.pumpAndSettle();

      // Verify that the expanded content is now visible
      expect(find.text('url'), findsOneWidget);
      expect(find.text('description'), findsOneWidget);

      // Get the height after expansion
      final expandedHeight = await tester.runAsync(heightlookup) ?? -1;

      expect(initialHeight, greaterThan(0));
      expect(expandedHeight, greaterThan(0));
      // Verify that the height increased when the row was expanded
      expect(expandedHeight, greaterThan(initialHeight));

      expect(tester.takeException(), isNull);
    });

    testWidgets('Expanded row content should be visible after tap', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SingleChildScrollView(child: ListSearchable(search: mockSearch)),
      );

      // Wait for initial load
      await tester.pumpAndSettle();

      // Verify initial state
      expect(find.text('Test Feed 1'), findsOneWidget);
      expect(find.text('Test Feed 2'), findsOneWidget);

      // Verify that expanded content is not visible initially
      expect(find.text('url'), findsNothing);
      expect(find.text('description'), findsNothing);

      // Find the first feed row and tap it to expand
      final firstFeedRow = find.byType(ds.TableRow).first;
      await tester.tap(firstFeedRow);
      await tester.pumpAndSettle();

      // Verify that the expanded content is now visible
      expect(find.text('url'), findsOneWidget);
      expect(find.text('description'), findsOneWidget);

      expect(tester.takeException(), isNull);
    });
  });
}
