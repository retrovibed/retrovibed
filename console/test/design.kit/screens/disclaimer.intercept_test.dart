import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/modals.dart' as modals;
import 'package:retrovibed/design.kit/screens/disclaimer.intercept.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

Widget _overlay(void Function(bool) complete) {
  return Column(
    mainAxisSize: MainAxisSize.min,
    children: [
      TextButton(onPressed: () => complete(false), child: const Text('Decline')),
      TextButton(onPressed: () => complete(true), child: const Text('Accept')),
    ],
  );
}

void main() {
  group('DisclaimerIntercept', () {
    testWidgets('blocks the tap from reaching the child until the disclaimer is accepted', (tester) async {
      bool tapped = false;
      await tester.pumpApp(
        modals.Node(
          DisclaimerIntercept(
            TextButton(onPressed: () => tapped = true, child: const Text('child')),
            cacheid: 'test.disclaimer',
            cached: (_) => false,
            overlay: _overlay,
          ),
        ),
      );
      await tester.pumpAndSettle();

      // the catcher's opaque overlay sits in front of the child, so the tap
      // target is intentionally not the Text widget itself.
      await tester.tap(find.text('child'), warnIfMissed: false);
      await tester.pumpAndSettle();

      expect(tapped, isFalse);
      expect(find.text('Accept'), findsOneWidget);

      await tester.tap(find.text('Decline'));
      await tester.pumpAndSettle();

      expect(tapped, isFalse);
      expect(tester.takeException(), isNull);
    });

    testWidgets('consumes the gating tap on accept instead of forwarding it to the child', (tester) async {
      bool tapped = false;
      String? acknowledgedId;
      await tester.pumpApp(
        modals.Node(
          DisclaimerIntercept(
            TextButton(onPressed: () => tapped = true, child: const Text('child')),
            cacheid: 'test.disclaimer',
            cached: (_) => false,
            acknowledge: (id) => acknowledgedId = id,
            overlay: _overlay,
          ),
        ),
      );
      await tester.pumpAndSettle();

      // same overlay-over-child situation as above.
      await tester.tap(find.text('child'), warnIfMissed: false);
      await tester.pumpAndSettle();

      expect(tapped, isFalse);

      await tester.tap(find.text('Accept'));
      await tester.pumpAndSettle();

      expect(tapped, isFalse, reason: 'the gating tap is consumed, not replayed onto the child');
      expect(acknowledgedId, equals('test.disclaimer'));
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders the child natively with no catcher once already cached', (tester) async {
      bool tapped = false;
      await tester.pumpApp(
        modals.Node(
          DisclaimerIntercept(
            TextButton(onPressed: () => tapped = true, child: const Text('child')),
            cacheid: 'test.disclaimer',
            cached: (_) => true,
            overlay: _overlay,
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('child'));
      await tester.pumpAndSettle();

      expect(tapped, isTrue);
      expect(find.text('Accept'), findsNothing);
      expect(tester.takeException(), isNull);
    });
  });
}
