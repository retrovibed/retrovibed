import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/carousel.dart';
import 'package:retrovibed/design.kit/repeat.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('CarouselRow', () {
    testWidgets('renders title widget', (tester) async {
      await tester.pumpApp(
        CarouselRow(
          title: const Text('Title'),
          items: const [
            SizedBox(width: 100, height: 150, child: Text('a')),
            SizedBox(width: 100, height: 150, child: Text('b')),
          ],
        ),
      );
      await tester.pumpAndSettle();
      expect(find.text('Title'), findsOneWidget);
    });

    testWidgets('renders horizontal scroll with items', (tester) async {
      await tester.pumpApp(
        CarouselRow(
          title: const SizedBox(),
          items: [
            SizedBox(width: 100, height: 150, child: Text('1')),
            SizedBox(width: 100, height: 150, child: Text('2')),
            SizedBox(width: 100, height: 150, child: Text('3')),
          ],
        ),
      );
      await tester.pumpAndSettle();
      expect(find.text('1'), findsOneWidget);
      expect(find.text('2'), findsOneWidget);
      expect(find.text('3'), findsOneWidget);
    });

    testWidgets('shows loading indicator when loading', (tester) async {
      await tester.pumpApp(
        CarouselRow(
          title: const SizedBox(),
          items: [const SizedBox()],
          loading: true,
        ),
      );
      await tester.pump();
      expect(find.byType(CircularProgressIndicator), findsOneWidget);
    });

    testWidgets('horizontal scroll has correct direction', (tester) async {
      await tester.pumpApp(
        CarouselRow(
          title: const SizedBox(),
          items: List.generate(10, (i) => Text('Item $i')),
        ),
      );
      await tester.pumpAndSettle();
      final scroll = tester.widget<SingleChildScrollView>(
        find.byType(SingleChildScrollView),
      );
      expect(scroll.scrollDirection, equals(Axis.horizontal));
    });

    testWidgets('renders at all resolutions', (tester) async {
      await tester.pumpApp(
        CarouselRow(
          title: const Text('Test'),
          items: [const SizedBox(width: 128, height: 192)],
        ),
      );
      await tester.pumpAndSettle();
      expect(find.text('Test'), findsOneWidget);
    });

    testWidgets('horizontal scroll with many items', (tester) async {
      await tester.pumpApp(
        CarouselRow(
          title: const Text('Title'),
          items: List.generate(20, (i) => Text('Item $i')),
        ),
      );
      await tester.pumpAndSettle();
      final scroll = tester.widget<SingleChildScrollView>(
        find.byType(SingleChildScrollView),
      );
      expect(scroll.scrollDirection, equals(Axis.horizontal));
    });

    testWidgets('mouse wheel scrolls content horizontally', (tester) async {
      await tester.pumpApp(
        CarouselRow(
          title: const SizedBox(),
          items: List.generate(
            20,
            (i) => SizedBox(width: 100, height: 100, child: Text('Item $i')),
          ),
        ),
      );
      await tester.pumpAndSettle();

      final initialX = tester.getTopLeft(find.text('Item 0')).dx;

      final center = tester.getCenter(find.byType(SingleChildScrollView));
      await tester.sendEventToBinding(
        PointerScrollEvent(position: center, scrollDelta: const Offset(0, 200)),
      );
      await tester.pump();

      expect(tester.getTopLeft(find.text('Item 0')).dx, lessThan(initialX));
    });

    testWidgets('arrow right key scrolls content forward', (tester) async {
      await tester.pumpApp(
        CarouselRow(
          title: const SizedBox(),
          items: List.generate(
            20,
            (i) => SizedBox(width: 100, height: 100, child: Text('Item $i')),
          ),
        ),
      );
      await tester.pumpAndSettle();

      final initialX = tester.getTopLeft(find.text('Item 0')).dx;

      await tester.tap(find.byType(CarouselRow));
      await tester.pump();
      await tester.sendKeyEvent(LogicalKeyboardKey.arrowRight);
      await tester.pump();

      expect(tester.getTopLeft(find.text('Item 0')).dx, lessThan(initialX));
    });

    testWidgets('arrow left key scrolls content backward', (tester) async {
      await tester.pumpApp(
        CarouselRow(
          title: const SizedBox(),
          items: List.generate(
            20,
            (i) => SizedBox(width: 100, height: 100, child: Text('Item $i')),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byType(CarouselRow));
      await tester.pump();
      await tester.sendKeyEvent(LogicalKeyboardKey.arrowRight);
      await tester.pump();

      final afterRightX = tester.getTopLeft(find.text('Item 0')).dx;

      await tester.sendKeyEvent(LogicalKeyboardKey.arrowLeft);
      await tester.pump();

      expect(tester.getTopLeft(find.text('Item 0')).dx, greaterThan(afterRightX));
    });

    testWidgets('does not overflow when given a fixed-height constraint with tall items', (tester) async {
      await tester.pumpApp(
        CarouselRow(
          constraints: const BoxConstraints.tightFor(height: 256),
          title: const Text('Title'),
          items: [
            SizedBox(width: 200, height: 974, child: Text('tall')),
          ],
        ),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    });

    testWidgets('background Repeat tiles fixed-size children when items is empty', (tester) async {
      // Outer container padding=16, height=256 → column 768×224
      // title=SizedBox()=0, spacing=10 → Expanded 768×214
      // Repeat spacing=5: ceil(768/105)=8 × floor(214/105)=2 = 16
      await tester.pumpApp(
        CarouselRow(
          constraints: const BoxConstraints.tightFor(height: 256),
          title: const SizedBox(),
          background: Repeat(
            () => const SizedBox(width: 100, height: 100, child: Text('bg')),
          ),
          items: const [],
        ),
      );
      await tester.pumpAndSettle();
      expect(find.text('bg'), findsExactly(16));
    });

    testWidgets('background Repeat tiles expanding children when items is empty', (tester) async {
      // Expanding children (e.g. cards that fill available space) must still
      // tile — Repeat needs a fixed-size wrapper to measure from.
      // Repeat spacing=5: ceil(768/155)=5 × floor(214/205)=1 = 5
      await tester.pumpApp(
        CarouselRow(
          constraints: const BoxConstraints.tightFor(height: 256),
          title: const SizedBox(),
          background: Repeat(
            () => const SizedBox(
              width: 150,
              height: 200,
              child: SizedBox.expand(child: Text('card')),
            ),
          ),
          items: const [],
        ),
      );
      await tester.pumpAndSettle();
      expect(find.text('card'), findsExactly(5));
    });

    testWidgets('buttons inside items are tappable', (tester) async {
      bool tapped = false;
      await tester.pumpApp(
        CarouselRow(
          title: const SizedBox(),
          items: [
            ElevatedButton(
              onPressed: () => tapped = true,
              child: const Text('tap me'),
            ),
          ],
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.text('tap me'));
      await tester.pump();
      expect(tapped, isTrue);
    });
  });
}
