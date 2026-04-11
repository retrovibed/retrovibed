import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/container.dart' as ds;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('Container shrinking behavior', () {
    testWidgets('Container shrinks to parent size when child is larger', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 100,
          height: 100,
          child: ds.Container(
            SizedBox(width: 200, height: 200, child: const Text('Large child')),
            padding: EdgeInsets.zero,
            margin: EdgeInsets.zero,
          ),
        ),
      );
      await tester.pumpAndSettle();

      final containerSize = tester.getSize(find.byType(ds.Container));
      expect(containerSize.width, equals(100));
      expect(containerSize.height, equals(100));
      expect(tester.takeException(), isNull);
    });

    testWidgets('Container in Column shrinks to available height', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Column(
          children: [
            SizedBox(
              height: 50,
              child: ds.Container(
                const SizedBox(height: 100, child: Text('Tall child')),
                padding: EdgeInsets.zero,
                margin: EdgeInsets.zero,
              ),
            ),
            const Expanded(child: SizedBox()),
          ],
        ),
      );
      await tester.pumpAndSettle();

      final containerSize = tester.getSize(find.byType(ds.Container));
      expect(containerSize.height, equals(50));
      expect(tester.takeException(), isNull);
    });

    testWidgets('Container in Row shrinks to available width', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Row(
          children: [
            SizedBox(
              width: 50,
              child: ds.Container(
                const SizedBox(width: 100, child: Text('Wide')),
                padding: EdgeInsets.zero,
                margin: EdgeInsets.zero,
              ),
            ),
            const Expanded(child: SizedBox()),
          ],
        ),
      );
      await tester.pumpAndSettle();

      final containerSize = tester.getSize(find.byType(ds.Container));
      expect(containerSize.width, equals(50));
      expect(tester.takeException(), isNull);
    });

    testWidgets('Container with Column child overflows when child exceeds constraint', (
      WidgetTester tester,
    ) async {
      // Column with mainAxisSize.min cannot shrink below its children's combined height.
      // This is a Flutter layout constraint that cannot be resolved by Container.
      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 100,
          child: ds.Container(
            const Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                SizedBox(height: 50, child: Text('Row 1')),
                SizedBox(height: 50, child: Text('Row 2')),
                SizedBox(height: 50, child: Text('Row 3')),
              ],
            ),
            padding: EdgeInsets.zero,
            margin: EdgeInsets.zero,
          ),
        ),
      );
      await tester.pumpAndSettle();

      final containerSize = tester.getSize(find.byType(ds.Container));
      expect(containerSize.width, equals(200));
      expect(containerSize.height, equals(100));

      final exception = tester.takeException();
      expect(exception, isA<FlutterError>());
      expect(exception.toString(), contains('overflowed'));
    });

    testWidgets('Container in GridView cell respects aspect ratio', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 400,
          height: 400,
          child: GridView.builder(
            padding: EdgeInsets.zero,
            gridDelegate: SliverGridDelegateWithMaxCrossAxisExtent(
              maxCrossAxisExtent: 200,
              childAspectRatio: 1,
            ),
            itemCount: 1,
            itemBuilder: (context, index) {
              return ds.Container(
                const SizedBox(height: 300, child: Text('Grid item')),
                padding: EdgeInsets.zero,
                margin: EdgeInsets.zero,
              );
            },
          ),
        ),
      );
      await tester.pumpAndSettle();

      final containerSize = tester.getSize(find.byType(ds.Container));
      expect(containerSize.width, equals(200));
      expect(containerSize.height, equals(200));
      expect(tester.takeException(), isNull);
    });
  });

  group('Container expansion behavior', () {
    testWidgets('Container expands to fill parent tight constraints when child is smaller', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 300,
          height: 200,
          child: ds.Container(
            const SizedBox(width: 10, height: 10, child: Text('small')),
            padding: EdgeInsets.zero,
            margin: EdgeInsets.zero,
          ),
        ),
      );
      await tester.pumpAndSettle();

      final containerSize = tester.getSize(find.byType(ds.Container));
      expect(containerSize.width, equals(300));
      expect(containerSize.height, equals(200));
      expect(tester.takeException(), isNull);
    });

    testWidgets('Container inside Expanded fills remaining Row space', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 400,
          height: 100,
          child: Row(
            children: [
              const SizedBox(width: 100),
              Expanded(
                child: ds.Container(
                  const SizedBox(width: 10, child: Text('small')),
                  padding: EdgeInsets.zero,
                  margin: EdgeInsets.zero,
                ),
              ),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      final containerSize = tester.getSize(find.byType(ds.Container));
      expect(containerSize.width, equals(300));
      expect(tester.takeException(), isNull);
    });

    testWidgets('Container inside Expanded fills remaining Column space', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 100,
          height: 400,
          child: Column(
            children: [
              const SizedBox(height: 100),
              Expanded(
                child: ds.Container(
                  const SizedBox(height: 10, child: Text('small')),
                  padding: EdgeInsets.zero,
                  margin: EdgeInsets.zero,
                ),
              ),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      final containerSize = tester.getSize(find.byType(ds.Container));
      expect(containerSize.height, equals(300));
      expect(tester.takeException(), isNull);
    });
  });

  group('Container sizing with padding and margin', () {
    testWidgets('Container size includes padding within constraint', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 100,
          height: 100,
          child: ds.Container(
            const Text('Padded'),
            padding: EdgeInsets.all(10),
            margin: EdgeInsets.zero,
          ),
        ),
      );
      await tester.pumpAndSettle();

      final containerSize = tester.getSize(find.byType(ds.Container));
      expect(containerSize.width, equals(100));
      expect(containerSize.height, equals(100));
      expect(tester.takeException(), isNull);
    });

    testWidgets('Container size respects margin within constraint', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 100,
          height: 100,
          child: ds.Container(
            const Text('Margined'),
            padding: EdgeInsets.zero,
            margin: EdgeInsets.all(10),
          ),
        ),
      );
      await tester.pumpAndSettle();

      // ds.Container fills parent, margin is applied to inner m.Container
      final containerSize = tester.getSize(find.byType(ds.Container));
      expect(containerSize.width, equals(100));
      expect(containerSize.height, equals(100));
      expect(tester.takeException(), isNull);
    });
  });
}
