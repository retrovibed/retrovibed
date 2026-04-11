import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/profiles.dart' as profiles;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

meta.Profile _profile({
  String id = 'test-id-1',
  String display = 'Test User',
  String email = 'test@test.com',
}) {
  return meta.Profile(
    id: id,
    display: display,
    email: email,
    updatedAt: '2025-01-01T00:00:00Z',
  );
}

void main() {
  group('Profile Typography Row Tests', () {
    group('default size', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        await tester.pumpApp(
          profiles.Typography(_profile(id: 'id-1', display: 'Alice')),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
        expect(find.text('Alice'), findsOneWidget);
        expect(find.text('id-1'), findsOneWidget);
      });

      testWidgets('renders long names without overflow', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          profiles.Typography(
            _profile(
              id: 'very-long-identifier-that-could-potentially-overflow',
              display: 'A User With An Extremely Long Display Name',
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('shows all columns', (WidgetTester tester) async {
        await tester.pumpApp(
          profiles.Typography(_profile(id: 'id-1', display: 'Alice')),
        );
        await tester.pumpAndSettle();

        expect(find.text('id-1'), findsOneWidget);
        expect(find.text('Alice'), findsOneWidget);
      });
    });

    group('compact layout (< 400px)', () {
      testWidgets('hides id and updated columns below 400px', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          profiles.Typography(_profile(id: 'id-1', display: 'Alice')),
          physicalSize: const Size(380, 600),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
        expect(find.text('id-1'), findsNothing);
        expect(find.text('Alice'), findsOneWidget);
      });

      testWidgets('renders without overflow at 200px with long name', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          profiles.Typography(
            _profile(
              id: 'very-long-identifier',
              display: 'A User With An Extremely Long Display Name',
            ),
          ),
          physicalSize: const Size(200, 600),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
        expect(find.text('very-long-identifier'), findsNothing);
      });

      testWidgets('renders without overflow at 150px', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          profiles.Typography(_profile(id: 'id-1', display: 'Alice')),
          physicalSize: const Size(150, 600),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });
    });

    group('shows all columns above breakpoint', () {
      testWidgets('shows all columns at exactly 400px', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          profiles.Typography(_profile(id: 'id-1', display: 'Alice')),
          physicalSize: const Size(400, 600),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
        expect(find.text('id-1'), findsOneWidget);
        expect(find.text('Alice'), findsOneWidget);
      });

      testWidgets('shows all columns at 400px', (WidgetTester tester) async {
        await tester.pumpApp(
          profiles.Typography(_profile(id: 'id-1', display: 'Alice')),
          physicalSize: const Size(400, 600),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
        expect(find.text('id-1'), findsOneWidget);
        expect(find.text('Alice'), findsOneWidget);
      });
    });

    group('resolutions', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          profiles.Typography(_profile(id: 'id-1', display: 'Alice')),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      }, variant: _resolutions);

      testWidgets('renders without overflow with long name', (
        WidgetTester tester,
      ) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          profiles.Typography(
            _profile(
              id: 'very-long-identifier',
              display: 'A User With An Extremely Long Display Name',
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });
  });
}
