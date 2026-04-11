import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/lucene.dart' as lucene;
import 'package:retrovibed/timex.dart' as timex;

void main() {
  group('Queryer renders', () {
    testWidgets('shows search hint text', (tester) async {
      await tester.pumpApp(lucene.Queryer((_) {}, []));
      await tester.pumpAndSettle();

      expect(find.text('Search… (@ for filters)'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders no chips initially', (tester) async {
      await tester.pumpApp(
        lucene.Queryer((_) {}, [lucene.Boolean('hd', false, false, (_) {})]),
      );
      await tester.pumpAndSettle();

      expect(find.byType(lucene.FilterChip), findsNothing);
      expect(tester.takeException(), isNull);
    });
  });

  group('Queryer autocomplete dropdown', () {
    testWidgets('typing @ shows field suggestions', (tester) async {
      await tester.pumpApp(
        lucene.Queryer((_) {}, [
          lucene.Boolean('hd', false, false, (_) {}),
          lucene.Boolean('completed', false, false, (_) {}),
        ]),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), '@');
      await tester.pump();

      expect(find.text('hd'), findsOneWidget);
      expect(find.text('completed'), findsOneWidget);
    });

    testWidgets('typing @h narrows to matching fields', (tester) async {
      await tester.pumpApp(
        lucene.Queryer((_) {}, [
          lucene.Boolean('hd', false, false, (_) {}),
          lucene.Boolean('completed', false, false, (_) {}),
        ]),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), '@h');
      await tester.pump();

      expect(find.text('hd'), findsOneWidget);
      expect(find.text('completed'), findsNothing);
    });

    testWidgets('no matching prefix shows no dropdown', (tester) async {
      await tester.pumpApp(
        lucene.Queryer((_) {}, [lucene.Boolean('hd', false, false, (_) {})]),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), '@z');
      await tester.pump();

      expect(find.text('hd'), findsNothing);
    });

    testWidgets('escape key closes dropdown', (tester) async {
      await tester.pumpApp(
        lucene.Queryer((_) {}, [lucene.Boolean('hd', false, false, (_) {})]),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), '@');
      await tester.pump();
      expect(find.text('hd'), findsOneWidget);

      await tester.sendKeyEvent(LogicalKeyboardKey.escape);
      await tester.pump();
      expect(find.text('hd'), findsNothing);
    });
  });

  group('Queryer chip management', () {
    testWidgets('tapping Boolean suggestion adds chip', (tester) async {
      await tester.pumpApp(
        lucene.Queryer((_) {}, [lucene.Boolean('hd', false, false, (_) {})]),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), '@');
      await tester.pumpAndSettle();

      await tester.tap(find.text('hd'));
      await tester.pumpAndSettle();

      expect(find.byType(lucene.FilterChip), findsOneWidget);
    });

    testWidgets('tapping Boolean suggestion clears @ token from text field', (
      tester,
    ) async {
      await tester.pumpApp(
        lucene.Queryer((_) {}, [lucene.Boolean('hd', false, false, (_) {})]),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), '@');
      await tester.pumpAndSettle();
      await tester.tap(find.text('hd'));
      await tester.pumpAndSettle();

      final controller = tester.widget<TextField>(find.byType(TextField)).controller;
      expect(controller?.text ?? '', isNot(contains('@')));
    });

    testWidgets('deleting a chip removes it', (tester) async {
      await tester.pumpApp(
        lucene.Queryer((_) {}, [lucene.Boolean('hd', false, false, (_) {})]),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), '@');
      await tester.pumpAndSettle();
      await tester.tap(find.text('hd'));
      await tester.pumpAndSettle();
      expect(find.byType(lucene.FilterChip), findsOneWidget);

      await tester.tap(find.byIcon(Icons.close));
      await tester.pumpAndSettle();

      expect(find.byType(lucene.FilterChip), findsNothing);
    });

    testWidgets('deleted chip field reappears in suggestions', (tester) async {
      await tester.pumpApp(
        lucene.Queryer((_) {}, [lucene.Boolean('hd', false, false, (_) {})]),
      );
      await tester.pumpAndSettle();

      // Add the chip.
      await tester.enterText(find.byType(TextField), '@');
      await tester.pumpAndSettle();
      await tester.tap(find.text('hd'));
      await tester.pumpAndSettle();
      expect(find.byType(lucene.FilterChip), findsOneWidget);

      // Delete the chip.
      await tester.tap(find.byIcon(Icons.close));
      await tester.pumpAndSettle();
      expect(find.byType(lucene.FilterChip), findsNothing);

      // Field should be available again in suggestions.
      await tester.enterText(find.byType(TextField), '@');
      await tester.pumpAndSettle();
      expect(find.text('hd'), findsOneWidget);
    });

    testWidgets(
      'committed non-autocomplete field does not reappear in suggestions',
      (tester) async {
        await tester.pumpApp(
          lucene.Queryer((_) {}, [
            lucene.Elapsed.auto('runtime', Duration.zero, (_) {}),
            lucene.Boolean('hd', false, false, (_) {}),
          ]),
        );
        await tester.pumpAndSettle();

        // Commit runtime by typing @runtime: then selecting the first value
        // suggestion (Enter on highlighted first preset).
        await tester.enterText(find.byType(TextField), '@runtime:');
        await tester.pumpAndSettle();
        await tester.sendKeyEvent(LogicalKeyboardKey.enter);
        await tester.pumpAndSettle();
        expect(find.byType(lucene.FilterChip), findsOneWidget);

        // Re-open suggestions — runtime must not appear.
        await tester.enterText(find.byType(TextField), '@');
        await tester.pumpAndSettle();

        expect(find.text('hd'), findsOneWidget);
        expect(
          tester.widgetList<ListTile>(find.byType(ListTile)).map((t) {
            final text = t.title;
            return text is Text ? text.data : null;
          }),
          isNot(contains('runtime')),
        );
      },
    );

    testWidgets('each field can only appear once as a chip', (tester) async {
      await tester.pumpApp(
        lucene.Queryer((_) {}, [
          lucene.Boolean('hd', false, false, (_) {}),
          lucene.Boolean('completed', false, false, (_) {}),
        ]),
      );
      await tester.pumpAndSettle();

      // Add hd chip.
      await tester.enterText(find.byType(TextField), '@');
      await tester.pump();
      await tester.tap(find.text('hd'));
      await tester.pumpAndSettle();
      expect(find.text('hd'), findsOneWidget);
      await tester.pumpAndSettle();

      // Open dropdown again — hd should not appear.
      await tester.enterText(find.byType(TextField), '@');
      await tester.pumpAndSettle();

      // 'completed' is the only suggestion; 'hd' appears only in its chip.
      expect(find.text('completed'), findsOneWidget);
      expect(find.text('hd'), findsOneWidget); // chip label, not a suggestion
    });
  });

  group('Queryer onQuery callback', () {
    testWidgets('submitting free text calls onQuery', (tester) async {
      String? emitted;
      await tester.pumpApp(lucene.Queryer((q) => emitted = q, []));
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), 'hello world');
      await tester.testTextInput.receiveAction(TextInputAction.done);
      await tester.pump();

      expect(emitted, 'hello world');
    });

    testWidgets('query combines chips and free text', (tester) async {
      String emitted = "";
      await tester.pumpApp(
        lucene.Queryer((q) => emitted += q, [
          lucene.Boolean.auto('hd', false, (v) {
            emitted += "hd:true|";
          }),
        ]),
      );
      await tester.pumpAndSettle();

      // Add chip.
      await tester.enterText(find.byType(TextField), '@');
      await tester.pumpAndSettle();
      await tester.tap(find.text('hd'));
      await tester.pumpAndSettle();

      // Type free text and submit.
      await tester.enterText(find.byType(TextField), 'batman');
      await tester.testTextInput.receiveAction(TextInputAction.done);
      await tester.pumpAndSettle();

      expect(emitted, equals('hd:true|batman'));
    });
  });

  group('Queryer field setters', () {
    testWidgets('adding Boolean chip calls setter with true', (tester) async {
      bool? captured;
      await tester.pumpApp(
        lucene.Queryer((_) {}, [
          lucene.Boolean('hd', false, false, (v) => captured = v),
        ]),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), '@');
      await tester.pumpAndSettle();
      await tester.tap(find.text('hd'));
      await tester.pumpAndSettle();

      expect(captured, isTrue);
    });

    testWidgets('removing Boolean chip resets setter to defaultValue', (
      tester,
    ) async {
      bool? captured;
      await tester.pumpApp(
        lucene.Queryer((_) {}, [
          lucene.Boolean('hd', false, false, (v) => captured = v),
        ]),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), '@');
      await tester.pumpAndSettle();
      await tester.tap(find.text('hd'));
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.close));
      await tester.pumpAndSettle();

      expect(captured, isFalse); // reset to defaultValue
    });

    testWidgets('free text submit does not call field setters', (tester) async {
      bool? captured;
      await tester.pumpApp(
        lucene.Queryer((_) {}, [
          lucene.Boolean('hd', false, false, (v) => captured = v),
        ]),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), 'batman');
      await tester.testTextInput.receiveAction(TextInputAction.done);
      await tester.pump();

      expect(captured, isNull);
    });
  });

  group('Queryer keyboard navigation', () {
    testWidgets('first suggestion is highlighted by default', (tester) async {
      await tester.pumpApp(
        lucene.Queryer((_) {}, [
          lucene.Boolean('hd', false, false, (_) {}),
          lucene.Boolean('completed', false, false, (_) {}),
        ]),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), '@');
      await tester.pump();

      final tiles = tester.widgetList<ListTile>(find.byType(ListTile)).toList();
      expect(tiles.first.selected, isTrue);
      expect(tiles.last.selected, isFalse);
      expect(tester.takeException(), isNull);
    });

    testWidgets('arrow down cycles to next suggestion', (tester) async {
      await tester.pumpApp(
        lucene.Queryer((_) {}, [
          lucene.Boolean('hd', false, false, (_) {}),
          lucene.Boolean('completed', false, false, (_) {}),
        ]),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), '@');
      await tester.pump();

      await tester.sendKeyEvent(LogicalKeyboardKey.arrowDown);
      await tester.pump();

      final tiles = tester.widgetList<ListTile>(find.byType(ListTile)).toList();
      expect(tiles.first.selected, isFalse);
      expect(tiles.last.selected, isTrue);
      expect(tester.takeException(), isNull);
    });

    testWidgets('arrow up wraps to last suggestion from first', (tester) async {
      await tester.pumpApp(
        lucene.Queryer((_) {}, [
          lucene.Boolean('hd', false, false, (_) {}),
          lucene.Boolean('completed', false, false, (_) {}),
        ]),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), '@');
      await tester.pump();

      await tester.sendKeyEvent(LogicalKeyboardKey.arrowUp);
      await tester.pump();

      final tiles = tester.widgetList<ListTile>(find.byType(ListTile)).toList();
      expect(tiles.first.selected, isFalse);
      expect(tiles.last.selected, isTrue);
      expect(tester.takeException(), isNull);
    });

    testWidgets('arrow down wraps back to first suggestion from last', (tester) async {
      await tester.pumpApp(
        lucene.Queryer((_) {}, [
          lucene.Boolean('hd', false, false, (_) {}),
          lucene.Boolean('completed', false, false, (_) {}),
        ]),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), '@');
      await tester.pump();

      // Two arrow downs wraps back to first.
      await tester.sendKeyEvent(LogicalKeyboardKey.arrowDown);
      await tester.sendKeyEvent(LogicalKeyboardKey.arrowDown);
      await tester.pump();

      final tiles = tester.widgetList<ListTile>(find.byType(ListTile)).toList();
      expect(tiles.first.selected, isTrue);
      expect(tester.takeException(), isNull);
    });

    testWidgets('enter selects the highlighted suggestion', (tester) async {
      bool? captured;
      await tester.pumpApp(
        lucene.Queryer((_) {}, [
          lucene.Boolean('hd', false, false, (v) => captured = v),
          lucene.Boolean('completed', false, false, (_) {}),
        ]),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), '@');
      await tester.pumpAndSettle();

      // First suggestion (hd) is highlighted by default — press Enter to select.
      await tester.sendKeyEvent(LogicalKeyboardKey.enter);
      await tester.pumpAndSettle();

      expect(find.byType(lucene.FilterChip), findsOneWidget);
      expect(captured, isTrue);
    });

    testWidgets('enter selects second suggestion after arrow down', (tester) async {
      bool? captured;
      await tester.pumpApp(
        lucene.Queryer((_) {}, [
          lucene.Boolean('hd', false, false, (_) {}),
          lucene.Boolean('completed', false, false, (v) => captured = v),
        ]),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), '@');
      await tester.pumpAndSettle();

      await tester.sendKeyEvent(LogicalKeyboardKey.arrowDown);
      await tester.pumpAndSettle();
      await tester.sendKeyEvent(LogicalKeyboardKey.enter);
      await tester.pumpAndSettle();

      expect(find.byType(lucene.FilterChip), findsOneWidget);
      expect(captured, isTrue);
    });

    testWidgets('highlight resets to first when suggestions change', (
      tester,
    ) async {
      await tester.pumpApp(
        lucene.Queryer((_) {}, [
          lucene.Boolean('hd', false, false, (_) {}),
          lucene.Boolean('completed', false, false, (_) {}),
        ]),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), '@');
      await tester.pump();
      await tester.sendKeyEvent(LogicalKeyboardKey.tab);
      await tester.pump();

      // Narrow the list — second item disappears, index should reset.
      await tester.enterText(find.byType(TextField), '@h');
      await tester.pump();

      final tiles = tester.widgetList<ListTile>(find.byType(ListTile)).toList();
      expect(tiles.first.selected, isTrue);
      expect(tester.takeException(), isNull);
    });

    testWidgets('enter submits query when no suggestions are shown', (
      tester,
    ) async {
      String? emitted;
      await tester.pumpApp(lucene.Queryer((q) => emitted = q, []));
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), 'hello world');
      await tester.testTextInput.receiveAction(TextInputAction.done);
      await tester.pumpAndSettle();

      expect(emitted, 'hello world');
      expect(tester.takeException(), isNull);
    });
  });

  group('Queryer Value state Timestamp', () {
    testWidgets('pressing enter after picking a date commits chip', (tester) async {
      await tester.pumpApp(
        lucene.Queryer((_) {}, [
          lucene.Timestamp.auto('published', timex.epoch, (_) {}),
        ]),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), '@published:');
      await tester.pumpAndSettle();
      expect(find.byType(CalendarDatePicker), findsOneWidget);

      // Pick a date.
      tester.widget<CalendarDatePicker>(find.byType(CalendarDatePicker)).onDateChanged.call(DateTime(2025, 6, 1));
      await tester.pump();

      // Enter should commit the pending value as a chip.
      await tester.sendKeyEvent(LogicalKeyboardKey.enter);
      await tester.pumpAndSettle();

      expect(find.byType(lucene.FilterChip), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('enter is not consumed by queryer when no suggestions shown', (tester) async {
      await tester.pumpApp(
        lucene.Queryer((_) {}, [
          lucene.Timestamp.auto('published', timex.epoch, (_) {}),
        ]),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), '@published:');
      await tester.pumpAndSettle();
      expect(find.byType(CalendarDatePicker), findsOneWidget);
      expect(find.byType(YearPicker), findsNothing);

      // Tap the mode toggle (the month/year label InkWell) to open year picker.
      // If the queryer were consuming Enter, the year picker would not open.
      await tester.tap(find.bySemanticsLabel(RegExp(r'Select year')));
      await tester.pumpAndSettle();
      expect(find.byType(YearPicker), findsOneWidget);

      // Press Enter inside the year picker — queryer must not consume it.
      await tester.sendKeyEvent(LogicalKeyboardKey.enter);
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });
  });

  group('Queryer chip editing', () {
    testWidgets('Timestamp edit panel stays open after picking a date', (tester) async {
      await tester.pumpApp(
        lucene.Queryer((_) {}, [
          lucene.Timestamp.auto('published', timex.epoch, (_) {}),
        ]),
      );
      await tester.pumpAndSettle();

      // Commit a Timestamp chip.
      await tester.enterText(find.byType(TextField), '@published:');
      await tester.pumpAndSettle();
      tester.widget<CalendarDatePicker>(find.byType(CalendarDatePicker)).onDateChanged.call(DateTime(2025, 6, 1));
      await tester.pump();
      final ctrl = tester.widget<TextField>(find.byType(TextField)).controller!;
      ctrl.selection = const TextSelection.collapsed(offset: 0);
      await tester.pumpAndSettle();
      expect(find.byType(lucene.FilterChip), findsOneWidget);

      // Open the edit panel via chip press.
      await tester.tap(find.byType(lucene.FilterChip));
      await tester.pumpAndSettle();
      expect(find.byType(CalendarDatePicker), findsOneWidget);

      // Pick a different date — the panel must remain open.
      tester.widget<CalendarDatePicker>(find.byType(CalendarDatePicker)).onDateChanged.call(DateTime(2025, 7, 15));
      await tester.pump();

      expect(find.byType(CalendarDatePicker), findsOneWidget, reason: 'edit panel must stay open after picking a date');
      expect(tester.takeException(), isNull);
    });

    testWidgets('enter accepts date and closes Timestamp edit panel', (tester) async {
      await tester.pumpApp(
        lucene.Queryer((_) {}, [
          lucene.Timestamp.auto('published', timex.epoch, (_) {}),
        ]),
      );
      await tester.pumpAndSettle();

      // Commit a Timestamp chip.
      await tester.enterText(find.byType(TextField), '@published:');
      await tester.pumpAndSettle();
      tester.widget<CalendarDatePicker>(find.byType(CalendarDatePicker)).onDateChanged.call(DateTime(2025, 6, 1));
      await tester.pump();
      final ctrl = tester.widget<TextField>(find.byType(TextField)).controller!;
      ctrl.selection = const TextSelection.collapsed(offset: 0);
      await tester.pumpAndSettle();
      expect(find.byType(lucene.FilterChip), findsOneWidget);

      // Open the edit panel via chip press.
      await tester.tap(find.byType(lucene.FilterChip));
      await tester.pumpAndSettle();
      expect(find.byType(CalendarDatePicker), findsOneWidget);

      // Pick a different date.
      tester.widget<CalendarDatePicker>(find.byType(CalendarDatePicker)).onDateChanged.call(DateTime(2025, 7, 15));
      await tester.pump();

      // Press Enter — should accept the picked date and close the panel.
      await tester.sendKeyEvent(LogicalKeyboardKey.enter);
      await tester.pumpAndSettle();

      expect(
        find.byType(CalendarDatePicker),
        findsNothing,
        reason: 'edit panel must close after Enter accepts the date',
      );
      expect(find.byType(lucene.FilterChip), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('tab after closing Timestamp edit panel does not throw', (tester) async {
      await tester.pumpApp(
        lucene.Queryer((_) {}, [
          lucene.Timestamp.auto('published', timex.epoch, (_) {}),
        ]),
      );
      await tester.pumpAndSettle();

      // Show the Timestamp render widget (CalendarDatePicker) by typing @published:
      await tester.enterText(find.byType(TextField), '@published:');
      await tester.pumpAndSettle();
      expect(find.byType(CalendarDatePicker), findsOneWidget);

      // Pick a date — DateInput buffers and commits via deactivate when the
      // Value state's render widget is removed. Clearing the field exits Value.
      tester.widget<CalendarDatePicker>(find.byType(CalendarDatePicker)).onDateChanged.call(DateTime(2025, 6, 1));
      await tester.pump();
      // Exit Value state by moving cursor before the @ anchor.
      final ctrl = tester.widget<TextField>(find.byType(TextField)).controller!;
      ctrl.selection = const TextSelection.collapsed(offset: 0);
      await tester.pumpAndSettle();
      expect(find.byType(lucene.FilterChip), findsOneWidget);

      // Open the edit panel via chip press.
      await tester.tap(find.byType(lucene.FilterChip));
      await tester.pumpAndSettle();
      expect(find.byType(CalendarDatePicker), findsOneWidget);

      // Close the edit panel by pressing the chip again (toggle).
      await tester.tap(find.byType(lucene.FilterChip));
      await tester.pumpAndSettle();
      expect(find.byType(CalendarDatePicker), findsNothing);

      // Tab focus traversal must not throw after the panel is removed.
      await tester.sendKeyEvent(LogicalKeyboardKey.tab);
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });
  });
}
