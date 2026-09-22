import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/profiles/create.inlined.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

void main() {
  group('profiles.CreateInlined', () {
    group('layout', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        await tester.pumpApp(
          CreateInlined(onClose: () async {}),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
        expect(find.byType(CreateInlined), findsOneWidget);
      });

      testWidgets('renders without overflow in constrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 400, maxHeight: 600),
            child: SingleChildScrollView(
              child: CreateInlined(onClose: () async {}),
            ),
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
          SingleChildScrollView(
            child: CreateInlined(onClose: () async {}),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });

    group('validation', () {
      testWidgets('submitting with an empty public key shows an error and does not close', (
        WidgetTester tester,
      ) async {
        var closed = false;

        await tester.pumpApp(
          CreateInlined(onClose: () async => closed = true),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.widgetWithText(OutlinedButton, 'Add'));
        await tester.pumpAndSettle();

        expect(find.text('Public key is required'), findsOneWidget);
        expect(closed, isFalse);
        expect(tester.takeException(), isNull);
      });
    });

    group('cancel', () {
      testWidgets('tapping Cancel calls onClose', (WidgetTester tester) async {
        var closed = false;

        await tester.pumpApp(
          CreateInlined(onClose: () async => closed = true),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.widgetWithText(OutlinedButton, 'Cancel'));
        await tester.pumpAndSettle();

        expect(closed, isTrue);
        expect(tester.takeException(), isNull);
      });
    });
  });
}
