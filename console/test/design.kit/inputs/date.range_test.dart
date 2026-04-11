import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/inputs/date.range.dart';
import 'package:retrovibed/design.kit/typography/timestamp.dart' as typography;
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/timex.dart' as timex;

void main() {
  final begin = DateTime.utc(2027, 1, 15);
  final end = DateTime.utc(2027, 3, 20);

  Widget app(timex.Range value, ValueChanged<timex.Range> onChanged) => Scaffold(
    body: SingleChildScrollView(
      child: DateRangeInput(value: value, onChanged: onChanged),
    ),
  );

  Finder beginButton() => find.ancestor(
    of: find.byWidgetPredicate(
      (w) => w is typography.Timestamp && w.timestamp == begin,
    ),
    matching: find.byType(TextButton),
  );

  Finder endButton(DateTime end) => find.ancestor(
    of: find.byWidgetPredicate(
      (w) => w is typography.Timestamp && w.timestamp == end,
    ),
    matching: find.byType(TextButton),
  );

  group('DateRangeInput renders', () {
    testWidgets('shows Timestamp widgets for begin and end', (tester) async {
      await tester.pumpApp(app(timex.Range(begin, end), (_) {}));
      await tester.pumpAndSettle();

      expect(
        find.byWidgetPredicate(
          (w) => w is typography.Timestamp && w.timestamp == begin,
        ),
        findsOneWidget,
      );
      expect(
        find.byWidgetPredicate(
          (w) => w is typography.Timestamp && w.timestamp == end,
        ),
        findsOneWidget,
      );
      expect(tester.takeException(), isNull);
    });

    testWidgets('shows CalendarDatePicker for begin on init', (tester) async {
      await tester.pumpApp(app(timex.Range(begin, end), (_) {}));
      await tester.pumpAndSettle();

      expect(find.byType(CalendarDatePicker), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('displays never for timex.inf end', (tester) async {
      await tester.pumpApp(app(timex.Range(begin, timex.inf), (_) {}));
      await tester.pumpAndSettle();

      expect(find.text('never'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('DateRangeInput autofocus', () {
    testWidgets('begin button receives focus when autofocus is true', (tester) async {
      await tester.pumpApp(
        Scaffold(
          body: SingleChildScrollView(
            child: DateRangeInput(
              value: timex.Range(begin, end),
              onChanged: (_) {},
              autofocus: true,
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      final focusedContext = tester.binding.focusManager.primaryFocus?.context;
      final beginBtnElement = tester.element(beginButton());
      bool isFocusWithinBeginBtn = false;
      focusedContext?.visitAncestorElements((el) {
        if (el == beginBtnElement) {
          isFocusWithinBeginBtn = true;
          return false;
        }
        return true;
      });

      expect(isFocusWithinBeginBtn, isTrue);
      expect(tester.takeException(), isNull);
    });

    testWidgets('begin button does not steal focus when autofocus is false', (tester) async {
      await tester.pumpApp(
        Scaffold(
          body: SingleChildScrollView(
            child: DateRangeInput(
              value: timex.Range(begin, end),
              onChanged: (_) {},
              autofocus: false,
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      final focusedContext = tester.binding.focusManager.primaryFocus?.context;
      final beginBtnElement = tester.element(beginButton());
      bool isFocusWithinBeginBtn = false;
      focusedContext?.visitAncestorElements((el) {
        if (el == beginBtnElement) {
          isFocusWithinBeginBtn = true;
          return false;
        }
        return true;
      });

      expect(isFocusWithinBeginBtn, isFalse);
      expect(tester.takeException(), isNull);
    });
  });

  group('DateRangeInput tab focus', () {
    // Finds the node whose nearest tappable ancestor is [type].
    FocusNode? _focusedNode(WidgetTester tester) => tester.binding.focusManager.primaryFocus;

    bool _focusIsWithin(WidgetTester tester, Finder finder) {
      final focused = tester.binding.focusManager.primaryFocus?.context;
      final target = tester.element(finder);
      bool found = false;
      focused?.visitAncestorElements((el) {
        if (el == target) {
          found = true;
          return false;
        }
        return true;
      });
      return found;
    }

    testWidgets('tab from begin button moves focus to end button', (tester) async {
      await tester.pumpApp(
        Scaffold(
          body: SingleChildScrollView(
            child: DateRangeInput(
              value: timex.Range(begin, end),
              onChanged: (_) {},
              autofocus: true,
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      // Start: begin button has focus.
      expect(_focusIsWithin(tester, beginButton()), isTrue);

      await tester.sendKeyEvent(LogicalKeyboardKey.tab);
      await tester.pump();

      expect(_focusIsWithin(tester, endButton(end)), isTrue);
      expect(tester.takeException(), isNull);
    });

    // Full tab order: begin → end → month/year header → prev month → next month → day grid → (wrap)
    testWidgets('tab from end button moves focus to prev month arrow', (tester) async {
      await tester.pumpApp(
        Scaffold(
          body: SingleChildScrollView(
            child: DateRangeInput(
              value: timex.Range(begin, end),
              onChanged: (_) {},
              autofocus: true,
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      // begin → end → month/year header → prev month
      for (int i = 0; i < 3; i++) {
        await tester.sendKeyEvent(LogicalKeyboardKey.tab);
        await tester.pump();
      }

      final prevMonthFocused =
          _focusedNode(tester)?.context?.findAncestorWidgetOfExactType<Tooltip>()?.message == 'Previous month';
      expect(prevMonthFocused, isTrue);
      expect(tester.takeException(), isNull);
    });

    testWidgets('tab from prev month moves focus to next month arrow', (tester) async {
      await tester.pumpApp(
        Scaffold(
          body: SingleChildScrollView(
            child: DateRangeInput(
              value: timex.Range(begin, end),
              onChanged: (_) {},
              autofocus: true,
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      // begin → end → month/year header → prev month → next month
      for (int i = 0; i < 4; i++) {
        await tester.sendKeyEvent(LogicalKeyboardKey.tab);
        await tester.pump();
      }

      final nextMonthFocused =
          _focusedNode(tester)?.context?.findAncestorWidgetOfExactType<Tooltip>()?.message == 'Next month';
      expect(nextMonthFocused, isTrue);
      expect(tester.takeException(), isNull);
    });

    testWidgets('tab from next month moves focus into day grid', (tester) async {
      await tester.pumpApp(
        Scaffold(
          body: SingleChildScrollView(
            child: DateRangeInput(
              value: timex.Range(begin, end),
              onChanged: (_) {},
              autofocus: true,
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      // begin → end → month/year header → prev month → next month → day grid
      for (int i = 0; i < 5; i++) {
        await tester.sendKeyEvent(LogicalKeyboardKey.tab);
        await tester.pump();
      }

      final label = _focusedNode(tester)?.debugLabel ?? '';
      expect(
        label.startsWith('Day') || label == 'Day Grid',
        isTrue,
        reason: 'Expected focus on day grid or a day, got: $label',
      );
      expect(tester.takeException(), isNull);
    });
  });

  group('DateRangeInput focus loss submits', () {
    testWidgets('losing focus calls onChanged with pending range', (tester) async {
      timex.Range? captured;
      await tester.pumpApp(
        Scaffold(
          body: Column(
            children: [
              DateRangeInput(
                value: timex.Range(begin, end),
                onChanged: (r) => captured = r,
                autofocus: true,
              ),
              const TextField(key: ValueKey('other')),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      // Pick a new begin date.
      final picked = DateTime(2027, 2, 1);
      tester.widget<CalendarDatePicker>(find.byType(CalendarDatePicker)).onDateChanged(picked);
      await tester.pump();

      // Tap outside to move focus away.
      await tester.tap(find.byKey(const ValueKey('other')));
      await tester.pump();

      expect(captured, isNotNull);
      expect(captured!.begin, equals(picked.toUtc()));
      expect(tester.takeException(), isNull);
    });
  });

  group('DateRangeInput epoch begin', () {
    testWidgets('tapping begin with epoch begin does not throw', (tester) async {
      await tester.pumpApp(
        app(timex.Range(timex.epoch, timex.inf), (_) {}),
      );
      await tester.pumpAndSettle();

      expect(find.byType(CalendarDatePicker), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('DateRangeInput picker', () {
    testWidgets('tapping begin date shows CalendarDatePicker', (tester) async {
      await tester.pumpApp(app(timex.Range(begin, end), (_) {}));
      await tester.pumpAndSettle();

      await tester.tap(beginButton());
      await tester.pump();

      expect(find.byType(CalendarDatePicker), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('tapping end date shows CalendarDatePicker', (tester) async {
      await tester.pumpApp(app(timex.Range(begin, end), (_) {}));
      await tester.pumpAndSettle();

      await tester.tap(endButton(end));
      await tester.pump();

      expect(find.byType(CalendarDatePicker), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('selecting a begin date updates pending but does not call onChanged', (tester) async {
      timex.Range? captured;
      await tester.pumpApp(
        app(timex.Range(begin, end), (r) => captured = r),
      );
      await tester.pumpAndSettle();

      final picked = DateTime(2027, 2, 1);
      tester.widget<CalendarDatePicker>(find.byType(CalendarDatePicker)).onDateChanged(picked);
      await tester.pump();

      expect(captured, isNull);
      expect(find.byType(CalendarDatePicker), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('onChanged is called with pending range on deactivate', (tester) async {
      timex.Range? captured;
      await tester.pumpApp(
        app(timex.Range(begin, end), (r) => captured = r),
      );
      await tester.pumpAndSettle();

      final picked = DateTime(2027, 2, 1);
      tester.widget<CalendarDatePicker>(find.byType(CalendarDatePicker)).onDateChanged(picked);
      await tester.pump();

      await tester.pumpWidget(const SizedBox());

      expect(captured!.begin, equals(picked.toUtc()));
      expect(captured!.end, equals(end));
      expect(tester.takeException(), isNull);
    });

    testWidgets('selecting an end date updates pending but does not call onChanged', (tester) async {
      timex.Range? captured;
      await tester.pumpApp(
        app(timex.Range(begin, end), (r) => captured = r),
      );
      await tester.pumpAndSettle();

      await tester.tap(endButton(end));
      await tester.pump();

      final picked = DateTime(2027, 4, 10);
      tester.widget<CalendarDatePicker>(find.byType(CalendarDatePicker)).onDateChanged(picked);
      await tester.pump();

      expect(captured, isNull);
      expect(find.byType(CalendarDatePicker), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('DateRangeInput deactivate during build', () {
    testWidgets('onChanged via deactivate during a build pass does not throw', (tester) async {
      // Regression: deactivate was calling onChanged synchronously, which
      // triggered setState on an ancestor that was mid-build (e.g. Queryer
      // rebuilding its Value parser state). The fix defers onChanged to a
      // post-frame callback so it never fires during a build pass.
      timex.Range? captured;
      bool built = false;

      late StateSetter outerSetState;
      await tester.pumpWidget(
        MaterialApp(
          home: StatefulBuilder(
            builder: (context, setState) {
              outerSetState = setState;
              return Scaffold(
                body:
                    built
                        ? const SizedBox()
                        : DateRangeInput(
                          value: timex.Range(begin, end),
                          onChanged: (r) {
                            captured = r;
                            // Simulate what Queryer does: call setState on a
                            // parent that may be building.
                            outerSetState(() {});
                          },
                        ),
              );
            },
          ),
        ),
      );
      await tester.pumpAndSettle();

      // Select a new date so _pending != widget.value.
      final picked = DateTime(2027, 2, 1);
      tester.widget<CalendarDatePicker>(find.byType(CalendarDatePicker)).onDateChanged(picked);
      await tester.pump();

      // Trigger a rebuild that deactivates DateRangeInput mid-build.
      outerSetState(() => built = true);
      await tester.pump(); // must not throw

      expect(tester.takeException(), isNull);
      await tester.pump(); // post-frame: deferred onChanged fires
      expect(captured, isNotNull);
      expect(captured!.begin, equals(picked.toUtc()));
    });
  });

  group('DateRangeInput month navigation', () {
    testWidgets('pressing next month arrow does not call onChanged', (tester) async {
      timex.Range? captured;
      await tester.pumpApp(
        app(timex.Range(begin, end), (r) => captured = r),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byTooltip('Next month'));
      await tester.pump();

      expect(captured, isNull);
      expect(tester.takeException(), isNull);
    });

    testWidgets('pressing previous month arrow does not call onChanged', (tester) async {
      timex.Range? captured;
      await tester.pumpApp(
        app(timex.Range(begin, end), (r) => captured = r),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byTooltip('Previous month'));
      await tester.pump();

      expect(captured, isNull);
      expect(tester.takeException(), isNull);
    });

    testWidgets('navigating months preserves the calendar picker', (tester) async {
      await tester.pumpApp(app(timex.Range(begin, end), (_) {}));
      await tester.pumpAndSettle();

      await tester.tap(find.byTooltip('Next month'));
      await tester.pump();

      expect(find.byType(CalendarDatePicker), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('selecting a date after navigating months submits the navigated date on deactivate', (tester) async {
      timex.Range? captured;
      await tester.pumpApp(
        app(timex.Range(begin, end), (r) => captured = r),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byTooltip('Next month'));
      await tester.pump();

      final picked = DateTime(2027, 2, 10);
      tester.widget<CalendarDatePicker>(find.byType(CalendarDatePicker)).onDateChanged(picked);
      await tester.pump();

      await tester.pumpWidget(const SizedBox());

      expect(captured!.begin, equals(picked.toUtc()));
      expect(captured!.end, equals(end));
      expect(tester.takeException(), isNull);
    });
  });
}
