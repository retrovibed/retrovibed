import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/timex.dart' as timex;
import 'package:retrovibed/design.kit/inputs/time.range.dart';

final _now = DateTime(2024, 6, 15);
final _segments = [
  timex.Range(DateTime(_now.year, _now.month - 1, _now.day), _now),
  timex.Range(DateTime(_now.year, _now.month - 3, _now.day), _now),
  timex.Range(DateTime(_now.year - 1, _now.month, _now.day), _now),
  timex.Range(DateTime(_now.year - 3, _now.month, _now.day), _now),
];

Widget buildTestWidget({
  required timex.Range selected,
  required ValueChanged<timex.Range> onChanged,
}) {
  return MaterialApp(
    home: Scaffold(
      body: TimeRange(
        segments: _segments,
        selected: selected,
        onChanged: onChanged,
      ),
    ),
  );
}

void main() {
  group('TimeRange input', () {
    group('screen resolutions', () {
      testWidgets('renders dropdown at minimum width (300x568)', (
        WidgetTester tester,
      ) async {
        tester.view.physicalSize = Size(300, 568);
        tester.view.devicePixelRatio = 1.0;
        addTearDown(() => tester.view.reset());

        await tester.pumpWidget(buildTestWidget(
          selected: _segments.first,
          onChanged: (_) {},
        ));

        expect(find.byType(DropdownButton<timex.Range>), findsOneWidget);
        expect(find.byType(SegmentedButton<timex.Range>), findsNothing);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders dropdown on small mobile (320x568)', (
        WidgetTester tester,
      ) async {
        tester.view.physicalSize = Size(320, 568);
        tester.view.devicePixelRatio = 1.0;
        addTearDown(() => tester.view.reset());

        await tester.pumpWidget(buildTestWidget(
          selected: _segments.first,
          onChanged: (_) {},
        ));

        expect(find.byType(DropdownButton<timex.Range>), findsOneWidget);
        expect(find.byType(SegmentedButton<timex.Range>), findsNothing);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders segmented button on iPhone SE (375x667)', (
        WidgetTester tester,
      ) async {
        tester.view.physicalSize = Size(375, 667);
        tester.view.devicePixelRatio = 1.0;
        addTearDown(() => tester.view.reset());

        await tester.pumpWidget(buildTestWidget(
          selected: _segments.first,
          onChanged: (_) {},
        ));

        expect(find.byType(SegmentedButton<timex.Range>), findsOneWidget);
        expect(find.byType(DropdownButton<timex.Range>), findsNothing);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders segmented button on tablet (768x1024)', (
        WidgetTester tester,
      ) async {
        tester.view.physicalSize = Size(768, 1024);
        tester.view.devicePixelRatio = 1.0;
        addTearDown(() => tester.view.reset());

        await tester.pumpWidget(buildTestWidget(
          selected: _segments.first,
          onChanged: (_) {},
        ));

        expect(find.byType(SegmentedButton<timex.Range>), findsOneWidget);
        expect(find.byType(DropdownButton<timex.Range>), findsNothing);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders segmented button on desktop (1920x1080)', (
        WidgetTester tester,
      ) async {
        tester.view.physicalSize = Size(1920, 1080);
        tester.view.devicePixelRatio = 1.0;
        addTearDown(() => tester.view.reset());

        await tester.pumpWidget(buildTestWidget(
          selected: _segments.first,
          onChanged: (_) {},
        ));

        expect(find.byType(SegmentedButton<timex.Range>), findsOneWidget);
        expect(find.byType(DropdownButton<timex.Range>), findsNothing);
        expect(tester.takeException(), isNull);
      });
    });

    testWidgets('displays all segment labels', (WidgetTester tester) async {
      await tester.pumpWidget(buildTestWidget(
        selected: _segments.first,
        onChanged: (_) {},
      ));

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

      await tester.pumpWidget(buildTestWidget(
        selected: _segments.first,
        onChanged: (value) => changed = value,
      ));

      await tester.tap(find.text('1 Year'));
      await tester.pumpAndSettle();

      expect(changed, equals(_segments[2]));
      expect(tester.takeException(), isNull);
    });

    testWidgets('calls onChanged when a dropdown item is selected', (
      WidgetTester tester,
    ) async {
      tester.view.physicalSize = Size(300, 568);
      tester.view.devicePixelRatio = 1.0;
      addTearDown(() => tester.view.reset());

      timex.Range? changed;

      await tester.pumpWidget(buildTestWidget(
        selected: _segments.first,
        onChanged: (value) => changed = value,
      ));

      // Tap the dropdown to open it.
      await tester.tap(find.byType(DropdownButton<timex.Range>));
      await tester.pumpAndSettle();

      // Select '1 Year' from the dropdown menu.
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
      tester.view.physicalSize = Size(300, 568);
      tester.view.devicePixelRatio = 1.0;
      addTearDown(() => tester.view.reset());

      final staleSelection = timex.Range(
        DateTime(2020, 1, 1),
        DateTime(2020, 12, 31),
      );

      await tester.pumpWidget(MaterialApp(
        home: Scaffold(
          body: TimeRange(
            segments: _segments,
            selected: staleSelection,
            onChanged: (_) {},
          ),
        ),
      ));

      expect(find.byType(DropdownButton<timex.Range>), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
