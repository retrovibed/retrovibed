import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/timex.dart' as timex;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  // A date that produces a long formatted string with the default format.
  final verbose = DateTime.utc(2024, 1, 15, 14, 30);

  group('overflow - tight width constraints', () {
    testWidgets('ellipsis when constrained narrower than text', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: SizedBox(
            width: 50,
            height: 50,
            child: ds.Timestamp(verbose),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('no overflow in Row with competing siblings', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: Row(
            children: [
              const SizedBox(width: 300),
              Flexible(child: ds.Timestamp(verbose)),
              const SizedBox(width: 300),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('no overflow with leading and trailing in narrow container', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: SizedBox(
            width: 80,
            child: ds.Timestamp(
              verbose,
              leading: const Icon(Icons.calendar_today),
              trailing: const Icon(Icons.chevron_right),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });
  });

  group('overflow - flex containers', () {
    testWidgets('no overflow when wrapped in Flexible inside Row', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: Row(
            children: [
              Flexible(child: ds.Timestamp(verbose)),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('no overflow when wrapped in Expanded inside Row', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: Row(
            children: [
              Expanded(child: ds.Timestamp(verbose)),
              const SizedBox(width: 50),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('no overflow in narrow ListView', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: SizedBox(
            width: 100,
            child: ListView(
              children: [
                ds.Timestamp(verbose),
              ],
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });
  });

  group('adaptive default format', () {
    testWidgets('uses short numeric format below 120px', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: SizedBox(
            width: 100,
            height: 50,
            child: ds.Timestamp(verbose),
          ),
        ),
      );
      await tester.pumpAndSettle();

      // M/d/y format
      expect(find.textContaining('1/15/2024'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('uses abbreviated month format between 120px and 200px', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: SizedBox(
            width: 180,
            height: 50,
            child: ds.Timestamp(verbose),
          ),
        ),
      );
      await tester.pumpAndSettle();

      // MMM d, y format
      expect(find.textContaining('Jan 15, 2024'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('uses abbreviated month with time between 200px and 300px', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: SizedBox(
            width: 250,
            height: 50,
            child: ds.Timestamp(verbose),
          ),
        ),
      );
      await tester.pumpAndSettle();

      // MMM d, y hh:mm a format
      expect(find.textContaining('Jan 15, 2024'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('uses full verbose format at 300px and above', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: SizedBox(
            width: 400,
            height: 50,
            child: ds.Timestamp(verbose),
          ),
        ),
      );
      await tester.pumpAndSettle();

      // y MMMM EEEE d hh:mm a format
      expect(find.textContaining('2024 January'), findsOneWidget);
      expect(find.textContaining('Monday'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('custom format overrides adaptive default', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: SizedBox(
            width: 400,
            height: 50,
            child: ds.Timestamp(
              verbose,
              format: (dt) => '${dt.month}/${dt.day}',
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      // Custom format used even though there's plenty of space
      expect(find.text('1/15'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('overflow - special values in tight constraints', () {
    testWidgets('inf renders "never" without overflow in narrow container', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: SizedBox(
            width: 80,
            height: 50,
            child: ds.Timestamp(timex.inf),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('never'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('neginf renders "always" without overflow in narrow container', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: SizedBox(
            width: 80,
            height: 50,
            child: ds.Timestamp(timex.neginf),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('always'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('overflow - scrollable containers', () {
    testWidgets('no overflow in horizontal ScrollView', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            child: ds.Timestamp(verbose),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });
  });

  group('default rendering', () {
    testWidgets('neginf default rendering', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: SizedBox(
            width: 80,
            height: 50,
            child: ds.Timestamp(timex.neginf),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('always'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('inf default rendering', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: SizedBox(
            width: 80,
            height: 50,
            child: ds.Timestamp(timex.inf),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('never'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('empty string default rendering', (
      WidgetTester tester,
    ) async {
      // empty strings should render identically to neginf.
      await tester.pumpApp(
        Scaffold(
          body: SizedBox(
            width: 80,
            height: 50,
            child: ds.Timestamp.iso8601(""),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('always'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
