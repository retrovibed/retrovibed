import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/screens/guarded.mask.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  Finder maskIgnorePointer() => find.descendant(
        of: find.byType(GuardedMask),
        matching: find.byType(IgnorePointer),
      ).first;

  Finder maskOpacity() => find.descendant(
        of: find.byType(GuardedMask),
        matching: find.byType(Opacity),
      ).first;

  group('GuardedMask', () {
    group('unprotected', () {
      testWidgets('full opacity', (tester) async {
        await tester.pumpApp(const GuardedMask(child: SizedBox(width: 100, height: 100)));
        await tester.pumpAndSettle();

        expect(tester.widget<Opacity>(maskOpacity()).opacity, 1.0);
      });

      testWidgets('IgnorePointer is not ignoring', (tester) async {
        await tester.pumpApp(const GuardedMask(child: SizedBox(width: 100, height: 100)));
        await tester.pumpAndSettle();

        expect(tester.widget<IgnorePointer>(maskIgnorePointer()).ignoring, isFalse);
      });

      testWidgets('pointer events pass through', (tester) async {
        var tapped = false;
        await tester.pumpApp(GuardedMask(
          child: GestureDetector(
            behavior: HitTestBehavior.opaque,
            onTap: () => tapped = true,
            child: const SizedBox(width: 100, height: 100),
          ),
        ));
        await tester.pumpAndSettle();

        await tester.tap(find.byType(GestureDetector));
        expect(tapped, isTrue);
      });
    });

    group('protected', () {
      testWidgets('half opacity', (tester) async {
        await tester.pumpApp(const GuardedMask(protected: true, child: SizedBox(width: 100, height: 100)));
        await tester.pumpAndSettle();

        expect(tester.widget<Opacity>(maskOpacity()).opacity, 0.5);
      });

      testWidgets('IgnorePointer is ignoring', (tester) async {
        await tester.pumpApp(const GuardedMask(protected: true, child: SizedBox(width: 100, height: 100)));
        await tester.pumpAndSettle();

        expect(tester.widget<IgnorePointer>(maskIgnorePointer()).ignoring, isTrue);
      });

      testWidgets('pointer events are blocked', (tester) async {
        var tapped = false;
        await tester.pumpApp(GuardedMask(
          protected: true,
          child: GestureDetector(
            behavior: HitTestBehavior.opaque,
            onTap: () => tapped = true,
            child: const SizedBox(width: 100, height: 100),
          ),
        ));
        await tester.pumpAndSettle();

        await tester.tap(find.byType(GestureDetector), warnIfMissed: false);
        expect(tapped, isFalse);
      });
    });
  });
}
