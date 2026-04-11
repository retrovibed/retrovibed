import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/meta/daemon.item.dart';
import 'package:retrovibed/meta/api.dart' as api;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  final daemon = api.Daemon(description: 'Test Library');

  group('DaemonDropdownItem constrained parent', () {
    testWidgets('renders within fixed SizedBox constraints', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: SizedBox(
            width: 200,
            height: 150,
            child: DaemonDropdownItem(daemon: daemon, onTap: () {}),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Test Library'), findsOneWidget);
      expect(tester.takeException(), isNull);

      final size = tester.getSize(find.byType(DaemonDropdownItem));
      expect(size.width, equals(200));
      expect(size.height, equals(150));
    });

    testWidgets('renders in Column with fixed height', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: Column(
            children: [
              SizedBox(
                height: 100,
                child: DaemonDropdownItem(daemon: daemon, onTap: () {}),
              ),
              Expanded(child: Container()),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Test Library'), findsOneWidget);
      expect(tester.takeException(), isNull);

      final size = tester.getSize(find.byType(DaemonDropdownItem));
      expect(size.height, equals(100));
    });

    testWidgets('renders in Row with fixed width', (WidgetTester tester) async {
      await tester.pumpApp(
        Scaffold(
          body: Row(
            children: [
              SizedBox(
                width: 150,
                child: DaemonDropdownItem(daemon: daemon, onTap: () {}),
              ),
              Expanded(child: Container()),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Test Library'), findsOneWidget);
      expect(tester.takeException(), isNull);

      final size = tester.getSize(find.byType(DaemonDropdownItem));
      expect(size.width, equals(150));
    });

    testWidgets('renders with small dimensions', (WidgetTester tester) async {
      await tester.pumpApp(
        Scaffold(
          body: Center(
            child: SizedBox(
              width: 50,
              height: 50,
              child: DaemonDropdownItem(daemon: daemon, onTap: () {}),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);

      final size = tester.getSize(find.byType(DaemonDropdownItem));
      expect(size.width, equals(50));
      expect(size.height, equals(50));
    });

    testWidgets('renders with zero width constraint', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: Center(
            child: SizedBox(
              width: 0,
              height: 100,
              child: DaemonDropdownItem(daemon: daemon, onTap: () {}),
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
        Scaffold(
          body: Center(
            child: SizedBox(
              width: 100,
              height: 0,
              child: DaemonDropdownItem(daemon: daemon, onTap: () {}),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });
  });

  group('DaemonDropdownItem unconstrained parent', () {
    testWidgets('renders in ListView with fixed height child', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: ListView(
            children: [
              SizedBox(
                height: 200,
                child: DaemonDropdownItem(daemon: daemon, onTap: () {}),
              ),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Test Library'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in SingleChildScrollView with fixed height', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: SingleChildScrollView(
            child: SizedBox(
              height: 300,
              child: DaemonDropdownItem(daemon: daemon, onTap: () {}),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Test Library'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in horizontal ListView with fixed width child', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: SizedBox(
            height: 100,
            child: ListView(
              scrollDirection: Axis.horizontal,
              children: [
                SizedBox(
                  width: 200,
                  child: DaemonDropdownItem(daemon: daemon, onTap: () {}),
                ),
              ],
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Test Library'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
