import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

Widget buildMenu({
  List<Widget> leading = const [],
  List<Widget> trailing = const [],
  Widget child = const Text('child'),
}) {
  return ds.CompactingMenu(
    child,
    leading: leading,
    trailing: trailing,
  );
}

void main() {
  group('CompactingMenu non-compact', () {
    testWidgets('renders child, leading, and trailing inline', (tester) async {
      await tester.pumpApp(
        buildMenu(
          leading: [const Text('lead')],
          trailing: [const Text('trail')],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('child'), findsOneWidget);
      expect(find.text('lead'), findsOneWidget);
      expect(find.text('trail'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('does not show toggle button', (tester) async {
      await tester.pumpApp(buildMenu());
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.more_vert), findsNothing);
      expect(tester.takeException(), isNull);
    });
  });

  group('CompactingMenu compact', () {
    testWidgets('shows toggle button and child', (tester) async {
      await tester.pumpApp(
        buildMenu(
          leading: [const Text('lead')],
          trailing: [const Text('trail')],
        ),
        physicalSize: const Size(300, 600),
      );
      await tester.pumpAndSettle();

      expect(find.text('child'), findsOneWidget);
      expect(find.byIcon(Icons.more_vert), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('hides leading and trailing until toggled', (tester) async {
      await tester.pumpApp(
        buildMenu(
          leading: [const Text('lead')],
          trailing: [const Text('trail')],
        ),
        physicalSize: const Size(300, 600),
      );
      await tester.pumpAndSettle();

      expect(find.text('lead'), findsNothing);
      expect(find.text('trail'), findsNothing);
    });

    testWidgets('shows leading and trailing after tapping toggle', (tester) async {
      await tester.pumpApp(
        buildMenu(
          leading: [const Text('lead')],
          trailing: [const Text('trail')],
        ),
        physicalSize: const Size(300, 600),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.more_vert));
      await tester.pumpAndSettle();

      expect(find.text('lead'), findsOneWidget);
      expect(find.text('trail'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('hides leading and trailing again after second tap', (tester) async {
      await tester.pumpApp(
        buildMenu(
          leading: [const Text('lead')],
          trailing: [const Text('trail')],
        ),
        physicalSize: const Size(300, 600),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.more_vert));
      await tester.pumpAndSettle();
      await tester.tap(find.byIcon(Icons.more_vert));
      await tester.pumpAndSettle();

      expect(find.text('lead'), findsNothing);
      expect(find.text('trail'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('no toggle button when leading and trailing are empty', (tester) async {
      await tester.pumpApp(
        buildMenu(),
        physicalSize: const Size(300, 600),
      );
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.more_vert), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('pinned leading stays in row without toggling', (tester) async {
      await tester.pumpApp(
        buildMenu(
          leading: [
            ds.CompactingMenu.pinned(const Text('pinned')),
            const Text('menu item'),
          ],
        ),
        physicalSize: const Size(300, 600),
      );
      await tester.pumpAndSettle();

      expect(find.text('pinned'), findsOneWidget);
      expect(find.text('menu item'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('pinned trailing stays in row without toggling', (tester) async {
      await tester.pumpApp(
        buildMenu(
          trailing: [
            const Text('menu item'),
            ds.CompactingMenu.pinned(const Text('pinned')),
          ],
        ),
        physicalSize: const Size(300, 600),
      );
      await tester.pumpAndSettle();

      expect(find.text('pinned'), findsOneWidget);
      expect(find.text('menu item'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('non-pinned items still appear in menu alongside pinned', (tester) async {
      await tester.pumpApp(
        buildMenu(
          leading: [
            ds.CompactingMenu.pinned(const Text('pinned')),
            const Text('menu item'),
          ],
        ),
        physicalSize: const Size(300, 600),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.more_vert));
      await tester.pumpAndSettle();

      expect(find.text('pinned'), findsOneWidget);
      expect(find.text('menu item'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders child and all items inline when all items are pinned', (tester) async {
      await tester.pumpApp(
        buildMenu(
          leading: [ds.CompactingMenu.pinned(const Text('pinned lead'))],
          trailing: [ds.CompactingMenu.pinned(const Text('pinned trail'))],
        ),
        physicalSize: const Size(300, 600),
      );
      await tester.pumpAndSettle();

      // menuItems is empty when all items are pinned, so the full inline row is rendered
      expect(find.text('child'), findsOneWidget);
      expect(find.text('pinned lead'), findsOneWidget);
      expect(find.text('pinned trail'), findsOneWidget);
      expect(find.byIcon(Icons.more_vert), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('scrolls into view on toggle when menu extends below viewport', (tester) async {
      final controller = ScrollController();
      addTearDown(controller.dispose);

      await tester.pumpApp(
        SingleChildScrollView(
          controller: controller,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              const SizedBox(height: 550),
              buildMenu(
                leading: [const Text('lead')],
                trailing: [const Text('trail')],
              ),
              const SizedBox(height: 200),
            ],
          ),
        ),
        physicalSize: const Size(300, 600),
      );
      await tester.pumpAndSettle();

      expect(controller.offset, 0.0);

      await tester.tap(find.byIcon(Icons.more_vert));
      await tester.pumpAndSettle();

      // ensureVisible should have scrolled to reveal the expanded row
      expect(controller.offset, greaterThan(0.0));
      expect(tester.takeException(), isNull);
    });

    testWidgets('no toggle button when all items are pinned', (tester) async {
      await tester.pumpApp(
        buildMenu(
          leading: [ds.CompactingMenu.pinned(const Text('pinned lead'))],
          trailing: [ds.CompactingMenu.pinned(const Text('pinned trail'))],
        ),
        physicalSize: const Size(300, 600),
      );
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.more_vert), findsNothing);
      expect(find.text('pinned lead'), findsOneWidget);
      expect(find.text('pinned trail'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
