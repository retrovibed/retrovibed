import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/discovery/settings.dart' as discovery;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

void main() {
  group('Discovery Settings layout behaviors', () {
    testWidgets('renders within constrained size', (WidgetTester tester) async {
      await tester.pumpApp(
        SizedBox(width: 512, height: 454.0, child: discovery.Settings()),
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
            SizedBox(height: 452, child: discovery.Settings()),
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
    // Settings must fit without overflow at that height.
    testWidgets('renders without overflow at 300px wide card height (403px)', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(width: 300, height: 403, child: discovery.Settings()),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    // The GridSettings card body at 300px screen width constrains Settings to
    // ~403.7px tall via Flexible(tight). We measure the widget's natural height
    // (excluding text, which is font-dependent in tests) by asserting the
    // non-text overhead — padding, margin, spacing — stays within the budget.
    // This catches regressions to spacing/padding that cause overflow at runtime.
    // Measure natural height inside SingleChildScrollView (unbounded) so
    // Column(min) sizes to content — not to the screen. The result must fit
    // within the ~403.7px Flexible(tight) card body in GridSettings.
    testWidgets(
      'natural height fits GridSettings card body constraint (226x403.7)',
      (WidgetTester tester) async {
        await tester.pumpApp(
          SingleChildScrollView(
            child: SizedBox(width: 226, child: discovery.Settings()),
          ),
        );
        await tester.pumpAndSettle();

        final size = tester.getSize(find.byType(discovery.Settings));
        expect(size.height, lessThanOrEqualTo(403.7));
      },
    );

    testWidgets('renders without overflow at 300x402', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(width: 300, height: 402, child: discovery.Settings()),
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
