import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/inputs/rate.limit.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('RateLimit input', () {
    final presets = [
      (label: 'low', value: 10, unit: 'sec'),
      (label: 'medium', value: 100, unit: 'sec'),
      (label: 'high', value: 1000, unit: 'sec'),
    ];

    group('screen resolutions', () {
      testWidgets('renders at minimum width (300x568)', (
        WidgetTester tester,
      ) async {
        tester.view.physicalSize = Size(300, 568);
        tester.view.devicePixelRatio = 1.0;
        addTearDown(() => tester.view.reset());

        await tester.pumpApp(
          RateLimit(
            value: 0,
            onChanged: (_) {},
            presets: presets,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(TextFormField), findsOneWidget);
        expect(find.text('sec'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders on small mobile (320x568)', (
        WidgetTester tester,
      ) async {
        tester.view.physicalSize = Size(320, 568);
        tester.view.devicePixelRatio = 1.0;
        addTearDown(() => tester.view.reset());

        await tester.pumpApp(
          RateLimit(
            value: 0,
            onChanged: (_) {},
            presets: presets,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(TextFormField), findsOneWidget);
        expect(find.byType(IconButton), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders on iPhone SE (375x667)', (
        WidgetTester tester,
      ) async {
        tester.view.physicalSize = Size(375, 667);
        tester.view.devicePixelRatio = 1.0;
        addTearDown(() => tester.view.reset());

        await tester.pumpApp(
          RateLimit(
            value: 0,
            onChanged: (_) {},
            presets: presets,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(TextFormField), findsOneWidget);
        expect(find.text('0 means no limit'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders on tablet (768x1024)', (
        WidgetTester tester,
      ) async {
        tester.view.physicalSize = Size(768, 1024);
        tester.view.devicePixelRatio = 1.0;
        addTearDown(() => tester.view.reset());

        await tester.pumpApp(
          RateLimit(
            value: 0,
            onChanged: (_) {},
            presets: presets,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(TextFormField), findsOneWidget);
        expect(find.byType(IconButton), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders on desktop (1920x1080)', (
        WidgetTester tester,
      ) async {
        tester.view.physicalSize = Size(1920, 1080);
        tester.view.devicePixelRatio = 1.0;
        addTearDown(() => tester.view.reset());

        await tester.pumpApp(
          RateLimit(
            value: 0,
            onChanged: (_) {},
            presets: presets,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(TextFormField), findsOneWidget);
        expect(find.byType(IconButton), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    testWidgets('displays current value in text field', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        RateLimit(
          value: 100,
          onChanged: (_) {},
          presets: presets,
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('100'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('displays helper text', (WidgetTester tester) async {
      await tester.pumpApp(
        RateLimit(
          value: 0,
          onChanged: (_) {},
          presets: presets,
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('0 means no limit'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('displays unit label', (WidgetTester tester) async {
      await tester.pumpApp(
        RateLimit(
          value: 0,
          onChanged: (_) {},
          presets: presets,
          units: const ['sec'],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('sec'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('calls onChanged when text changes', (
      WidgetTester tester,
    ) async {
      int? changed;

      await tester.pumpApp(
        RateLimit(
          value: 0,
          onChanged: (value) => changed = value,
          presets: presets,
        ),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextFormField), '50');
      await tester.pumpAndSettle();

      expect(changed, equals(50));
      expect(tester.takeException(), isNull);
    });

    testWidgets('dropdown button toggles expanded state', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        RateLimit(
          value: 0,
          onChanged: (_) {},
          presets: presets,
        ),
      );
      await tester.pumpAndSettle();

      // Initially expanded should be false - presets not visible
      expect(find.byKey(const ValueKey('presets')), findsNothing);

      // Tap dropdown to expand
      await tester.tap(find.byType(IconButton));
      await tester.pumpAndSettle();

      // Now expanded should be true - presets visible
      expect(find.byKey(const ValueKey('presets')), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('displays unit selection chips when multiple units', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        RateLimit(
          value: 0,
          onChanged: (_) {},
          presets: presets,
          units: const ['sec', 'min', 'hour'],
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byType(IconButton));
      await tester.pumpAndSettle();

      expect(find.byType(Chip), findsNothing);
      expect(find.byType(OutlinedButton).at(0), findsOneWidget);
      expect(find.byType(OutlinedButton).at(1), findsOneWidget);
      expect(find.byType(OutlinedButton).at(2), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('displays preset buttons', (WidgetTester tester) async {
      await tester.pumpApp(
        RateLimit(
          value: 0,
          onChanged: (_) {},
          presets: presets,
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byType(IconButton));
      await tester.pumpAndSettle();

      expect(find.text('low'), findsOneWidget);
      expect(find.text('medium'), findsOneWidget);
      expect(find.text('high'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('preset selection calls onChanged', (
      WidgetTester tester,
    ) async {
      int? changed;

      await tester.pumpApp(
        RateLimit(
          value: 0,
          onChanged: (value) => changed = value,
          presets: presets,
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byType(IconButton));
      await tester.pumpAndSettle();

      await tester.tap(find.text('low'));
      await tester.pumpAndSettle();

      expect(changed, equals(10));
      expect(tester.takeException(), isNull);
    });

    testWidgets('preset selection updates value display', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        RateLimit(
          value: 0,
          onChanged: (_) {},
          presets: presets,
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byType(IconButton));
      await tester.pumpAndSettle();

      await tester.tap(find.text('medium'));
      await tester.pumpAndSettle();

      expect(find.text('100'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('preset selection sets unit', (WidgetTester tester) async {
      await tester.pumpApp(
        RateLimit(
          value: 0,
          onChanged: (_) {},
          presets: [
            (label: 'low', value: 10, unit: 'sec'),
            (label: 'medium', value: 100, unit: 'min'),
            (label: 'high', value: 1000, unit: 'hour'),
          ],
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byType(IconButton));
      await tester.pumpAndSettle();

      await tester.tap(find.text('low'));
      await tester.pumpAndSettle();

      expect(find.text('sec'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders empty presets list gracefully', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        RateLimit(
          value: 0,
          onChanged: (_) {},
          presets: [],
          units: const ['sec'],
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byType(IconButton));
      await tester.pumpAndSettle();

      expect(find.byKey(const ValueKey('presets')), findsNothing);
      expect(find.text('sec'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('validates numeric input only', (
      WidgetTester tester,
    ) async {
      int? changed;

      await tester.pumpApp(
        RateLimit(
          value: 0,
          onChanged: (value) => changed = value,
          presets: presets,
        ),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextFormField), 'abc');
      await tester.pumpAndSettle();

      expect(changed, isNull);
      expect(tester.takeException(), isNull);
    });

    testWidgets('retains focus when typing after preset selection', (
      WidgetTester tester,
    ) async {
      // Regression: key: ValueKey((_preset, _unit)) changed when _onTextChanged
      // cleared _preset, destroying the TextFormField and losing focus.
      await tester.pumpApp(
        RateLimit(
          value: 0,
          onChanged: (_) {},
          presets: presets,
        ),
      );
      await tester.pumpAndSettle();

      // Select a preset so _preset becomes non-null.
      await tester.tap(find.byType(IconButton));
      await tester.pumpAndSettle();
      await tester.tap(find.text('low'));
      await tester.pumpAndSettle();

      // Focus the text field.
      await tester.tap(find.byType(TextFormField));
      await tester.pump();

      final focusBefore = tester.binding.focusManager.primaryFocus;
      expect(focusBefore, isNotNull);

      // Type without re-tapping — internally clears _preset, used to change the key.
      tester.testTextInput.enterText('5');
      await tester.pump();

      expect(
        tester.binding.focusManager.primaryFocus,
        same(focusBefore),
        reason: 'text field should retain focus when typing clears a previously selected preset',
      );
      expect(tester.takeException(), isNull);
    });
  });
}
