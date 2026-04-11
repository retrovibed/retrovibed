import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('Shortcuts', () {
    testWidgets('calls callback when activator matches', (tester) async {
      var called = false;
      await tester.pumpApp(
        ds.Shortcuts(
          Focus(autofocus: true, child: const SizedBox()),
          bindings: {
            const SingleActivator(LogicalKeyboardKey.escape): (
              const SizedBox(),
              () {
                called = true;
                return KeyEventResult.handled;
              },
            ),
          },
        ),
      );
      await tester.pump();
      await tester.pumpAndSettle();

      await tester.sendKeyEvent(LogicalKeyboardKey.escape);
      expect(called, isTrue);
    });

    testWidgets('does not call callback when activator does not match', (
      tester,
    ) async {
      var called = false;
      await tester.pumpApp(
        ds.Shortcuts(
          Focus(autofocus: true, child: const SizedBox()),
          bindings: {
            const SingleActivator(LogicalKeyboardKey.escape): (
              const SizedBox(),
              () {
                called = true;
                return KeyEventResult.handled;
              },
            ),
          },
        ),
      );
      await tester.pump();
      await tester.pumpAndSettle();

      await tester.sendKeyEvent(LogicalKeyboardKey.arrowDown);
      expect(called, isFalse);
    });

    testWidgets('returning handled consumes the event', (tester) async {
      var outerCalled = false;
      await tester.pumpApp(
        Focus(
          onKeyEvent: (node, event) {
            if (event is KeyDownEvent &&
                event.logicalKey == LogicalKeyboardKey.escape) {
              outerCalled = true;
            }
            return KeyEventResult.ignored;
          },
          child: ds.Shortcuts(
            Focus(autofocus: true, child: const SizedBox()),
            bindings: {
              const SingleActivator(LogicalKeyboardKey.escape): (
                const SizedBox(),
                () => KeyEventResult.handled,
              ),
            },
          ),
        ),
      );
      await tester.pump();
      await tester.pumpAndSettle();

      await tester.sendKeyEvent(LogicalKeyboardKey.escape);
      expect(outerCalled, isFalse);
    });

    testWidgets('returning ignored allows event to propagate', (tester) async {
      var outerCalled = false;
      await tester.pumpApp(
        Focus(
          onKeyEvent: (node, event) {
            if (event is KeyDownEvent &&
                event.logicalKey == LogicalKeyboardKey.escape) {
              outerCalled = true;
            }
            return KeyEventResult.ignored;
          },
          child: ds.Shortcuts(
            Focus(autofocus: true, child: const SizedBox()),
            bindings: {
              const SingleActivator(LogicalKeyboardKey.escape): (
                const SizedBox(),
                () => KeyEventResult.ignored,
              ),
            },
          ),
        ),
      );
      await tester.pump();
      await tester.pumpAndSettle();

      await tester.sendKeyEvent(LogicalKeyboardKey.escape);
      expect(outerCalled, isTrue);
    });

    testWidgets('does not request focus itself', (tester) async {
      final focusNode = FocusNode();
      await tester.pumpApp(
        ds.Shortcuts(
          Focus(focusNode: focusNode, autofocus: true, child: const SizedBox()),
          bindings: {
            const SingleActivator(LogicalKeyboardKey.escape): (
              const SizedBox(),
              () => KeyEventResult.handled,
            ),
          },
        ),
      );
      await tester.pump();
      await tester.pumpAndSettle();

      expect(focusNode.hasFocus, isTrue);
    });
  });

  group('Shortcuts focus-based routing', () {
    testWidgets('does not fire when focus is outside the subtree', (
      tester,
    ) async {
      var called = false;
      final outsideFocus = FocusNode();

      await tester.pumpApp(
        Column(
          children: [
            Focus(focusNode: outsideFocus, child: const SizedBox()),
            ds.Shortcuts(
              Focus(child: const SizedBox()),
              bindings: {
                const SingleActivator(LogicalKeyboardKey.escape): (
                  const SizedBox(),
                  () {
                    called = true;
                    return KeyEventResult.handled;
                  },
                ),
              },
            ),
          ],
        ),
      );
      outsideFocus.requestFocus();
      await tester.pump();
      await tester.pumpAndSettle();

      await tester.sendKeyEvent(LogicalKeyboardKey.escape);
      expect(called, isFalse);
    });

    testWidgets('fires when focus is inside the subtree', (tester) async {
      var called = false;
      final insideFocus = FocusNode();

      await tester.pumpApp(
        Column(
          children: [
            Focus(child: const SizedBox()),
            ds.Shortcuts(
              Focus(focusNode: insideFocus, child: const SizedBox()),
              bindings: {
                const SingleActivator(LogicalKeyboardKey.escape): (
                  const SizedBox(),
                  () {
                    called = true;
                    return KeyEventResult.handled;
                  },
                ),
              },
            ),
          ],
        ),
      );
      insideFocus.requestFocus();
      await tester.pump();
      await tester.pumpAndSettle();

      await tester.sendKeyEvent(LogicalKeyboardKey.escape);
      expect(called, isTrue);
    });

    testWidgets('stops firing after focus moves outside the subtree', (
      tester,
    ) async {
      var called = false;
      final insideFocus = FocusNode();
      final outsideFocus = FocusNode();

      await tester.pumpApp(
        Column(
          children: [
            Focus(focusNode: outsideFocus, child: const SizedBox()),
            ds.Shortcuts(
              Focus(focusNode: insideFocus, child: const SizedBox()),
              bindings: {
                const SingleActivator(LogicalKeyboardKey.escape): (
                  const SizedBox(),
                  () {
                    called = true;
                    return KeyEventResult.handled;
                  },
                ),
              },
            ),
          ],
        ),
      );

      insideFocus.requestFocus();
      await tester.pump();
      await tester.pumpAndSettle();
      await tester.sendKeyEvent(LogicalKeyboardKey.escape);
      expect(called, isTrue);

      called = false;
      outsideFocus.requestFocus();
      await tester.pump();
      await tester.pumpAndSettle();
      await tester.sendKeyEvent(LogicalKeyboardKey.escape);
      expect(called, isFalse);
    });
  });
}
