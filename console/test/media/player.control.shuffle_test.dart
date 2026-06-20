import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/media/player.control.shuffle.dart';
import 'package:retrovibed/media/play.queue.dart';
import 'package:retrovibed/media/media.row.display.dart';
import 'package:retrovibed/media/api.dart' as api;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

PlayableMedia _media(String id, String title) => PlayableMedia(api.Media(id: id, description: title));

Finder _strategyIcon(String tooltip) => find.byWidgetPredicate((w) => w is IconButton && w.tooltip == tooltip);

void main() {
  group('PlaybackStrategyModal', () {
    testWidgets('renders an icon button for each strategy option', (tester) async {
      await tester.pumpApp(
        PlaybackStrategyModal(
          current: search,
          history: const [],
          onModeSelected: (_) {},
          onHistorySelected: (_) {},
        ),
      );
      await tester.pumpAndSettle();

      expect(_strategyIcon('Search'), findsOneWidget);
      expect(_strategyIcon('Random'), findsOneWidget);
      expect(_strategyIcon('Auto'), findsOneWidget);
    });

    testWidgets('tints the icon matching current with the theme primary color', (tester) async {
      await tester.pumpApp(
        PlaybackStrategyModal(
          current: random,
          history: const [],
          onModeSelected: (_) {},
          onHistorySelected: (_) {},
        ),
      );
      await tester.pumpAndSettle();

      final theme = Theme.of(tester.element(_strategyIcon('Random')));
      expect(tester.widget<IconButton>(_strategyIcon('Random')).color, theme.colorScheme.primary);
      expect(tester.widget<IconButton>(_strategyIcon('Search')).color, isNull);
    });

    testWidgets('invokes onModeSelected with the tapped mode', (tester) async {
      RangeFn? selected;
      await tester.pumpApp(
        PlaybackStrategyModal(
          current: search,
          history: const [],
          onModeSelected: (mode) => selected = mode,
          onHistorySelected: (_) {},
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(_strategyIcon('Auto'));
      expect(selected, acoustic);
    });

    testWidgets('renders without history rows when history is empty', (tester) async {
      await tester.pumpApp(
        PlaybackStrategyModal(
          current: search,
          history: const [],
          onModeSelected: (_) {},
          onHistorySelected: (_) {},
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(RowDisplay), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders one row per history entry, most-recently-played first', (tester) async {
      final history = [_media('a', 'Song A'), _media('b', 'Song B'), _media('c', 'Song C')];

      await tester.pumpApp(
        PlaybackStrategyModal(
          current: search,
          history: history,
          onModeSelected: (_) {},
          onHistorySelected: (_) {},
        ),
      );
      await tester.pumpAndSettle();

      final rows = tester.widgetList<RowDisplay>(find.byType(RowDisplay)).toList();
      expect(rows.map((r) => r.media.description), ['Song C', 'Song B', 'Song A']);
    });

    testWidgets('invokes onHistorySelected with the tapped track', (tester) async {
      final history = [_media('a', 'Song A'), _media('b', 'Song B')];
      PlayableMedia? selected;

      await tester.pumpApp(
        PlaybackStrategyModal(
          current: search,
          history: history,
          onModeSelected: (_) {},
          onHistorySelected: (item) => selected = item,
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('Song A'));
      expect(selected?.current.id, 'a');
    });
  });
}
