import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/buttons.loading.icon.dart';

Widget _wrap(Widget child) => MaterialApp(home: Scaffold(body: child));

void main() {
  group('LoadingIconButton toggled state', () {
    testWidgets('toggled=true applies primary color to IconButton', (
      WidgetTester tester,
    ) async {
      await tester.pumpWidget(_wrap(
        LoadingIconButton(
          onPressed: () async {},
          icon: const Icon(Icons.star),
          toggled: true,
        ),
      ));

      final context = tester.element(find.byType(LoadingIconButton));
      final primaryColor = Theme.of(context).colorScheme.primary;

      final iconButton = tester.widget<IconButton>(find.byType(IconButton));
      expect(iconButton.color, equals(primaryColor));
      expect(tester.takeException(), isNull);
    });

    testWidgets('toggled=false applies no color to IconButton', (
      WidgetTester tester,
    ) async {
      await tester.pumpWidget(_wrap(
        LoadingIconButton(
          onPressed: () async {},
          icon: const Icon(Icons.star),
          toggled: false,
        ),
      ));

      final iconButton = tester.widget<IconButton>(find.byType(IconButton));
      expect(iconButton.color, isNull);
      expect(tester.takeException(), isNull);
    });

    testWidgets('toggled omitted (null) applies no color to IconButton', (
      WidgetTester tester,
    ) async {
      await tester.pumpWidget(_wrap(
        LoadingIconButton(
          onPressed: () async {},
          icon: const Icon(Icons.star),
        ),
      ));

      final iconButton = tester.widget<IconButton>(find.byType(IconButton));
      expect(iconButton.color, isNull);
      expect(tester.takeException(), isNull);
    });

    testWidgets('toggled can change from false to true', (
      WidgetTester tester,
    ) async {
      bool toggled = false;

      await tester.pumpWidget(
        StatefulBuilder(
          builder: (context, setState) => _wrap(
            Column(
              children: [
                LoadingIconButton(
                  onPressed: () async {},
                  icon: const Icon(Icons.star),
                  toggled: toggled,
                ),
                ElevatedButton(
                  onPressed: () => setState(() => toggled = !toggled),
                  child: const Text('toggle'),
                ),
              ],
            ),
          ),
        ),
      );

      final context = tester.element(find.byType(LoadingIconButton));
      final primaryColor = Theme.of(context).colorScheme.primary;

      expect(tester.widget<IconButton>(find.byType(IconButton)).color, isNull);

      await tester.tap(find.text('toggle'));
      await tester.pump();

      expect(tester.widget<IconButton>(find.byType(IconButton)).color, equals(primaryColor));
      expect(tester.takeException(), isNull);
    });

    testWidgets('toggled button still fires onPressed', (
      WidgetTester tester,
    ) async {
      int callCount = 0;

      await tester.pumpWidget(_wrap(
        LoadingIconButton(
          onPressed: () async {
            callCount++;
          },
          icon: const Icon(Icons.star),
          toggled: true,
        ),
      ));

      await tester.tap(find.byType(IconButton));
      await tester.pumpAndSettle();

      expect(callCount, equals(1));
      expect(tester.takeException(), isNull);
    });
  });
}
