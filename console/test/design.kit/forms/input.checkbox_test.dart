import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/forms/input.checkbox.dart' as c;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('Checkbox Widget Tests', () {
    group('layout', () {
      testWidgets('renders without overflow in unconstrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Material(
            child: Center(
              child: c.Checkbox(const Text('Test Checkbox'), value: false),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('Test Checkbox'), findsOneWidget);
        expect(tester.takeException(), isNull);
        final renderBox = tester.renderObject<RenderBox>(
          find.byType(c.Checkbox),
        );

        expect(renderBox.size.width, 800);
        expect(renderBox.size.height, 48);
      });

      testWidgets(
        'renders without overflow in constrained environment',
        (WidgetTester tester) async {
          await tester.pumpApp(
            MediaQuery(
              data: const MediaQueryData(textScaler: TextScaler.linear(1.0)),
              child: Material(
                child: Center(
                  child: ConstrainedBox(
                    constraints: const BoxConstraints(
                      maxWidth: 200,
                      maxHeight: 100,
                    ),
                    child: c.Checkbox(
                      const Text('Test Checkbox'),
                      value: false,
                    ),
                  ),
                ),
              ),
            ),
            theme: ThemeData(
              fontFamily: 'Ahem',
            ), // ensure no differences due to platform specific font behaviors.
          );
          await tester.pumpAndSettle();

          expect(find.text('Test Checkbox'), findsOneWidget);
          expect(tester.takeException(), isNull);

          final renderBox = tester.renderObject<RenderBox>(
            find.byType(c.Checkbox),
          );

          expect(renderBox.size.width, 200);
          expect(renderBox.size.height, 64);
        },
        variant: TargetPlatformVariant.all(),
      );
    });

    testWidgets('renders multiple checkboxes in Wrap without overflow', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Material(
          child: Wrap(
            children: [
              c.Checkbox(const Text('Checkbox 1'), value: false),
              c.Checkbox(const Text('Checkbox 2'), value: false),
              c.Checkbox(const Text('Checkbox 3'), value: false),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Checkbox 1'), findsOneWidget);
      expect(find.text('Checkbox 2'), findsOneWidget);
      expect(find.text('Checkbox 3'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets(
      'renders multiple checkboxes in constrained Wrap without overflow',
      (WidgetTester tester) async {
        await tester.pumpApp(
          Material(
            child: ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 200, maxHeight: 100),
              child: Wrap(
                children: [
                  c.Checkbox(const Text('Checkbox 1'), value: false),
                  c.Checkbox(const Text('Checkbox 2'), value: false),
                  c.Checkbox(const Text('Checkbox 3'), value: false),
                ],
              ),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('Checkbox 1'), findsOneWidget);
        expect(find.text('Checkbox 2'), findsOneWidget);
        expect(find.text('Checkbox 3'), findsOneWidget);
        expect(tester.takeException(), isNull);
      },
    );

    testWidgets('renders multiple checkboxes in Column without overflow', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Material(
          child: Column(
            children: [
              c.Checkbox(const Text('Checkbox 1'), value: false),
              c.Checkbox(const Text('Checkbox 2'), value: false),
              c.Checkbox(const Text('Checkbox 3'), value: false),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Checkbox 1'), findsOneWidget);
      expect(find.text('Checkbox 2'), findsOneWidget);
      expect(find.text('Checkbox 3'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets(
      'renders multiple checkboxes in constrained Column without overflow',
      (WidgetTester tester) async {
        await tester.pumpApp(
          Material(
            child: ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 200, maxHeight: 150),
              child: Column(
                children: [
                  c.Checkbox(const Text('Checkbox 1'), value: false),
                  c.Checkbox(const Text('Checkbox 2'), value: false),
                  c.Checkbox(const Text('Checkbox 3'), value: false),
                ],
              ),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('Checkbox 1'), findsOneWidget);
        expect(find.text('Checkbox 2'), findsOneWidget);
        expect(find.text('Checkbox 3'), findsOneWidget);
        expect(tester.takeException(), isNull);
      },
    );
  });
}
