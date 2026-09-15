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

  // The Focus wrapper DateRangeInput itself adds around the calendar picker
  // (autofocus target) — restricted to descendants of DateRangeInput so it
  // can't accidentally match an unrelated Focus widget further up the tree
  // (e.g. the app/route's own focus scope).
  Finder calendarFocusWrapper() => find.ancestor(
    of: find.byType(CalendarDatePicker),
    matching: find.descendant(
      of: find.byType(DateRangeInput),
      matching: find.byType(Focus),
    ),
  ).last;

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
    testWidgets('calendar picker receives focus when autofocus is true', (tester) async {
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
      final calendarWrapperElement = tester.element(calendarFocusWrapper());
      bool isFocusWithinCalendar = focusedContext == calendarWrapperElement;
      focusedContext?.visitAncestorElements((el) {
        if (el == calendarWrapperElement) {
          isFocusWithinCalendar = true;
          return false;
        }
        return true;
      });

      expect(isFocusWithinCalendar, isTrue);
      expect(tester.takeException(), isNull);
    });

    testWidgets('calendar picker does not steal focus when autofocus is false', (tester) async {
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
      final calendarWrapperElement = tester.element(calendarFocusWrapper());
      bool isFocusWithinCalendar = false;
      focusedContext?.visitAncestorElements((el) {
        if (el == calendarWrapperElement) {
          isFocusWithinCalendar = true;
          return false;
        }
        return true;
      });

      expect(isFocusWithinCalendar, isFalse);
      expect(tester.takeException(), isNull);
    });
  });

  group('DateRangeInput tab focus', () {
    // Finds the node whose nearest tappable ancestor is [type].
    FocusNode? _focusedNode(WidgetTester tester) => tester.binding.focusManager.primaryFocus;

    bool _focusIsWithin(WidgetTester tester, Finder finder) {
      final focused = tester.binding.focusManager.primaryFocus?.context;
      final target = tester.element(finder);
      bool found = focused == target;
      focused?.visitAncestorElements((el) {
        if (el == target) {
          found = true;
          return false;
        }
        return true;
      });
      return found;
    }

    testWidgets('tab from calendar (autofocused) moves focus to month/year header', (tester) async {
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

      // Start: the calendar itself has focus (autofocus target).
      expect(_focusIsWithin(tester, calendarFocusWrapper()), isTrue);

      await tester.sendKeyEvent(LogicalKeyboardKey.tab);
      await tester.pump();

      // Not the calendar's own top-level wrapper anymore — moved into its
      // first internal control (the month/year header, no Tooltip).
      final tt = tester.binding.focusManager.primaryFocus?.context?.findAncestorWidgetOfExactType<Tooltip>();
      expect(tt, isNull);
      expect(tester.takeException(), isNull);
    });

    // Full tab order: calendar (autofocus) → month/year header → prev month → next month → day grid → begin → end → (wrap)
    testWidgets('tab from month/year header moves focus to prev month arrow', (tester) async {
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

      // calendar → month/year header → prev month
      for (int i = 0; i < 2; i++) {
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

      // calendar → month/year header → prev month → next month
      for (int i = 0; i < 3; i++) {
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

      // calendar → month/year header → prev month → next month → day grid
      for (int i = 0; i < 4; i++) {
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

  group('DateRangeInput focus retention across begin/end swap', () {
    bool focusIsWithin(WidgetTester tester, Finder finder) {
      final focused = tester.binding.focusManager.primaryFocus?.context;
      final target = tester.element(finder);
      bool found = focused == target;
      focused?.visitAncestorElements((el) {
        if (el == target) {
          found = true;
          return false;
        }
        return true;
      });
      return found;
    }

    testWidgets('switching from begin to end picker keeps focus inside the calendar', (tester) async {
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

      // Sanity check: begin picker starts with focus inside the calendar.
      expect(focusIsWithin(tester, calendarFocusWrapper()), isTrue);

      // Switch from the begin picker to the end picker, exactly as a user
      // clicking the end-date button would.
      await tester.tap(endButton(end));
      await tester.pumpAndSettle();

      // Regression: the calendar Focus wrapper is reused across the swap (its
      // `autofocus` only fires once, on first attach), so focus can drift
      // outside DateRangeInput's own FocusScope once the begin picker's
      // CalendarDatePicker (and its internal focus node) is torn down and
      // replaced by the end picker's.
      expect(
        focusIsWithin(tester, calendarFocusWrapper()),
        isTrue,
        reason: 'focus should stay inside the calendar after switching from the begin to the end picker',
      );
      expect(tester.takeException(), isNull);
    });

    testWidgets('pressing enter after switching to the end picker and choosing a date commits the pending range', (
      tester,
    ) async {
      timex.Range? captured;
      await tester.pumpApp(
        Scaffold(
          body: SingleChildScrollView(
            child: DateRangeInput(
              value: timex.Range(begin, end),
              onChanged: (r) => captured = r,
              autofocus: true,
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      // Manually pick a begin date.
      final pickedBegin = DateTime(2027, 2, 1);
      tester.widget<CalendarDatePicker>(find.byType(CalendarDatePicker)).onDateChanged(pickedBegin);
      await tester.pumpAndSettle();

      // Switch to the end picker and manually pick an end date.
      await tester.tap(endButton(end));
      await tester.pumpAndSettle();
      final pickedEnd = DateTime(2027, 4, 10);
      tester.widget<CalendarDatePicker>(find.byType(CalendarDatePicker)).onDateChanged(pickedEnd);
      await tester.pumpAndSettle();

      // Commit via Enter, as a user would after finishing their manual pick —
      // this should be intercepted by DateRangeInput's own FocusScope, not
      // fall through to some outer handler.
      await tester.sendKeyEvent(LogicalKeyboardKey.enter);
      await tester.pumpAndSettle();

      expect(captured, isNotNull);
      expect(captured!.begin, equals(pickedBegin.toUtc()));
      expect(captured!.end, equals(pickedEnd.toUtc()));
      expect(tester.takeException(), isNull);
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
