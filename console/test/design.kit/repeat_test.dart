import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/repeat.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('Repeat', () {
    testWidgets(
      'does not throw when placed inside a Column with unbounded height',
      (tester) async {
        await tester.pumpApp(
          Column(
            children: [
              Repeat(
                () => const SizedBox(width: 200, height: 200, child: Text('X')),
              ),
            ],
          ),
        );
        await tester.pumpAndSettle();
        expect(find.text('X'), findsWidgets);
        expect(tester.takeException(), isNull);
      },
    );

    testWidgets('renders single child when child fills viewport', (
      tester,
    ) async {
      // Child (1200x1200) clamped by measurement constraint (795x595) → 1 tile
      await tester.pumpApp(
        Repeat(
          () => const SizedBox(width: 1200, height: 1200, child: Text('Large')),
        ),
      );
      await tester.pumpAndSettle();
      expect(find.text('Large'), findsExactly(1));
      expect(tester.takeException(), isNull);
    });

    testWidgets('repeats child in grid based on viewport size', (tester) async {
      // With 200x200 child, spacing=5, and 800x600 viewport:
      // width: ceil(800/205) = ceil(3.9) = 4 clones
      // height: floor(600/205) = floor(2.9) = 2 clones
      // total: 4 * 2 = 8 clones
      await tester.pumpApp(
        Repeat(() => const SizedBox(width: 200, height: 200, child: Text('X'))),
      );

      await tester.pumpAndSettle();
      expect(find.text('X'), findsExactly(8));
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders child when constraints are bounded', (tester) async {
      // 100x100 child, spacing=5, in 800x600 viewport:
      // ceil(800/105)=8, floor(600/105)=5, total=40 clones
      await tester.pumpApp(
        Repeat(
          () => const SizedBox(width: 100, height: 100, child: Text('Bounded')),
        ),
      );
      await tester.pumpAndSettle();
      expect(find.text('Bounded'), findsExactly(40));
      expect(tester.takeException(), isNull);
    });

    testWidgets('handles small child size correctly', (tester) async {
      // With 50x50 child, spacing=5, and 800x600 viewport:
      // width: ceil(800/55) = 15 clones
      // height: floor(600/55) = 10 clones
      // total: 15 * 10 = 150 clones
      await tester.pumpApp(
        Repeat(() => const SizedBox(width: 50, height: 50, child: Text('A'))),
      );

      await tester.pumpAndSettle();
      expect(find.text('A'), findsExactly(150));
      expect(tester.takeException(), isNull);
    });

    testWidgets('repeats child wrapped in ConstrainedBox', (tester) async {
      // 150x150 constrained child, spacing=5, in 800x600 viewport:
      // width: ceil(800/155) = ceil(5.16) = 6 clones
      // height: floor(600/155) = floor(3.87) = 3 clones
      // total: 6 * 3 = 18 clones
      await tester.pumpApp(
        Repeat(
          () => ConstrainedBox(
            constraints: const BoxConstraints.tightFor(width: 150, height: 150),
            child: const Text('Constrained'),
          ),
        ),
      );
      await tester.pumpAndSettle();
      expect(find.text('Constrained'), findsExactly(18));
      expect(tester.takeException(), isNull);
    });

    testWidgets('repeats child with BoxConstraints.expand height', (
      tester,
    ) async {
      // BoxConstraints.expand(height: 256): measurement constrained to 795x595
      // → detected 795x256, slot 800x261
      // width: floor(800/800)=1, height: floor(600/261)=2, total: 2
      await tester.pumpApp(
        Repeat(
          () => ConstrainedBox(
            constraints: const BoxConstraints.expand(height: 256),
            child: const Text('Expanded'),
          ),
        ),
      );
      await tester.pumpAndSettle();
      expect(find.text('Expanded'), findsExactly(2));
      expect(tester.takeException(), isNull);
    });

    testWidgets('no exceptions with different resolutions', (tester) async {
      await tester.pumpApp(
        Repeat(
          () => const SizedBox(
            width: 100,
            height: 100,
            child: Text('Resolution'),
          ),
        ),
      );

      await tester.pumpAndSettle();
      expect(find.text('Resolution'), findsWidgets);
      expect(tester.takeException(), isNull);
    });

    // The following tests surface the issue where widgets that expand to fill
    // their constraints (e.g. SizedBox.expand, or cards with SizedBox.expand
    // inside) cause _MeasureSize to detect the full available area as the
    // child size, resulting in only 1 tile instead of a tiled pattern.

    testWidgets(
      'expandable child without size wrapper fills viewport and renders once',
      (tester) async {
        // SizedBox.expand: measurement constrained to 795x595 → detected 795x595
        // slot 800x600: floor(800/800)=1, floor(600/600)=1, total: 1
        await tester.pumpApp(
          Repeat(() => const SizedBox.expand(child: Text('Expand'))),
        );
        await tester.pumpAndSettle();
        expect(find.text('Expand'), findsExactly(1));
        expect(tester.takeException(), isNull);
      },
    );

    testWidgets(
      'expandable child wrapped in fixed SizedBox tiles at fixed size',
      (tester) async {
        // SizedBox(200x200) constrains measured size to 200x200.
        // In 800x600 with spacing=5: ceil(800/205)=4, floor(600/205)=2, total: 8
        await tester.pumpApp(
          Repeat(
            () => const SizedBox(
              width: 200,
              height: 200,
              child: SizedBox.expand(child: Text('Wrapped')),
            ),
          ),
        );
        await tester.pumpAndSettle();
        expect(find.text('Wrapped'), findsExactly(8));
        expect(tester.takeException(), isNull);
      },
    );

    testWidgets(
      'Column(mainAxisSize.min) with expandable Flexible inside fixed SizedBox tiles correctly',
      (tester) async {
        // Simulates ds.Card behaviour: Column(mainAxisSize.min) containing a
        // Flexible(SizedBox.expand). The 200x200 SizedBox wrapper constrains size.
        // In 800x600 with spacing=5: ceil(800/205)=4, floor(600/205)=2, total: 8
        await tester.pumpApp(
          Repeat(
            () => SizedBox(
              width: 200,
              height: 200,
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: const [
                  Flexible(child: SizedBox.expand(child: Text('Column'))),
                ],
              ),
            ),
          ),
        );
        await tester.pumpAndSettle();
        expect(find.text('Column'), findsExactly(8));
        expect(tester.takeException(), isNull);
      },
    );

    testWidgets('AspectRatio inside fixed SizedBox tiles at fixed size', (
      tester,
    ) async {
      // Simulates KnownMediaCard: AspectRatio(2/3) with SizedBox.expand inside
      // a fixed 200x200 wrapper. The measured size must be 200x200.
      // In 800x600 with spacing=5: ceil(800/205)=4, floor(600/205)=2, total: 8
      await tester.pumpApp(
        Repeat(
          () => SizedBox(
            width: 200,
            height: 200,
            child: AspectRatio(
              aspectRatio: 2 / 3,
              child: const SizedBox.expand(child: Text('Aspect')),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();
      expect(find.text('Aspect'), findsExactly(8));
      expect(tester.takeException(), isNull);
    });

    testWidgets('AspectRatio child tiles in bounded constraints', (
      tester,
    ) async {
      // Replicates ds.Recommendations: KnownMediaCard uses AspectRatio(2/3).
      // Measurement constrained to 795x595 → height=595, width=2/3*595≈397
      // slot ≈402x600: ceil(800/402)=2, floor(600/600)=1, total: 2
      await tester.pumpApp(
        Repeat(
          () => const AspectRatio(
            aspectRatio: 2 / 3,
            child: SizedBox.expand(child: Text('card')),
          ),
        ),
      );
      await tester.pumpAndSettle();
      expect(find.text('card'), findsExactly(2));
      expect(tester.takeException(), isNull);
    });
  });
}
