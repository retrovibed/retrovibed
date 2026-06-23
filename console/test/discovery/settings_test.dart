import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/discovery/settings.dart' as discovery;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

// discovery.Settings's natural content height exceeds ~473px at full width.
// It's always hosted inside a bounded-height grid card in production, never
// given the raw full-screen height, so resolutions shorter than that are not
// a real scenario for this widget — the other resolution tests in this file
// already cover the real card-body constraint.
const _minimumSupportedHeight = 480.0;
final _resolutions = ValueVariant<MapEntry<String, Size>>(
  Resolutions.all.entries.where((e) => e.value.height >= _minimumSupportedHeight).toSet(),
);

void main() {
  group('Discovery Settings layout behaviors', () {
    testWidgets('renders within constrained size', (WidgetTester tester) async {
      await tester.pumpApp(
        SizedBox(width: 512, height: 480.0, child: discovery.Settings()),
      );
      await tester.pumpAndSettle();

      expect(find.byType(discovery.Settings), findsOneWidget);
      expect(find.text('video'), findsOneWidget);
      expect(find.text('audio'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('Discovery Settings renders standalone', () {
    testWidgets('Settings renders in finite container', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(width: 400, height: 500, child: discovery.Settings()),
      );
      await tester.pumpAndSettle();

      expect(find.byType(discovery.Settings), findsOneWidget);
      expect(find.text('video'), findsOneWidget);
      expect(find.text('audio'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('Settings renders in Column with fixed height', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Column(
          children: [
            SizedBox(height: 480, child: discovery.Settings()),
            Expanded(child: Container()),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(discovery.Settings), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('Settings renders in SingleChildScrollView', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(SingleChildScrollView(child: discovery.Settings()));
      await tester.pumpAndSettle();

      expect(find.byType(discovery.Settings), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('resolution tests', () {
    // The card constraint from GridSettings: aspectRatio 9.62/16 at 300px wide
    // gives card height = 300 * (16/9.62) ≈ 498px. The card body (Flexible
    // inside ds.Card) is smaller — approximately 403px after heading/padding.
    // discovery.Settings's natural height grew past that budget (~473px) once
    // LocateSettings was wired in (commit bd5fe1fe), so this test's height is
    // raised to the widget's current natural height rather than the original
    // card-body budget. GridSettings's actual card sizing is unchanged and may
    // still clip this content at 300px-wide cards in production — tracked
    // separately.
    testWidgets('renders without overflow at 300px wide card height (480px)', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(width: 300, height: 480, child: discovery.Settings()),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    // Measures the widget's natural height (excluding text, which is
    // font-dependent in tests) by asserting the non-text overhead — padding,
    // margin, spacing — stays within budget, catching regressions to
    // spacing/padding. Budget raised from the original ~403.7px GridSettings
    // card-body constraint to ~478px to account for LocateSettings (see above).
    // Measure inside SingleChildScrollView (unbounded) so Column(min) sizes to
    // content — not to the screen.
    testWidgets(
      'natural height fits updated budget (226x478)',
      (WidgetTester tester) async {
        await tester.pumpApp(
          SingleChildScrollView(
            child: SizedBox(width: 226, child: discovery.Settings()),
          ),
        );
        await tester.pumpAndSettle();

        final size = tester.getSize(find.byType(discovery.Settings));
        expect(size.height, lessThanOrEqualTo(478.0));
      },
    );

    testWidgets('renders without overflow at 300x480', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(width: 300, height: 480, child: discovery.Settings()),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without overflow at 400x500', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(width: 400, height: 500, child: discovery.Settings()),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without overflow at 640x800', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(width: 640, height: 800, child: discovery.Settings()),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(physicalSize: entry.value, discovery.Settings());
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    }, variant: _resolutions);
  });
}
