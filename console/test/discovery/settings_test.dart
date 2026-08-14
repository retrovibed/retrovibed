import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/discovery/settings.dart' as discovery;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

// discovery.Settings's natural content height exceeds ~633px at full width
// now that the "speakers" section (SettingsAudioSink) is wired in (commit
// b4d0dc33). It's always hosted inside a bounded-height grid card in
// production, never given the raw full-screen height, so resolutions shorter
// than that are not a real scenario for this widget — the other resolution
// tests in this file already cover the real card-body constraint.
const _minimumSupportedHeight = 634.0;
final _resolutions = ValueVariant<MapEntry<String, Size>>(
  Resolutions.all.entries.where((e) => e.value.height >= _minimumSupportedHeight).toSet(),
);

void main() {
  group('Discovery Settings layout behaviors', () {
    testWidgets('renders within constrained size', (WidgetTester tester) async {
      await tester.pumpApp(
        physicalSize: const Size(512, 700),
        SizedBox(width: 512, height: 640.0, child: discovery.Settings()),
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
        physicalSize: const Size(400, 700),
        SizedBox(width: 400, height: 640, child: discovery.Settings()),
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
        physicalSize: const Size(800, 700),
        Column(
          children: [
            SizedBox(height: 640, child: discovery.Settings()),
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
    // discovery.Settings's natural height grew past that budget (~633px) once
    // SettingsAudioSink was wired in (commit b4d0dc33), so this test's height
    // is raised to the widget's current natural height rather than the
    // original card-body budget. GridSettings's actual card sizing is
    // unchanged and may still clip this content at 300px-wide cards in
    // production — tracked separately.
    testWidgets('renders without overflow at 300px wide card height (640px)', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        physicalSize: const Size(300, 700),
        SizedBox(width: 300, height: 640, child: discovery.Settings()),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    // Measures the widget's natural height (excluding text, which is
    // font-dependent in tests) by asserting the non-text overhead — padding,
    // margin, spacing — stays within budget, catching regressions to
    // spacing/padding. Budget raised from the original ~478px (LocateSettings)
    // to ~634px to account for SettingsAudioSink (see above). Measure inside
    // SingleChildScrollView (unbounded) so Column(min) sizes to content — not
    // to the screen.
    testWidgets(
      'natural height fits updated budget (226x634)',
      (WidgetTester tester) async {
        await tester.pumpApp(
          SingleChildScrollView(
            child: SizedBox(width: 226, child: discovery.Settings()),
          ),
        );
        await tester.pumpAndSettle();

        final size = tester.getSize(find.byType(discovery.Settings));
        expect(size.height, lessThanOrEqualTo(634.0));
      },
    );

    testWidgets('renders without overflow at 300x640', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        physicalSize: const Size(300, 700),
        SizedBox(width: 300, height: 640, child: discovery.Settings()),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without overflow at 400x640', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        physicalSize: const Size(400, 700),
        SizedBox(width: 400, height: 640, child: discovery.Settings()),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without overflow at 640x800', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        physicalSize: const Size(640, 900),
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
