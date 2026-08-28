import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/screens/full.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/windowx.dart';
import 'package:window_manager/window_manager.dart';

void main() {
  group('Full initial state', () {
    testWidgets('seeds chromeless false when the plugin reports not-fullscreen', (tester) async {
      late BuildContext capturedContext;
      final wm = WindowManagerX.fake(fullscreen: false);

      await tester.pumpApp(Full(Builder(builder: (context) {
        capturedContext = context;
        return const SizedBox();
      }), windowManager: wm));
      await tester.pumpAndSettle();

      expect(Full.nochrome(capturedContext), isFalse);
      expect(tester.takeException(), isNull);
    });

    // Regression: window_manager's isFullScreen() is an unguarded channel
    // read that can report true before the window is actually settled
    // (see windowx_test.dart's waitUntilReadyToShow correction test). Full
    // must not get permanently stuck on a bad initial read - a subsequent
    // onWindowLeaveFullScreen event has to be able to correct it.
    testWidgets('a stale fullscreen:true initial report is corrected by a later leave-fullscreen event', (
      tester,
    ) async {
      late BuildContext capturedContext;
      final wm = WindowManagerX.fake(fullscreen: true);

      await tester.pumpApp(Full(Builder(builder: (context) {
        capturedContext = context;
        return const SizedBox();
      }), windowManager: wm));
      await tester.pumpAndSettle();

      expect(Full.nochrome(capturedContext), isTrue);

      (Full.of(capturedContext) as WindowListener).onWindowLeaveFullScreen();
      await tester.pump();

      expect(Full.nochrome(capturedContext), isFalse);
      expect(tester.takeException(), isNull);
    });
  });

  group('Full.toggle', () {
    testWidgets('flips chromeless and calls setFullScreen with the new value', (tester) async {
      final calls = <String>[];
      late BuildContext capturedContext;

      await tester.pumpApp(Full(
        Builder(builder: (context) {
          capturedContext = context;
          return const SizedBox();
        }),
        windowManager: WindowManagerX.fake(calls: calls),
      ));
      await tester.pumpAndSettle();
      calls.clear();

      expect(Full.nochrome(capturedContext), isFalse);
      Full.toggle(capturedContext);
      await tester.pumpAndSettle();

      expect(Full.nochrome(capturedContext), isTrue);
      expect(calls, contains('setFullScreen'));
    });

    testWidgets('a rejecting setFullScreen does not throw past toggle()', (tester) async {
      late BuildContext capturedContext;
      final wm = WindowManagerX(setFullScreen: (_) => Future.error('boom'));

      await tester.pumpApp(Full(
        Builder(builder: (context) {
          capturedContext = context;
          return const SizedBox();
        }),
        windowManager: wm,
      ));
      await tester.pumpAndSettle();

      Full.toggle(capturedContext);
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });
  });

  group('Full window listener sync', () {
    testWidgets('onWindowEnterFullScreen/onWindowLeaveFullScreen update chromeless independent of toggle', (
      tester,
    ) async {
      late BuildContext capturedContext;

      await tester.pumpApp(Full(
        Builder(builder: (context) {
          capturedContext = context;
          return const SizedBox();
        }),
        windowManager: WindowManagerX.fake(),
      ));
      await tester.pumpAndSettle();

      expect(Full.nochrome(capturedContext), isFalse);

      (Full.of(capturedContext) as WindowListener).onWindowEnterFullScreen();
      await tester.pump();
      expect(Full.nochrome(capturedContext), isTrue);

      (Full.of(capturedContext) as WindowListener).onWindowLeaveFullScreen();
      await tester.pump();
      expect(Full.nochrome(capturedContext), isFalse);
    });

    testWidgets('addListener/removeListener are called on mount/dispose', (tester) async {
      final calls = <String>[];

      await tester.pumpApp(Full(const SizedBox(), windowManager: WindowManagerX.fake(calls: calls)));
      await tester.pumpAndSettle();

      expect(calls, contains('addListener'));
      expect(calls, isNot(contains('removeListener')));

      await tester.pumpWidget(const SizedBox());
      await tester.pumpAndSettle();

      expect(calls, contains('removeListener'));
    });
  });
}
