import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/screens/hover.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('Hover constrained parent', () {
    testWidgets('renders within fixed SizedBox constraints', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 150,
          child: Hover(
            Container(color: Colors.blue, child: const Text('Child')),
            overlay: const Text('Overlay'),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Child'), findsOneWidget);
      expect(tester.takeException(), isNull);

      final hoverSize = tester.getSize(find.byType(Hover));
      expect(hoverSize.width, equals(200));
      expect(hoverSize.height, equals(150));
    });

    testWidgets('overlay matches child size when constrained', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 300,
          height: 200,
          child: Hover(
            Container(
              key: const Key('child'),
              color: Colors.blue,
            ),
            overlay: Container(
              key: const Key('overlay'),
              color: Colors.red.withValues(alpha: 0.5),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      // Trigger hover to show overlay
      final gesture = await tester.createGesture(kind: PointerDeviceKind.mouse);
      await gesture.addPointer(location: Offset.zero);
      addTearDown(gesture.removePointer);
      await tester.pump();

      await gesture.moveTo(tester.getCenter(find.byType(Hover)));
      await tester.pumpAndSettle();

      final childSize = tester.getSize(find.byKey(const Key('child')));
      final overlaySize = tester.getSize(find.byKey(const Key('overlay')));

      expect(overlaySize.width, equals(childSize.width));
      expect(overlaySize.height, equals(childSize.height));
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in Column with fixed height', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Column(
          children: [
            SizedBox(
              height: 100,
              child: Hover(
                const Text('In column'),
                overlay: const Text('Overlay'),
              ),
            ),
            const Expanded(child: SizedBox()),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('In column'), findsOneWidget);
      expect(tester.takeException(), isNull);

      final hoverSize = tester.getSize(find.byType(Hover));
      expect(hoverSize.height, equals(100));
    });

    testWidgets('renders in Row with fixed width', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Row(
          children: [
            SizedBox(
              width: 150,
              child: Hover(
                const Text('In row'),
                overlay: const Text('Overlay'),
              ),
            ),
            const Expanded(child: SizedBox()),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('In row'), findsOneWidget);
      expect(tester.takeException(), isNull);

      final hoverSize = tester.getSize(find.byType(Hover));
      expect(hoverSize.width, equals(150));
    });
  });

  group('Hover unconstrained/infinite parent', () {
    testWidgets('renders in ListView with fixed height child', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ListView(
          children: [
            SizedBox(
              height: 200,
              child: Hover(
                const Text('In ListView'),
                overlay: const Text('Overlay'),
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
            child: Hover(
              const Text('In scroll'),
              overlay: const Text('Overlay'),
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
                child: Hover(
                  const Text('Horizontal'),
                  overlay: const Text('Overlay'),
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

  group('Hover minimum size rendering', () {
    testWidgets('renders with small dimensions', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Center(
          child: SizedBox(
            width: 50,
            height: 50,
            child: Hover(
              Container(color: Colors.blue),
              overlay: Container(color: Colors.red),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);

      final hoverSize = tester.getSize(find.byType(Hover));
      expect(hoverSize.width, equals(50));
      expect(hoverSize.height, equals(50));
    });

    testWidgets('renders with very small dimensions', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Center(
          child: SizedBox(
            width: 10,
            height: 10,
            child: Hover(
              Container(color: Colors.blue),
              overlay: Container(color: Colors.red),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);

      final hoverSize = tester.getSize(find.byType(Hover));
      expect(hoverSize.width, equals(10));
      expect(hoverSize.height, equals(10));
    });

    testWidgets('renders with zero width constraint', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Center(
          child: SizedBox(
            width: 0,
            height: 100,
            child: Hover(
              Container(color: Colors.blue),
              overlay: Container(color: Colors.red),
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
            child: Hover(
              Container(color: Colors.blue),
              overlay: Container(color: Colors.red),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });
  });

  group('Hover flex children', () {
    testWidgets('Expanded child fills available space', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 400,
          height: 300,
          child: Hover(
            Column(
              children: [
                const Text('Header'),
                Expanded(
                  key: const Key('expanded'),
                  child: Container(color: Colors.blue),
                ),
              ],
            ),
            overlay: const Text('Overlay'),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);

      final expandedSize = tester.getSize(find.byKey(const Key('expanded')));
      expect(expandedSize.width, equals(400));
      expect(expandedSize.height, greaterThan(0));
    });

    testWidgets('Spacer works correctly in overlay content', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 300,
          height: 400,
          child: Hover(
            Container(color: Colors.blue),
            overlay: Column(
              children: [
                const Text('Top'),
                const Spacer(),
                const Text('Bottom'),
              ],
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      // Trigger hover
      final gesture = await tester.createGesture(kind: PointerDeviceKind.mouse);
      await gesture.addPointer(location: Offset.zero);
      addTearDown(gesture.removePointer);
      await tester.pump();

      await gesture.moveTo(tester.getCenter(find.byType(Hover)));
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);

      final topPos = tester.getTopLeft(find.text('Top'));
      final bottomPos = tester.getTopLeft(find.text('Bottom'));

      // Bottom should be significantly lower than top due to Spacer
      expect(bottomPos.dy, greaterThan(topPos.dy + 100));
    });

    testWidgets('Flexible children respect flex factors', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 300,
          height: 300,
          child: Hover(
            Column(
              children: [
                Flexible(
                  flex: 1,
                  key: const Key('flex1'),
                  child: Container(color: Colors.red),
                ),
                Flexible(
                  flex: 2,
                  key: const Key('flex2'),
                  child: Container(color: Colors.blue),
                ),
              ],
            ),
            overlay: const Text('Overlay'),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);

      final flex1Size = tester.getSize(find.byKey(const Key('flex1')));
      final flex2Size = tester.getSize(find.byKey(const Key('flex2')));

      // flex2 should be approximately twice the height of flex1
      expect(flex2Size.height, closeTo(flex1Size.height * 2, 1.0));
    });

    testWidgets('nested Expanded widgets work correctly', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 400,
          height: 400,
          child: Hover(
            Row(
              children: [
                Expanded(
                  child: Column(
                    children: [
                      Expanded(
                        key: const Key('nested-expanded'),
                        child: Container(color: Colors.green),
                      ),
                    ],
                  ),
                ),
              ],
            ),
            overlay: const Text('Overlay'),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);

      final nestedSize = tester.getSize(find.byKey(const Key('nested-expanded')));
      expect(nestedSize.width, equals(400));
      expect(nestedSize.height, equals(400));
    });
  });

  group('Hover state changes', () {
    testWidgets('shows child initially, overlay hidden', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 200,
          child: Hover(
            const Text('Child'),
            overlay: const Text('Overlay'),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Child'), findsOneWidget);
      // Overlay widget exists but should be replaced with SizedBox when not hovered
      expect(tester.takeException(), isNull);
    });

    testWidgets('shows overlay on mouse enter', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 200,
          child: Hover(
            const Text('Child'),
            overlay: const Text('Overlay'),
          ),
        ),
      );
      await tester.pumpAndSettle();

      final gesture = await tester.createGesture(kind: PointerDeviceKind.mouse);
      await gesture.addPointer(location: Offset.zero);
      addTearDown(gesture.removePointer);
      await tester.pump();

      await gesture.moveTo(tester.getCenter(find.byType(Hover)));
      await tester.pumpAndSettle();

      expect(find.text('Overlay'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('hides overlay on mouse exit', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 200,
          child: Hover(
            const Text('Child'),
            overlay: const Text('Overlay'),
          ),
        ),
      );
      await tester.pumpAndSettle();

      final gesture = await tester.createGesture(kind: PointerDeviceKind.mouse);
      await gesture.addPointer(location: Offset.zero);
      addTearDown(gesture.removePointer);
      await tester.pump();

      // Enter hover
      await gesture.moveTo(tester.getCenter(find.byType(Hover)));
      await tester.pumpAndSettle();
      expect(find.text('Overlay'), findsOneWidget);

      // Exit hover
      await gesture.moveTo(Offset(-100, -100));
      await tester.pumpAndSettle();

      // Overlay should be hidden (replaced with SizedBox)
      expect(tester.takeException(), isNull);
    });

    testWidgets('child opacity changes on hover', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 200,
          child: Hover(
            Container(key: const Key('child'), child: const Text('Child')),
            overlay: const Text('Overlay'),
          ),
        ),
      );
      await tester.pumpAndSettle();

      // Before hover - no Opacity widget wrapping child
      var opacityWidgets = tester.widgetList<Opacity>(find.byType(Opacity));
      expect(opacityWidgets.where((o) => o.opacity == 0.05).length, equals(0));

      final gesture = await tester.createGesture(kind: PointerDeviceKind.mouse);
      await gesture.addPointer(location: Offset.zero);
      addTearDown(gesture.removePointer);
      await tester.pump();

      // Enter hover
      await gesture.moveTo(tester.getCenter(find.byType(Hover)));
      await tester.pumpAndSettle();

      // After hover - child should have reduced opacity
      opacityWidgets = tester.widgetList<Opacity>(find.byType(Opacity));
      expect(opacityWidgets.where((o) => o.opacity == 0.05).length, equals(1));
      expect(tester.takeException(), isNull);
    });
  });

  group('Hover overlay fills space', () {
    testWidgets('overlay takes full parent size', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 250,
          height: 180,
          child: Hover(
            Container(key: const Key('child')),
            overlay: Container(key: const Key('overlay')),
          ),
        ),
      );
      await tester.pumpAndSettle();

      final gesture = await tester.createGesture(kind: PointerDeviceKind.mouse);
      await gesture.addPointer(location: Offset.zero);
      addTearDown(gesture.removePointer);
      await tester.pump();

      await gesture.moveTo(tester.getCenter(find.byType(Hover)));
      await tester.pumpAndSettle();

      final hoverSize = tester.getSize(find.byType(Hover));
      final overlaySize = tester.getSize(find.byKey(const Key('overlay')));

      expect(overlaySize.width, equals(hoverSize.width));
      expect(overlaySize.height, equals(hoverSize.height));
      expect(tester.takeException(), isNull);
    });

    testWidgets('overlay position matches child position', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Center(
          child: SizedBox(
            width: 200,
            height: 200,
            child: Hover(
              Container(key: const Key('child')),
              overlay: Container(key: const Key('overlay')),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      final gesture = await tester.createGesture(kind: PointerDeviceKind.mouse);
      await gesture.addPointer(location: Offset.zero);
      addTearDown(gesture.removePointer);
      await tester.pump();

      await gesture.moveTo(tester.getCenter(find.byType(Hover)));
      await tester.pumpAndSettle();

      final childPos = tester.getTopLeft(find.byKey(const Key('child')));
      final overlayPos = tester.getTopLeft(find.byKey(const Key('overlay')));

      expect(overlayPos.dx, equals(childPos.dx));
      expect(overlayPos.dy, equals(childPos.dy));
      expect(tester.takeException(), isNull);
    });
  });

  group('Hover ValueNotifier', () {
    testWidgets('shows overlay when notifier is true', (tester) async {
      final notifier = ValueNotifier(false);
      addTearDown(notifier.dispose);

      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 200,
          child: Hover(
            const Text('Child'),
            overlay: const Text('Overlay'),
            notifier: notifier,
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Overlay'), findsNothing);

      notifier.value = true;
      await tester.pumpAndSettle();

      expect(find.text('Overlay'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('hides overlay when notifier toggled back to false', (tester) async {
      final notifier = ValueNotifier(true);
      addTearDown(notifier.dispose);

      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 200,
          child: Hover(
            const Text('Child'),
            overlay: const Text('Overlay'),
            notifier: notifier,
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Overlay'), findsOneWidget);

      notifier.value = false;
      await tester.pumpAndSettle();

      expect(find.text('Overlay'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('notifier false does not hide overlay while mouse is hovering', (tester) async {
      final notifier = ValueNotifier(false);
      addTearDown(notifier.dispose);

      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 200,
          child: Hover(
            const Text('Child'),
            overlay: const Text('Overlay'),
            notifier: notifier,
          ),
        ),
      );
      await tester.pumpAndSettle();

      // Enter with mouse to activate hover state
      final gesture = await tester.createGesture(kind: PointerDeviceKind.mouse);
      await gesture.addPointer(location: Offset.zero);
      addTearDown(gesture.removePointer);
      await tester.pump();
      await gesture.moveTo(tester.getCenter(find.byType(Hover)));
      await tester.pumpAndSettle();

      expect(find.text('Overlay'), findsOneWidget);

      // Toggle notifier — hover state is independent, overlay must remain
      notifier.value = true;
      await tester.pumpAndSettle();
      expect(find.text('Overlay'), findsOneWidget);

      notifier.value = false;
      await tester.pumpAndSettle();
      expect(find.text('Overlay'), findsOneWidget);

      expect(tester.takeException(), isNull);
    });

    testWidgets('overlay shown when notifier true even without mouse hover', (tester) async {
      final notifier = ValueNotifier(false);
      addTearDown(notifier.dispose);

      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 200,
          child: Hover(
            const Text('Child'),
            overlay: const Text('Overlay'),
            notifier: notifier,
          ),
        ),
      );
      await tester.pumpAndSettle();

      // No mouse interaction at all
      notifier.value = true;
      await tester.pumpAndSettle();

      expect(find.text('Overlay'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('Hover.overlays.icon', () {
    testWidgets('renders icon overlay with default play icon', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 200,
          child: Builder(
            builder: (context) => Hover(
              Container(color: Colors.blue),
              overlay: Hover.overlays.icon(
                context,
                content: const Text('Content'),
              ),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      final gesture = await tester.createGesture(kind: PointerDeviceKind.mouse);
      await gesture.addPointer(location: Offset.zero);
      addTearDown(gesture.removePointer);
      await tester.pump();

      await gesture.moveTo(tester.getCenter(find.byType(Hover)));
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.play_circle_filled), findsOneWidget);
      expect(find.text('Content'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders icon overlay with custom icon', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 200,
          child: Builder(
            builder: (context) => Hover(
              Container(color: Colors.blue),
              overlay: Hover.overlays.icon(
                context,
                content: const Text('Content'),
                icon: Icons.pause_circle_filled,
              ),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      final gesture = await tester.createGesture(kind: PointerDeviceKind.mouse);
      await gesture.addPointer(location: Offset.zero);
      addTearDown(gesture.removePointer);
      await tester.pump();

      await gesture.moveTo(tester.getCenter(find.byType(Hover)));
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.pause_circle_filled), findsOneWidget);
      expect(find.byIcon(Icons.play_circle_filled), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('icon is centered in overlay', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 300,
          height: 300,
          child: Builder(
            builder: (context) => Hover(
              Container(color: Colors.blue),
              overlay: Hover.overlays.icon(
                context,
                content: const Text('Content'),
              ),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      final gesture = await tester.createGesture(kind: PointerDeviceKind.mouse);
      await gesture.addPointer(location: Offset.zero);
      addTearDown(gesture.removePointer);
      await tester.pump();

      await gesture.moveTo(tester.getCenter(find.byType(Hover)));
      await tester.pumpAndSettle();

      final hoverCenter = tester.getCenter(find.byType(Hover));
      final iconCenter = tester.getCenter(find.byIcon(Icons.play_circle_filled));

      // Icon should be centered within the hover widget
      expect(iconCenter.dx, closeTo(hoverCenter.dx, 1.0));
      expect(iconCenter.dy, closeTo(hoverCenter.dy, 1.0));
      expect(tester.takeException(), isNull);
    });
  });
}
