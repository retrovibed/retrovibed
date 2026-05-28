import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/buttons.loading.icon.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

// Pumps the widget into loading state using a Completer so the async
// operation stays pending. Returns the completer so callers can settle it.
Future<Completer<void>> _pumpLoading(
  WidgetTester tester, {
  required double? value,
}) async {
  final completer = Completer<void>();
  await tester.pumpApp(
    LoadingIconButton(
      onPressed: () => completer.future,
      icon: const Icon(Icons.star),
      value: value,
    ),
  );
  await tester.pumpAndSettle();
  await tester.tap(find.byType(IconButton));
  await tester.pump(); // trigger setState(_isLoading = true)
  return completer;
}

void main() {
  group('LoadingIconButton value edge cases', () {
    // Non-finite inputs are sanitized to null so CircularProgressIndicator
    // shows an indeterminate spinner instead of throwing.
    for (final (label, input) in [
      ('null', null),
      ('NaN', double.nan),
      ('infinity', double.infinity),
      ('negativeInfinity', double.negativeInfinity),
    ]) {
      testWidgets('value=$label shows indeterminate spinner without error', (tester) async {
        final completer = await _pumpLoading(tester, value: input);

        expect(find.byType(CircularProgressIndicator), findsOneWidget);
        final indicator = tester.widget<CircularProgressIndicator>(
          find.byType(CircularProgressIndicator),
        );
        expect(indicator.value, isNull);
        expect(tester.takeException(), isNull);

        completer.complete();
        await tester.pumpAndSettle();
        expect(tester.takeException(), isNull);
      });
    }

    // Finite values are passed through, clamped to [0, 1].
    for (final (label, input, expected) in [
      ('0.0', 0.0, 0.0),
      ('0.5', 0.5, 0.5),
      ('1.0', 1.0, 1.0),
      ('-0.5 clamped to 0.0', -0.5, 0.0),
      ('1.5 clamped to 1.0', 1.5, 1.0),
    ]) {
      testWidgets('value=$label shows determinate progress', (tester) async {
        final completer = await _pumpLoading(tester, value: input);

        expect(find.byType(CircularProgressIndicator), findsOneWidget);
        final indicator = tester.widget<CircularProgressIndicator>(
          find.byType(CircularProgressIndicator),
        );
        expect(indicator.value, equals(expected));
        expect(tester.takeException(), isNull);

        completer.complete();
        await tester.pumpAndSettle();
        expect(tester.takeException(), isNull);
      });
    }

    testWidgets('loading spinner replaced by icon after completion', (tester) async {
      final completer = await _pumpLoading(tester, value: null);
      expect(find.byType(CircularProgressIndicator), findsOneWidget);

      completer.complete();
      await tester.pumpAndSettle();

      expect(find.byType(CircularProgressIndicator), findsNothing);
      expect(find.byType(Icon), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

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
