import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/library/known.media.card.dart';
import 'package:retrovibed/library/api.dart' as api;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = ValueVariant({
  ...Resolutions.all.entries,
  const MapEntry('small (256x256)', Size(256, 256)),
});

Future<void> _hover(WidgetTester tester) async {
  final gesture = await tester.createGesture(kind: PointerDeviceKind.mouse);
  await gesture.addPointer(location: Offset.zero);
  addTearDown(gesture.removePointer);
  await gesture.moveTo(tester.getCenter(find.byType(KnownMediaCard)));
  await tester.pumpAndSettle();
}

void main() {
  group('KnownMediaCard in horizontal scroll (carousel-like layout)', () {
    testWidgets('cards have consistent size regardless of description length', (tester) async {
      await tester.pumpApp(
        SizedBox(
          height: 256,
          child: SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            child: Row(
              children: [
                KnownMediaCard(
                  api.Known(description: 'Short', summary: 'summary'),
                  key: const Key('card1'),
                ),
                KnownMediaCard(
                  api.Known(
                    description: 'A Very Long Title That Should Not Change Card Width At All',
                    summary: 'summary',
                  ),
                  key: const Key('card2'),
                ),
              ],
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      final size1 = tester.getSize(find.byKey(const Key('card1')));
      final size2 = tester.getSize(find.byKey(const Key('card2')));

      expect(size1, equals(size2));
      expect(tester.takeException(), isNull);
    });

    testWidgets('does not overflow when hovering a card with a long summary', (tester) async {
      await tester.pumpApp(
        SizedBox(
          height: 256,
          child: SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            child: Row(
              children: [
                KnownMediaCard(
                  api.Known(description: 'Test', summary: 'A' * 500),
                ),
              ],
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);

      await _hover(tester);
      expect(tester.takeException(), isNull);
    });
  });

  group('KnownMediaCard', () {
    testWidgets('renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        physicalSize: entry.value,
        KnownMediaCard(api.Known(description: 'Test Known', summary: 'Test summary')),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);

      await _hover(tester);
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders without overflow when highlighted', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        physicalSize: entry.value,
        KnownMediaCard(
          api.Known(description: 'Test Known', summary: 'Test summary'),
          highlighted: true,
        ),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);

      await _hover(tester);
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders without overflow with onTap', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        physicalSize: entry.value,
        KnownMediaCard(
          api.Known(description: 'Test Known', summary: 'Test summary'),
          onTap: () {},
        ),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);

      await _hover(tester);
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('onTap is invoked when tapped', (tester) async {
      var tapped = false;
      await tester.pumpApp(
        KnownMediaCard(
          api.Known(description: 'Test', summary: 'summary'),
          onTap: () { tapped = true; },
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(KnownMediaCard));
      await tester.pump();
      expect(tapped, isTrue);
    });

    testWidgets('renders without overflow with onDoubleTap', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        physicalSize: entry.value,
        KnownMediaCard(
          api.Known(description: 'Test Known', summary: 'Test summary'),
          onDoubleTap: () {},
        ),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);

      await _hover(tester);
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders without overflow with trailing widget', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        physicalSize: entry.value,
        KnownMediaCard(
          api.Known(description: 'Test Known', summary: 'Test summary'),
          trailing: const Text('trailing'),
        ),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);

      await _hover(tester);
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders without overflow with custom icon', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        physicalSize: entry.value,
        KnownMediaCard(
          api.Known(description: 'Test Known', summary: 'Test summary'),
          icon: Icons.download,
        ),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);

      await _hover(tester);
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('future factory renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        physicalSize: entry.value,
        KnownMediaCard.future(
          Future.value(api.Known(description: 'Test Known', summary: 'Test summary')),
        ),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);

      await _hover(tester);
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);
  });

  group('KnownMediaCard cursor', () {
    testWidgets('shows click cursor when onTap is set', (tester) async {
      await tester.pumpApp(
        KnownMediaCard(
          api.Known(description: 'Test', summary: 'summary'),
          onTap: () {},
        ),
      );
      await tester.pumpAndSettle();

      expect(
        tester.resolvedCursorAt(find.byType(KnownMediaCard)),
        SystemMouseCursors.click,
      );
    });

    testWidgets('shows click cursor when onDoubleTap is set', (tester) async {
      await tester.pumpApp(
        KnownMediaCard(
          api.Known(description: 'Test', summary: 'summary'),
          onDoubleTap: () {},
        ),
      );
      await tester.pumpAndSettle();

      expect(
        tester.resolvedCursorAt(find.byType(KnownMediaCard)),
        SystemMouseCursors.click,
      );
    });

    testWidgets('shows basic cursor when onDoubleTap is not set', (tester) async {
      await tester.pumpApp(
        KnownMediaCard(api.Known(description: 'Test', summary: 'summary')),
      );
      await tester.pumpAndSettle();

      expect(
        tester.resolvedCursorAt(find.byType(KnownMediaCard)),
        SystemMouseCursors.basic,
      );
    });
  });
}
