import 'package:flutter/material.dart';
import 'package:fixnum/fixnum.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/inputs.dart' as inputs;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('Bytes', () {
    testWidgets('initializes with default values and displays correctly', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(inputs.Bytes());
      await tester.pumpAndSettle();

      expect(find.byType(TextField), findsOneWidget);
      expect(find.text('0'), findsOneWidget); // Default initialBytes is 0
      expect(find.byType(DropdownButton<int>), findsOneWidget);
      expect(
        find.text('GiB'),
        findsOneWidget,
      ); // Default initialMagnitude is bytes
    });

    testWidgets('initializes with provided initialBytes and initialMagnitude', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        inputs.Bytes(value: 2 * ds.bytesx.KiB, magnitude: ds.bytesx.KiB),
      );
      await tester.pumpAndSettle();

      expect(find.byType(TextField), findsOneWidget);
      expect(find.text('2'), findsOneWidget);
      expect(find.byType(DropdownButton<int>), findsOneWidget);
      expect(find.text('KiB'), findsOneWidget);
    });

    testWidgets('updates TextField value when magnitude changes', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        inputs.Bytes(value: 5 * ds.bytesx.MiB, magnitude: ds.bytesx.MiB),
      );
      await tester.pumpAndSettle();

      // Verify initial state
      expect(find.text('5'), findsOneWidget);
      expect(find.text('MiB'), findsOneWidget);

      // Tap on the dropdown
      await tester.tap(find.text('MiB'));
      await tester.pumpAndSettle(); // Wait for the dropdown to open

      // Select KiB
      await tester.tap(
        find.text('KiB').last,
      ); // Use .last if multiple KiB are found in the list
      await tester.pumpAndSettle(); // Wait for the dropdown to close and widget to update

      // Verify updated values
      expect(find.text('5'), findsOneWidget);
      expect(find.text('KiB'), findsOneWidget);
    });

    testWidgets(
      'triggers onBytesChanged callback when TextField value changes',
      (WidgetTester tester) async {
        Int64? capturedBytes;
        await tester.pumpApp(
          inputs.Bytes(
            onChange: (bytes) {
              capturedBytes = bytes;
            },
          ),
        );
        await tester.pumpAndSettle();

        // Enter '100' in the TextField
        await tester.enterText(find.byType(TextField), '100');
        await tester.pumpAndSettle();

        expect(capturedBytes, Int64(100 * ds.bytesx.GiB)); // 100 bytes

        // Change magnitude to KiB
        await tester.tap(find.text('GiB'));
        await tester.pumpAndSettle();
        await tester.tap(find.text('KiB').last);
        await tester.pumpAndSettle();

        // Enter '2' in the TextField (now in KiB)
        await tester.enterText(find.byType(TextField), '2');
        await tester.pumpAndSettle();

        expect(capturedBytes, Int64(2 * ds.bytesx.KiB)); // 2 KiB = 2048 bytes
      },
    );

    testWidgets('triggers onBytesChanged callback when magnitude changes', (
      WidgetTester tester,
    ) async {
      Int64? capturedBytes;
      await tester.pumpApp(
        inputs.Bytes(
          value: ds.bytesx.KiB, // 1 KiB
          magnitude: ds.bytesx.KiB,
          onChange: (bytes) {
            capturedBytes = bytes;
          },
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('1'), findsOneWidget); // 1 KiB
      expect(find.text('KiB'), findsOneWidget);
      expect(capturedBytes, isNull);

      // Tap on the dropdown
      await tester.tap(find.text('KiB'));
      await tester.pumpAndSettle();

      // Select MiB
      await tester.tap(find.text('MiB').last);
      await tester.pumpAndSettle();

      // The value should update to 0.00 in the TextField but the underlying bytes should remain 1024.
      // The callback should reflect the underlying bytes value.
      expect(
        find.text('1'),
        findsOneWidget,
      ); // 1 KiB is 0.00 MiB (rounded to 2 decimal places)
      expect(find.text('MiB'), findsOneWidget);
      expect(capturedBytes, Int64(ds.bytesx.MiB));
    });

    testWidgets('handles non-numeric input gracefully', (
      WidgetTester tester,
    ) async {
      Int64? capturedBytes;
      await tester.pumpApp(
        inputs.Bytes(
          onChange: (bytes) {
            capturedBytes = bytes;
          },
        ),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), 'abc');
      await tester.pumpAndSettle();

      expect(
        find.text('abc'),
        findsNothing,
      ); // TextField still shows invalid input
      expect(capturedBytes, isNull);
    });

    testWidgets('uses provided decoration', (WidgetTester tester) async {
      const customLabelText = 'Enter Capacity';
      await tester.pumpApp(
        inputs.Bytes(
          decoration: InputDecoration(
            labelText: customLabelText,
            hintText: 'e.g., 100',
            border: OutlineInputBorder(
              borderRadius: BorderRadius.all(Radius.circular(10.0)),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text(customLabelText), findsOneWidget);
      expect(find.byType(TextField), findsOneWidget);

      final TextField textField = tester.widget(find.byType(TextField));
      expect(textField.decoration?.hintText, 'e.g., 100');
      expect(textField.decoration?.border, isA<OutlineInputBorder>());
      expect(
        (textField.decoration?.border as OutlineInputBorder).borderRadius,
        const BorderRadius.all(Radius.circular(10.0)),
      );
    });
  });
}
