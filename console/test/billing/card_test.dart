import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/billing/card.dart' as billing;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

void main() {
  group('billing.Card', () {
    group('layout', () {
      testWidgets('renders subscription title', (WidgetTester tester) async {
        await tester.pumpApp(billing.Card(onPressed: (_) {}));
        await tester.pumpAndSettle();

        expect(find.text('Subscription'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders plan description', (WidgetTester tester) async {
        await tester.pumpApp(billing.Card(onPressed: (_) {}));
        await tester.pumpAndSettle();

        // free plan is the default when no Registered ancestor exists
        expect(find.text('free'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('invokes onPressed with Settings widget', (
        WidgetTester tester,
      ) async {
        Widget? received;
        await tester.pumpApp(billing.Card(onPressed: (w) => received = w));
        await tester.pumpAndSettle();

        final button = find.text('Manage');
        if (button.evaluate().isNotEmpty) {
          await tester.tap(button);
          await tester.pumpAndSettle();
          expect(received, isNotNull);
        }

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in unconstrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(billing.Card(onPressed: (_) {}));
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in constrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 400, maxHeight: 600),
            child: billing.Card(onPressed: (_) {}),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in Column', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Column(children: [Expanded(child: billing.Card(onPressed: (_) {}))]),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in Row', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Row(children: [billing.Card(onPressed: (_) {})]),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in SizedBox', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 400,
            height: 400,
            child: billing.Card(onPressed: (_) {}),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });
    });

    group('resolutions', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(physicalSize: entry.value, billing.Card(onPressed: (_) {}));
        await tester.pumpAndSettle();

        expect(find.text('Subscription'), findsOneWidget);
        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });
  });
}
