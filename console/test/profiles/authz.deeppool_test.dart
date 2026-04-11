import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/profiles.dart' as profiles;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

void main() {
  group('AuthzDeeppool', () {
    group('layout', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        await tester.pumpApp(
          const profiles.AuthzDeeppool(),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in constrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 400, maxHeight: 500),
            child: const profiles.AuthzDeeppool(),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in Column', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          const Column(
            children: [
              profiles.AuthzDeeppool(),
              Spacer(),
            ],
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in SizedBox', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          const SizedBox(
            width: 400,
            height: 400,
            child: profiles.AuthzDeeppool(),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });
    });

    group('resolutions', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          const SingleChildScrollView(child: profiles.AuthzDeeppool()),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });

    group('display values', () {
      testWidgets('renders permission labels', (WidgetTester tester) async {
        await tester.pumpApp(
          const SingleChildScrollView(child: profiles.AuthzDeeppool()),
        );
        await tester.pumpAndSettle();

        expect(find.text('User Management'), findsOneWidget);
        expect(find.text('Library Read'), findsOneWidget);
        expect(find.text('Library Modify'), findsOneWidget);
        expect(find.text('Community Modify'), findsOneWidget);
        expect(find.text('Billing Read'), findsOneWidget);
        expect(find.text('Billing Modify'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('all checkboxes are non-interactive', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(const profiles.AuthzDeeppool());
        await tester.pumpAndSettle();

        final checkboxes = tester.widgetList<Checkbox>(find.byType(Checkbox));
        for (final checkbox in checkboxes) {
          expect(checkbox.onChanged, isNull);
        }
        expect(tester.takeException(), isNull);
      });

      testWidgets('defaults to all permissions disabled when no cache', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(const profiles.AuthzDeeppool());
        await tester.pumpAndSettle();

        final checkboxes = tester.widgetList<Checkbox>(find.byType(Checkbox));
        for (final checkbox in checkboxes) {
          expect(checkbox.value, equals(false));
        }
        expect(tester.takeException(), isNull);
      });
    });
  });
}
