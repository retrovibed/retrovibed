import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/timex.dart' as timex;
import 'package:retrovibed/design.kit/inputs/time.range.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _now = DateTime(2024, 6, 15);
final _segments = [
  timex.Range(DateTime(_now.year, _now.month - 1, _now.day), _now),
  timex.Range(DateTime(_now.year, _now.month - 3, _now.day), _now),
  timex.Range(DateTime(_now.year - 1, _now.month, _now.day), _now),
  timex.Range(DateTime(_now.year - 3, _now.month, _now.day), _now),
];

void main() {
  group('TimeRange input', () {
    group('screen resolutions', () {
      testWidgets('renders dropdown at minimum width (300x568)', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          TimeRange(segments: _segments, selected: _segments.first, onChanged: (_) {}),
          physicalSize: Size(300, 568),
        );
        await tester.pumpAndSettle();

        expect(find.byType(DropdownButton<timex.Range>), findsOneWidget);
        expect(find.byType(SegmentedButton<timex.Range>), findsNothing);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders dropdown on small mobile (320x568)', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          TimeRange(segments: _segments, selected: _segments.first, onChanged: (_) {}),
          physicalSize: Size(320, 568),
        );
        await tester.pumpAndSettle();

        expect(find.byType(DropdownButton<timex.Range>), findsOneWidget);
        expect(find.byType(SegmentedButton<timex.Range>), findsNothing);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders segmented button on iPhone SE (375x667)', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          TimeRange(segments: _segments, selected: _segments.first, onChanged: (_) {}),
          physicalSize: Size(375, 667),
        );
        await tester.pumpAndSettle();

        expect(find.byType(SegmentedButton<timex.Range>), findsOneWidget);
        expect(find.byType(DropdownButton<timex.Range>), findsNothing);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders segmented button on tablet (768x1024)', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          TimeRange(segments: _segments, selected: _segments.first, onChanged: (_) {}),
          physicalSize: Size(768, 1024),
        );
        await tester.pumpAndSettle();

        expect(find.byType(SegmentedButton<timex.Range>), findsOneWidget);
        expect(find.byType(DropdownButton<timex.Range>), findsNothing);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders segmented button on desktop (1920x1080)', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          TimeRange(segments: _segments, selected: _segments.first, onChanged: (_) {}),
          physicalSize: Size(1920, 1080),
        );
        await tester.pumpAndSettle();

        expect(find.byType(SegmentedButton<timex.Range>), findsOneWidget);
        expect(find.byType(DropdownButton<timex.Range>), findsNothing);
        expect(tester.takeException(), isNull);
      });
    });

    testWidgets('displays all segment labels', (WidgetTester tester) async {
      await tester.pumpApp(
        TimeRange(segments: _segments, selected: _segments.first, onChanged: (_) {}),
      );
      await tester.pumpAndSettle();

      expect(find.text('1 Month'), findsOneWidget);
      expect(find.text('3 Months'), findsOneWidget);
      expect(find.text('1 Year'), findsOneWidget);
      expect(find.text('3 Years'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('calls onChanged when a segment is tapped', (
      WidgetTester tester,
    ) async {
      timex.Range? changed;

      await tester.pumpApp(
        TimeRange(
          segments: _segments,
          selected: _segments.first,
          onChanged: (value) => changed = value,
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('1 Year'));
      await tester.pumpAndSettle();

      expect(changed, equals(_segments[2]));
      expect(tester.takeException(), isNull);
    });

    testWidgets('calls onChanged when a dropdown item is selected', (
      WidgetTester tester,
    ) async {
      timex.Range? changed;

      await tester.pumpApp(
        TimeRange(
          segments: _segments,
          selected: _segments.first,
          onChanged: (value) => changed = value,
        ),
        physicalSize: Size(300, 568),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byType(DropdownButton<timex.Range>));
      await tester.pumpAndSettle();

      // DropdownButton renders duplicate items (one in button, one in overlay),
      // so tap the last match which is in the overlay.
      await tester.tap(find.text('1 Year').last);
      await tester.pumpAndSettle();

      expect(changed, equals(_segments[2]));
      expect(tester.takeException(), isNull);
    });

    testWidgets('dropdown falls back to first segment when selected is not in segments', (
      WidgetTester tester,
    ) async {
      final staleSelection = timex.Range(
        DateTime(2020, 1, 1),
        DateTime(2020, 12, 31),
      );

      await tester.pumpApp(
        TimeRange(
          segments: _segments,
          selected: staleSelection,
          onChanged: (_) {},
        ),
        physicalSize: Size(300, 568),
      );
      await tester.pumpAndSettle();

      expect(find.byType(DropdownButton<timex.Range>), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
