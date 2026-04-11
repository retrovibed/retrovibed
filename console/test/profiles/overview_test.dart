import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/profiles.dart' as profiles;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

void main() {
  group('Overview', () {
    group('layout', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        await tester.pumpApp(
          const SingleChildScrollView(child: profiles.Overview()),
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
            child: const SingleChildScrollView(child: profiles.Overview()),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in Column', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          const SingleChildScrollView(
            child: Column(children: [profiles.Overview()]),
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
            width: 600,
            height: 800,
            child: SingleChildScrollView(child: profiles.Overview()),
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
          const SingleChildScrollView(child: profiles.Overview()),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });

    group('content', () {
      testWidgets('renders permission labels from meta display', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          const SingleChildScrollView(child: profiles.Overview()),
        );
        await tester.pumpAndSettle();

        expect(find.text('User Management'), findsWidgets);
        expect(find.text('Library Read'), findsWidgets);
        expect(tester.takeException(), isNull);
      });
    });
  });
}
