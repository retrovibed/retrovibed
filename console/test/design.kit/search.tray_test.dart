import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/modals.dart' as modals;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

Widget buildSearchTray({
  Key? key,
  Widget help = ds.HelpScope.None,
  Future<void> Function(String)? onSubmitted,
  void Function(fixnum.Int64)? next,
  fixnum.Int64? current,
  bool empty = true,
  bool ensureVisible = false,
  List<Widget> leading = const [],
  List<Widget> trailing = const [],
  Widget? tuning,
  EdgeInsets? padding,
}) {
  return ds.SearchTray(
    key: key,
    onSubmitted: onSubmitted ?? (_) async {},
    next: next ?? (_) {},
    current: current ?? fixnum.Int64.ZERO,
    empty: empty,
    autoscroll: ensureVisible,
    leading: leading,
    trailing: trailing,
    tuning: tuning ?? ds.Empty,
    padding: padding,
    help: help,
  );
}

void main() {
  group('SearchTray finite constraints', () {
    testWidgets('renders in finite container', (WidgetTester tester) async {
      await tester.pumpApp(
        Container(width: 400, height: 100, child: buildSearchTray()),
      );
      await tester.pumpAndSettle();

      expect(find.byType(TextField), findsOneWidget);
      expect(find.byType(ds.SearchTray), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in SizedBox with fixed dimensions', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(width: 500, height: 80, child: buildSearchTray()),
      );
      await tester.pumpAndSettle();

      expect(find.byType(TextField), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in Column with constrained parent', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 400,
          height: 200,
          child: Column(children: [buildSearchTray(), Text('Other content')]),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(TextField), findsOneWidget);
      expect(find.text('Other content'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in Card with finite constraints', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 400,
          height: 150,
          child: Card(child: buildSearchTray()),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(TextField), findsOneWidget);
      expect(find.byType(Card), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('SearchTray infinite constraints', () {
    testWidgets('renders in ListView (infinite vertical)', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ListView(children: [buildSearchTray(), Text('Item 1'), Text('Item 2')]),
      );
      await tester.pumpAndSettle();

      expect(find.byType(TextField), findsOneWidget);
      expect(find.text('Item 1'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in SingleChildScrollView with Column', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SingleChildScrollView(
          child: Column(
            children: [buildSearchTray(), Text('Scrollable content')],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(TextField), findsOneWidget);
      expect(find.text('Scrollable content'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in CustomScrollView with SliverList', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        CustomScrollView(
          slivers: [
            SliverToBoxAdapter(child: buildSearchTray()),
            SliverList(
              delegate: SliverChildListDelegate([
                Text('Sliver item 1'),
                Text('Sliver item 2'),
              ]),
            ),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(TextField), findsOneWidget);
      expect(find.text('Sliver item 1'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in nested ScrollViews', (WidgetTester tester) async {
      await tester.pumpApp(
        SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(
                height: 200,
                child: ListView(children: [buildSearchTray()]),
              ),
              Text('Outside nested scroll'),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(TextField), findsOneWidget);
      expect(find.text('Outside nested scroll'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('SearchTray in Flex layouts', () {
    testWidgets('renders in Row with Expanded', (WidgetTester tester) async {
      await tester.pumpApp(
        Row(
          children: [
            Expanded(child: buildSearchTray()),
            SizedBox(width: 50, child: Text('Side')),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(TextField), findsOneWidget);
      expect(find.text('Side'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in Column with Expanded', (WidgetTester tester) async {
      await tester.pumpApp(
        Column(
          children: [
            buildSearchTray(),
            Expanded(child: Container(color: Colors.grey)),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(TextField), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in Flex with direction vertical', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Flex(
          direction: Axis.vertical,
          children: [
            buildSearchTray(),
            Expanded(child: Text('Flex content')),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(TextField), findsOneWidget);
      expect(find.text('Flex content'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders multiple SearchTrays in Column', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Column(
          children: [
            buildSearchTray(),
            SizedBox(height: 10),
            buildSearchTray(),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(TextField), findsNWidgets(2));
      expect(find.byType(ds.SearchTray), findsNWidgets(2));
      expect(tester.takeException(), isNull);
    });
  });

  group('SearchTray with leading and trailing widgets', () {
    testWidgets('renders with leading widget in constrained container', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 500,
          height: 100,
          child: buildSearchTray(leading: [Icon(Icons.search)]),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.search), findsOneWidget);
      expect(find.byType(TextField), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders with trailing widget in constrained container', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 500,
          height: 100,
          child: buildSearchTray(trailing: [Icon(Icons.filter_list)]),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.filter_list), findsOneWidget);
      expect(find.byType(TextField), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders with both leading and trailing in ListView', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ListView(
          children: [
            buildSearchTray(
              leading: [Icon(Icons.search)],
              trailing: [Icon(Icons.close)],
            ),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.search), findsOneWidget);
      expect(find.byIcon(Icons.close), findsOneWidget);
      expect(find.byType(TextField), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders with SearchFilters as leading', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 600,
          height: 100,
          child: buildSearchTray(
            leading: [
              Chip(label: Text('Filter 1')),
              Chip(label: Text('Filter 2')),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Filter 1'), findsOneWidget);
      expect(find.text('Filter 2'), findsOneWidget);
      expect(find.byType(TextField), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('SearchTray tuning panel', () {
    testWidgets('renders with tuning widget in constrained container', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 500,
          height: 200,
          child: buildSearchTray(
            tuning: ds.buttons.settings(onPressed: () {}),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(TextField), findsOneWidget);
      expect(find.byIcon(Icons.tune), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders with tuning widget in ListView', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ListView(
          children: [
            buildSearchTray(tuning: ds.buttons.settings(onPressed: () {})),
            Text('Other content'),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(TextField), findsOneWidget);
      expect(find.byIcon(Icons.tune), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('SearchTray help', () {
    testWidgets('registers description with HelpScope', (tester) async {
      await tester.pumpApp(
        ds.HelpScope(
          buildSearchTray(
            help: ds.Hint(const Text('filter results')),
          ),
        ),
      );
      await tester.pumpAndSettle();

      final scope = tester.state<ds.HelpScopeState>(find.byType(ds.HelpScope));
      expect(scope.descriptions, hasLength(2));
      expect(tester.takeException(), isNull);
    });

    testWidgets('description appears in help overlay', (tester) async {
      await tester.pumpApp(
        modals.Node(
          ds.HelpScope(
            buildSearchTray(
              key: Key('search-tray'),
              help: ds.Hint(const Text('filter results')),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyDownEvent(LogicalKeyboardKey.altLeft);
      await tester.sendKeyDownEvent(LogicalKeyboardKey.shiftLeft);
      await tester.sendKeyEvent(LogicalKeyboardKey.slash);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.shiftLeft);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.altLeft);
      await tester.pumpAndSettle();
      await tester.tap(find.byIcon(Icons.clear));
      await tester.pumpAndSettle();
      await tester.tap(find.descendant(of: find.byKey(Key('search-tray')), matching: find.byType(InkWell)).first);
      await tester.pumpAndSettle();

      expect(find.text('filter results'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('defaults to HelpScope.None and does not register', (
      tester,
    ) async {
      await tester.pumpApp(ds.HelpScope(buildSearchTray()));
      await tester.pumpAndSettle();

      final scope = tester.state<ds.HelpScopeState>(find.byType(ds.HelpScope));
      expect(scope.descriptions, hasLength(1));
      expect(tester.takeException(), isNull);
    });
  });

  group('SearchTray ensureVisible', () {
    testWidgets('defaults to false and does not scroll', (
      WidgetTester tester,
    ) async {
      final scrollController = ScrollController();
      await tester.pumpApp(
        ListView(
          controller: scrollController,
          children: [
            SizedBox(height: 800),
            buildSearchTray(),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(scrollController.offset, equals(0.0));
      expect(tester.takeException(), isNull);
      scrollController.dispose();
    });

    testWidgets('scrolls tray into view when ensureVisible is true', (
      WidgetTester tester,
    ) async {
      final scrollController = ScrollController();
      await tester.pumpApp(
        ListView(
          controller: scrollController,
          children: [
            SizedBox(height: 800),
            buildSearchTray(ensureVisible: true),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(scrollController.offset, greaterThan(0.0));
      expect(find.byType(TextField), findsOneWidget);
      expect(tester.takeException(), isNull);
      scrollController.dispose();
    });
  });

  group('SearchTray custom padding', () {
    testWidgets('renders with zero padding in constrained container', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 400,
          height: 100,
          child: buildSearchTray(padding: EdgeInsets.zero),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(TextField), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders with custom padding in ListView', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ListView(children: [buildSearchTray(padding: EdgeInsets.all(20))]),
      );
      await tester.pumpAndSettle();

      expect(find.byType(TextField), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
