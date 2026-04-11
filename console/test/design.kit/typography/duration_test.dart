import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('elapsed formatter', () {
    testWidgets('formats standard positive duration as HH:MM:SS', (
      WidgetTester tester,
    ) async {
      final duration = const Duration(hours: 2, minutes: 30, seconds: 45);

      await tester.pumpApp(
        Scaffold(body: ds.Duration(duration, formatter: ds.Duration.elapsed)),
      );
      await tester.pumpAndSettle();

      expect(find.text('02:30:45'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('formats negative duration with minus sign', (
      WidgetTester tester,
    ) async {
      final duration = const Duration(hours: -1, minutes: -15);

      await tester.pumpApp(
        Scaffold(body: ds.Duration(duration, formatter: ds.Duration.elapsed)),
      );
      await tester.pumpAndSettle();

      expect(find.text('-01:15:00'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('formats zero duration as 00:00:00', (
      WidgetTester tester,
    ) async {
      final duration = Duration.zero;

      await tester.pumpApp(
        Scaffold(body: ds.Duration(duration, formatter: ds.Duration.elapsed)),
      );
      await tester.pumpAndSettle();

      expect(find.text('00:00:00'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('formats hours exceeding 24 without wrapping', (
      WidgetTester tester,
    ) async {
      final duration = const Duration(hours: 50);

      await tester.pumpApp(
        Scaffold(body: ds.Duration(duration, formatter: ds.Duration.elapsed)),
      );
      await tester.pumpAndSettle();

      expect(find.text('50:00:00'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('pads single digit minutes and seconds', (
      WidgetTester tester,
    ) async {
      final duration = const Duration(hours: 1, minutes: 5, seconds: 9);

      await tester.pumpApp(
        Scaffold(body: ds.Duration(duration, formatter: ds.Duration.elapsed)),
      );
      await tester.pumpAndSettle();

      expect(find.text('01:05:09'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('ago formatter', () {
    testWidgets('formats days as Xd ago', (WidgetTester tester) async {
      final duration = const Duration(days: 3);

      await tester.pumpApp(
        Scaffold(body: ds.Duration(duration, formatter: ds.Duration.ago)),
      );
      await tester.pumpAndSettle();

      expect(find.text('3d ago'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('formats hours as Xh ago', (WidgetTester tester) async {
      final duration = const Duration(hours: 5);

      await tester.pumpApp(
        Scaffold(body: ds.Duration(duration, formatter: ds.Duration.ago)),
      );
      await tester.pumpAndSettle();

      expect(find.text('5h ago'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('formats minutes as Xm ago', (WidgetTester tester) async {
      final duration = const Duration(minutes: 45);

      await tester.pumpApp(
        Scaffold(body: ds.Duration(duration, formatter: ds.Duration.ago)),
      );
      await tester.pumpAndSettle();

      expect(find.text('45m ago'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('formats less than a minute as just now', (
      WidgetTester tester,
    ) async {
      final duration = const Duration(seconds: 30);

      await tester.pumpApp(
        Scaffold(body: ds.Duration(duration, formatter: ds.Duration.ago)),
      );
      await tester.pumpAndSettle();

      expect(find.text('just now'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('uses absolute value for negative duration', (
      WidgetTester tester,
    ) async {
      final duration = const Duration(hours: -5);

      await tester.pumpApp(
        Scaffold(body: ds.Duration(duration, formatter: ds.Duration.ago)),
      );
      await tester.pumpAndSettle();

      expect(find.text('5h ago'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('widget behavior', () {
    testWidgets('uses relative as default formatter', (
      WidgetTester tester,
    ) async {
      final duration = const Duration(hours: 2);

      await tester.pumpApp(Scaffold(body: ds.Duration(duration)));
      await tester.pumpAndSettle();

      expect(find.text('in 2h'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('accepts custom formatter', (WidgetTester tester) async {
      final duration = const Duration(minutes: 90);
      Widget customFormatter(Duration d) => Text('custom: ${d.inMinutes}m');

      await tester.pumpApp(
        Scaffold(body: ds.Duration(duration, formatter: customFormatter)),
      );
      await tester.pumpAndSettle();

      expect(find.text('custom: 90m'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('layout - finite constraints', () {
    testWidgets('renders in Container with fixed dimensions', (
      WidgetTester tester,
    ) async {
      final duration = const Duration(hours: 1);

      await tester.pumpApp(
        Scaffold(
          body: Container(
            width: 200,
            height: 100,
            child: ds.Duration(duration),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('in 1h'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in SizedBox with tight constraints', (
      WidgetTester tester,
    ) async {
      final duration = const Duration(minutes: 30);

      await tester.pumpApp(
        Scaffold(
          body: SizedBox(
            width: 100,
            height: 50,
            child: ds.Duration(duration),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('in 30m'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders with Expanded in Row', (WidgetTester tester) async {
      final duration = const Duration(hours: 2);

      await tester.pumpApp(
        Scaffold(
          body: Row(
            children: [
              Expanded(child: ds.Duration(duration)),
              const SizedBox(width: 50),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('in 2h'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders with Expanded in Column', (WidgetTester tester) async {
      final duration = const Duration(days: 1);

      await tester.pumpApp(
        Scaffold(
          body: Column(
            children: [
              Expanded(child: ds.Duration(duration)),
              const SizedBox(height: 50),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('in 1d'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('layout - flex containers', () {
    testWidgets('renders in Row without overflow', (WidgetTester tester) async {
      final duration = const Duration(hours: 1);

      await tester.pumpApp(
        Scaffold(
          body: Row(
            children: [ds.Duration(duration), const Text('other content')],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('in 1h'), findsOneWidget);
      expect(find.text('other content'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in Column without overflow', (
      WidgetTester tester,
    ) async {
      final duration = const Duration(minutes: 15);

      await tester.pumpApp(
        Scaffold(
          body: Column(
            children: [ds.Duration(duration), const Text('other content')],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('in 15m'), findsOneWidget);
      expect(find.text('other content'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in ListView', (WidgetTester tester) async {
      await tester.pumpApp(
        Scaffold(
          body: ListView(
            children: [
              ds.Duration(const Duration(hours: 1)),
              ds.Duration(const Duration(days: 2)),
              ds.Duration(const Duration(minutes: 30)),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('in 1h'), findsOneWidget);
      expect(find.text('in 2d'), findsOneWidget);
      expect(find.text('in 30m'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in horizontal ScrollView', (
      WidgetTester tester,
    ) async {
      final duration = const Duration(hours: 3);

      await tester.pumpApp(
        Scaffold(
          body: SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            child: Row(
              children: [
                ds.Duration(duration),
                const Text('scrollable content'),
              ],
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('in 3h'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('layout - intrinsic dimensions', () {
    testWidgets('provides valid intrinsic width', (WidgetTester tester) async {
      final duration = const Duration(hours: 1);

      await tester.pumpApp(
        Scaffold(body: IntrinsicWidth(child: ds.Duration(duration))),
      );
      await tester.pumpAndSettle();

      expect(find.text('in 1h'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('provides valid intrinsic height', (WidgetTester tester) async {
      final duration = const Duration(hours: 1);

      await tester.pumpApp(
        Scaffold(body: IntrinsicHeight(child: ds.Duration(duration))),
      );
      await tester.pumpAndSettle();

      expect(find.text('in 1h'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
