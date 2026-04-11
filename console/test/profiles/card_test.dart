import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/profiles.dart' as profiles;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/meta/api.dart' as meta;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

Future<authn.Session> _mockSession(BuildContext ctx) =>
    Future.value(authn.Session());

Future<meta.Authn> _mockCurrent() => Future.value(meta.Authn());

Future<authn.Session> _failSession(BuildContext ctx) =>
    Future.error(Exception('session error'));

Future<meta.Authn> _failCurrent() => Future.error(Exception('current error'));

void main() {
  group('profiles.Card', () {
    group('layout', () {
      testWidgets('renders Account title', (WidgetTester tester) async {
        await tester.pumpApp(
          profiles.Card(session: _mockSession, current: _mockCurrent),
        );
        await tester.pumpAndSettle();

        expect(find.text('Account'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in unconstrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          profiles.Card(session: _mockSession, current: _mockCurrent),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in constrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 400, maxHeight: 600),
            child: profiles.Card(session: _mockSession, current: _mockCurrent),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in Column', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Column(children: [
            Expanded(
              child: profiles.Card(
                session: _mockSession,
                current: _mockCurrent,
              ),
            ),
          ]),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in Row', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Row(children: [
            Expanded(
              child: profiles.Card(
                session: _mockSession,
                current: _mockCurrent,
              ),
            ),
          ]),
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
            child: profiles.Card(session: _mockSession, current: _mockCurrent),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders with margin without overflow', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 400, maxHeight: 600),
            child: profiles.Card(
              margin: const EdgeInsets.all(16),
              session: _mockSession,
              current: _mockCurrent,
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in small constrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 175, maxHeight: 215),
            child: profiles.Card(session: _mockSession, current: _mockCurrent),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in settings grid cell (192x192)', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 192,
            height: 192,
            child: profiles.Card(session: _mockSession, current: _mockCurrent),
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
          profiles.Card(session: _mockSession, current: _mockCurrent),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });

    group('errors', () {
      testWidgets('shows error when session fails', (WidgetTester tester) async {
        await tester.pumpApp(
          profiles.Card(session: _failSession, current: _mockCurrent),
        );
        await tester.pumpAndSettle();

        expect(find.text('an unexpected problem has occurred'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('shows error when current fails', (WidgetTester tester) async {
        await tester.pumpApp(
          profiles.Card(session: _mockSession, current: _failCurrent),
        );
        await tester.pumpAndSettle();

        expect(find.text('an unexpected problem has occurred'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('shows error when both fail', (WidgetTester tester) async {
        await tester.pumpApp(
          profiles.Card(session: _failSession, current: _failCurrent),
        );
        await tester.pumpAndSettle();

        expect(find.text('an unexpected problem has occurred'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('error renders without overflow in constrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 300, maxHeight: 400),
            child: profiles.Card(session: _failSession, current: _failCurrent),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('error renders without overflow in small constrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 175, maxHeight: 215),
            child: profiles.Card(session: _failSession, current: _failCurrent),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });
    });
  });
}
