import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('unconstrained parent', () {
    testWidgets('renders without overflow in unbounded horizontal space', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SingleChildScrollView(
          scrollDirection: Axis.horizontal,
          child: ds.Copyable(
            const Text('hello world'),
            onPressed: ds.Copyable.copy('hello world'),
            mainAxisSize: MainAxisSize.min,
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('hello world'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without overflow in unbounded vertical space', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: ListView(
            children: [
              ds.Copyable(
                const Text('hello world'),
                onPressed: ds.Copyable.copy('hello world'),
              ),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('hello world'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('constrained parent', () {
    testWidgets('renders without overflow in narrow SizedBox', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: SizedBox(
            width: 80,
            height: 50,
            child: ds.Copyable(
              const Text('hello world', overflow: TextOverflow.ellipsis),
              onPressed: ds.Copyable.copy('hello world'),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without overflow in Row with competing siblings', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: Row(
            children: [
              const SizedBox(width: 300),
              Flexible(
                child: ds.Copyable(
                  const Text('hello world', overflow: TextOverflow.ellipsis),
                  onPressed: ds.Copyable.copy('hello world'),
                ),
              ),
              const SizedBox(width: 300),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without overflow when Expanded in a Row', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: Row(
            children: [
              Expanded(
                child: ds.Copyable(
                  const Text('hello world', overflow: TextOverflow.ellipsis),
                  onPressed: ds.Copyable.copy('hello world'),
                ),
              ),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('copy button is present alongside content', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: SizedBox(
            width: 300,
            child: ds.Copyable(
              const Text('hello world'),
              onPressed: ds.Copyable.copy('hello world'),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('hello world'), findsOneWidget);
      expect(find.byIcon(Icons.copy), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('copy button interaction', () {
    testWidgets('onPressed is invoked when copy button is tapped', (
      WidgetTester tester,
    ) async {
      var pressed = false;
      await tester.pumpApp(
        Scaffold(
          body: ds.Copyable(
            const Text('hello world'),
            onPressed: () {
              pressed = true;
            },
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.copy));
      await tester.pump();

      expect(pressed, isTrue);
      expect(tester.takeException(), isNull);
    });

    testWidgets('null onPressed renders copy button as disabled', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        const Scaffold(body: ds.Copyable(Text('hello world'))),
      );
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.copy), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
