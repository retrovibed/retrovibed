import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/screens/glow.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('Glow constrained parent', () {
    testWidgets('renders within fixed SizedBox constraints', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 150,
          child: Glow(
            Container(color: Colors.blue, child: const Text('Child')),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Child'), findsOneWidget);
      expect(tester.takeException(), isNull);

      final size = tester.getSize(find.byType(Glow).last);
      expect(size.width, equals(200));
      expect(size.height, equals(150));
    });

    testWidgets('renders in Column with fixed height', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Column(
          children: [
            SizedBox(
              height: 100,
              child: Glow(
                const Text('In column'),
              ),
            ),
            const Expanded(child: SizedBox()),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('In column'), findsOneWidget);
      expect(tester.takeException(), isNull);

      final size = tester.getSize(find.byType(Glow).last);
      expect(size.height, equals(100));
    });

    testWidgets('renders in Row with fixed width', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Row(
          children: [
            SizedBox(
              width: 150,
              child: Glow(
                const Text('In row'),
              ),
            ),
            const Expanded(child: SizedBox()),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('In row'), findsOneWidget);
      expect(tester.takeException(), isNull);

      final size = tester.getSize(find.byType(Glow).last);
      expect(size.width, equals(150));
    });

    testWidgets('renders with small dimensions', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Center(
          child: SizedBox(
            width: 50,
            height: 50,
            child: Glow(
              Container(color: Colors.blue),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);

      final size = tester.getSize(find.byType(Glow).last);
      expect(size.width, equals(50));
      expect(size.height, equals(50));
    });

    testWidgets('renders with zero width constraint', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Center(
          child: SizedBox(
            width: 0,
            height: 100,
            child: Glow(
              Container(color: Colors.blue),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('renders with zero height constraint', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Center(
          child: SizedBox(
            width: 100,
            height: 0,
            child: Glow(
              Container(color: Colors.blue),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });
  });

  group('Glow unconstrained parent', () {
    testWidgets('renders in ListView with fixed height child', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ListView(
          children: [
            SizedBox(
              height: 200,
              child: Glow(
                const Text('In ListView'),
              ),
            ),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('In ListView'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in SingleChildScrollView with fixed height', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SingleChildScrollView(
          child: SizedBox(
            height: 300,
            child: Glow(
              const Text('In scroll'),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('In scroll'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in horizontal ListView with fixed width child', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          height: 100,
          child: ListView(
            scrollDirection: Axis.horizontal,
            children: [
              SizedBox(
                width: 200,
                child: Glow(
                  const Text('Horizontal'),
                ),
              ),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Horizontal'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
