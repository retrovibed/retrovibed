import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/buttons.loading.icon.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('LoadingIconButton toggled state', () {
    testWidgets('toggled=true applies primary color to IconButton', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        LoadingIconButton(
          onPressed: () async {},
          icon: const Icon(Icons.star),
          toggled: true,
        ),
      );
      await tester.pumpAndSettle();

      final context = tester.element(find.byType(LoadingIconButton));
      final primaryColor = Theme.of(context).colorScheme.primary;

      final iconButton = tester.widget<IconButton>(find.byType(IconButton));
      expect(iconButton.color, equals(primaryColor));
      expect(tester.takeException(), isNull);
    });

    testWidgets('toggled=false applies no color to IconButton', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        LoadingIconButton(
          onPressed: () async {},
          icon: const Icon(Icons.star),
          toggled: false,
        ),
      );
      await tester.pumpAndSettle();

      final iconButton = tester.widget<IconButton>(find.byType(IconButton));
      expect(iconButton.color, isNull);
      expect(tester.takeException(), isNull);
    });

    testWidgets('toggled omitted (null) applies no color to IconButton', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        LoadingIconButton(
          onPressed: () async {},
          icon: const Icon(Icons.star),
        ),
      );
      await tester.pumpAndSettle();

      final iconButton = tester.widget<IconButton>(find.byType(IconButton));
      expect(iconButton.color, isNull);
      expect(tester.takeException(), isNull);
    });

    testWidgets('toggled can change from false to true', (
      WidgetTester tester,
    ) async {
      bool toggled = false;

      await tester.pumpApp(
        StatefulBuilder(
          builder: (context, setState) => Column(
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
      );
      await tester.pumpAndSettle();

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

      await tester.pumpApp(
        LoadingIconButton(
          onPressed: () async {
            callCount++;
          },
          icon: const Icon(Icons.star),
          toggled: true,
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byType(IconButton));
      await tester.pumpAndSettle();

      expect(callCount, equals(1));
      expect(tester.takeException(), isNull);
    });
  });
}
