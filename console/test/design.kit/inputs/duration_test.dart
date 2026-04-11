import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/design.kit/inputs/duration.dart';

Future<void> _submit(WidgetTester tester) async {
  await tester.testTextInput.receiveAction(TextInputAction.done);
  await tester.pump();
}

void main() {
  group('DurationInput', () {
    group('initial display', () {
      testWidgets('renders a TextFormField and unit dropdown', (tester) async {
        await tester.pumpApp(
          DurationInput(value: Duration.zero, onChanged: (_) {}),
        );
        await tester.pumpAndSettle();

        expect(find.byType(TextFormField), findsOneWidget);
        expect(
          find.byType(DropdownButton<({String label, Duration duration})>),
          findsOneWidget,
        );
        expect(tester.takeException(), isNull);
      });

      testWidgets('shows s unit and correct amount for a seconds duration', (
        tester,
      ) async {
        await tester.pumpApp(
          DurationInput(
            value: const Duration(seconds: 45),
            onChanged: (_) {},
          ),
        );
        await tester.pumpAndSettle();

        final field = tester.widget<TextFormField>(find.byType(TextFormField));
        expect(field.controller?.text ?? field.initialValue, '45');
        final unit =
            tester
                .widget<DropdownButton<({String label, Duration duration})>>(
                  find.byType(DropdownButton<({String label, Duration duration})>),
                )
                .value!
                .label;
        expect(unit, 's');
      });

      testWidgets('shows m unit for a whole-minute duration', (tester) async {
        await tester.pumpApp(
          DurationInput(
            value: const Duration(minutes: 90),
            onChanged: (_) {},
          ),
        );
        await tester.pumpAndSettle();

        final field = tester.widget<TextFormField>(find.byType(TextFormField));
        expect(field.controller?.text ?? field.initialValue, '90');
        final unit =
            tester
                .widget<DropdownButton<({String label, Duration duration})>>(
                  find.byType(DropdownButton<({String label, Duration duration})>),
                )
                .value!
                .label;
        expect(unit, 'm');
      });

      testWidgets('shows hr unit for a whole-hour duration', (tester) async {
        await tester.pumpApp(
          DurationInput(
            value: const Duration(hours: 2),
            onChanged: (_) {},
          ),
        );
        await tester.pumpAndSettle();

        final field = tester.widget<TextFormField>(find.byType(TextFormField));
        expect(field.controller?.text ?? field.initialValue, '2');
        final unit =
            tester
                .widget<DropdownButton<({String label, Duration duration})>>(
                  find.byType(DropdownButton<({String label, Duration duration})>),
                )
                .value!
                .label;
        expect(unit, 'hr');
      });

      testWidgets('shows empty text field for zero duration', (tester) async {
        await tester.pumpApp(
          DurationInput(value: Duration.zero, onChanged: (_) {}),
        );
        await tester.pumpAndSettle();

        final field = tester.widget<TextFormField>(find.byType(TextFormField));
        expect(field.controller?.text ?? field.initialValue, '');
      });
    });

    group('onChanged callback', () {
      testWidgets('fires with correct Duration on submit in s', (tester) async {
        Duration? captured;
        await tester.pumpApp(
          DurationInput(
            value: Duration.zero,
            onChanged: (d) => captured = d,
          ),
        );
        await tester.pumpAndSettle();

        await tester.enterText(find.byType(TextFormField), '45');
        await _submit(tester);

        expect(captured, const Duration(seconds: 45));
      });

      testWidgets('fires with correct Duration on submit in m', (tester) async {
        Duration? captured;
        await tester.pumpApp(
          DurationInput(
            value: const Duration(minutes: 1),
            onChanged: (d) => captured = d,
          ),
        );
        await tester.pumpAndSettle();

        await tester.enterText(find.byType(TextFormField), '30');
        await _submit(tester);

        expect(captured, const Duration(minutes: 30));
      });

      testWidgets('fires in hours on submit when unit is hr', (tester) async {
        Duration? captured;
        await tester.pumpApp(
          DurationInput(
            value: const Duration(hours: 1),
            onChanged: (d) => captured = d,
          ),
        );
        await tester.pumpAndSettle();

        await tester.enterText(find.byType(TextFormField), '3');
        await _submit(tester);

        expect(captured, const Duration(hours: 3));
      });

      testWidgets('does not fire before submit', (tester) async {
        int callCount = 0;
        await tester.pumpApp(
          DurationInput(
            value: Duration.zero,
            onChanged: (_) => callCount++,
          ),
        );
        await tester.pumpAndSettle();

        await tester.enterText(find.byType(TextFormField), '45');
        await tester.pump();

        expect(callCount, 0);
      });

      testWidgets('does not fire for non-numeric input', (tester) async {
        int callCount = 0;
        await tester.pumpApp(
          DurationInput(
            value: Duration.zero,
            onChanged: (_) => callCount++,
          ),
        );
        await tester.pumpAndSettle();

        // FilteringTextInputFormatter blocks non-digits.
        await tester.enterText(find.byType(TextFormField), 'abc');
        await _submit(tester);

        expect(callCount, 0);
        expect(find.text('abc'), findsNothing);
      });

      testWidgets('does not fire for zero input', (tester) async {
        int callCount = 0;
        await tester.pumpApp(
          DurationInput(
            value: Duration.zero,
            onChanged: (_) => callCount++,
          ),
        );
        await tester.pumpAndSettle();

        await tester.enterText(find.byType(TextFormField), '0');
        await _submit(tester);

        expect(callCount, 0);
      });
    });

    group('unit switching', () {
      testWidgets('dropdown shows s, m, and hr options', (tester) async {
        await tester.pumpApp(
          DurationInput(value: Duration.zero, onChanged: (_) {}),
        );
        await tester.pumpAndSettle();

        await tester.tap(
          find.byType(DropdownButton<({String label, Duration duration})>),
        );
        await tester.pumpAndSettle();

        // 3 menu items + 1 selected item rendered in the button = 4 total.
        expect(
          find.byType(DropdownMenuItem<({String label, Duration duration})>),
          findsNWidgets(4),
        );
      });

      testWidgets('switching unit updates the dropdown label', (tester) async {
        await tester.pumpApp(
          DurationInput(
            value: const Duration(minutes: 90),
            onChanged: (_) {},
          ),
        );
        await tester.pumpAndSettle();

        final before =
            tester
                .widget<DropdownButton<({String label, Duration duration})>>(
                  find.byType(DropdownButton<({String label, Duration duration})>),
                )
                .value!
                .label;
        expect(before, 'm');

        await tester.tap(
          find.byType(DropdownButton<({String label, Duration duration})>),
        );
        await tester.pumpAndSettle();
        await tester.tap(find.text('hr').last);
        await tester.pumpAndSettle();

        final after =
            tester
                .widget<DropdownButton<({String label, Duration duration})>>(
                  find.byType(DropdownButton<({String label, Duration duration})>),
                )
                .value!
                .label;
        expect(after, 'hr');
      });

      testWidgets('entering value then switching unit submits with new unit', (
        tester,
      ) async {
        Duration? captured;
        await tester.pumpApp(
          DurationInput(
            value: Duration.zero,
            onChanged: (d) => captured = d,
          ),
        );
        await tester.pumpAndSettle();

        // Type and submit 30 in s.
        await tester.enterText(find.byType(TextFormField), '30');
        await _submit(tester);
        expect(captured, const Duration(seconds: 30));

        // Switch to m.
        await tester.tap(
          find.byType(DropdownButton<({String label, Duration duration})>),
        );
        await tester.pumpAndSettle();
        await tester.tap(find.text('m').last);
        await tester.pumpAndSettle();

        // Type and submit 5 in m.
        await tester.enterText(find.byType(TextFormField), '5');
        await _submit(tester);
        expect(captured, const Duration(minutes: 5));
      });
    });
  });
}
