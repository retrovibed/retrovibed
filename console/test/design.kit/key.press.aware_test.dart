import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('KeyPressAware.delete', () {
    testWidgets('calls onPress when delete key is pressed', (tester) async {
      var called = false;

      await tester.pumpApp(
        ds.KeyPressAware.delete(
          Focus(autofocus: true, child: const SizedBox()),
          onPress: () async => called = true,
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyEvent(LogicalKeyboardKey.delete);
      expect(called, isTrue);
    });

    testWidgets('calls onPress when backspace key is pressed', (tester) async {
      var called = false;

      await tester.pumpApp(
        ds.KeyPressAware.delete(
          Focus(autofocus: true, child: const SizedBox()),
          onPress: () async => called = true,
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyEvent(LogicalKeyboardKey.backspace);
      expect(called, isTrue);
    });

    testWidgets('does not call onPress for unrelated keys', (tester) async {
      var called = false;

      await tester.pumpApp(
        ds.KeyPressAware.delete(
          Focus(autofocus: true, child: const SizedBox()),
          onPress: () async => called = true,
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyEvent(LogicalKeyboardKey.enter);
      expect(called, isFalse);
    });

    testWidgets('focus is not lost after delete keypress', (tester) async {
      final focusNode = FocusNode();

      await tester.pumpApp(
        ds.KeyPressAware.delete(
          Focus(focusNode: focusNode, autofocus: true, child: const SizedBox()),
          onPress: () async {},
        ),
      );
      await tester.pumpAndSettle();
      expect(focusNode.hasFocus, isTrue);

      await tester.sendKeyEvent(LogicalKeyboardKey.delete);
      await tester.pump();
      expect(focusNode.hasFocus, isTrue);
    });
  });

  group('KeyPressAware.refresh', () {
    testWidgets('calls onPress when f5 is pressed', (tester) async {
      var called = false;

      await tester.pumpApp(
        ds.KeyPressAware.refresh(
          Focus(autofocus: true, child: const SizedBox()),
          onPress: () async => called = true,
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyEvent(LogicalKeyboardKey.f5);
      expect(called, isTrue);
    });

    testWidgets('does not call onPress for unrelated keys', (tester) async {
      var called = false;

      await tester.pumpApp(
        ds.KeyPressAware.refresh(
          Focus(autofocus: true, child: const SizedBox()),
          onPress: () async => called = true,
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyEvent(LogicalKeyboardKey.enter);
      expect(called, isFalse);
    });

    testWidgets('focus is not lost after f5 keypress', (tester) async {
      final focusNode = FocusNode();

      await tester.pumpApp(
        ds.KeyPressAware.refresh(
          Focus(focusNode: focusNode, autofocus: true, child: const SizedBox()),
          onPress: () async {},
        ),
      );
      await tester.pumpAndSettle();
      expect(focusNode.hasFocus, isTrue);

      await tester.sendKeyEvent(LogicalKeyboardKey.f5);
      await tester.pump();
      expect(focusNode.hasFocus, isTrue);
    });

    testWidgets('does not fire a second time while onPress is in flight', (
      tester,
    ) async {
      var callCount = 0;
      final completer = Completer<void>();

      await tester.pumpApp(
        ds.KeyPressAware.refresh(
          Focus(autofocus: true, child: const SizedBox()),
          onPress: () {
            callCount++;
            return completer.future;
          },
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyEvent(LogicalKeyboardKey.f5);
      await tester.sendKeyEvent(LogicalKeyboardKey.f5);
      expect(callCount, 1);

      completer.complete();
      await tester.pumpAndSettle();
    });
  });
}
