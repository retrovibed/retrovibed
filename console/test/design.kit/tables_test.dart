import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/modals.dart' as modals;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

// Compact breakpoint is 400 — below it the Table reverses vertical direction.
const double kWide = 600.0;
const double kNarrow = 300.0;

Widget buildTable({
  List<String> children = const ['alpha', 'beta', 'gamma'],
  Widget leading = const SizedBox(),
  Widget trailing = const SizedBox(),
  Widget empty = const SizedBox(),
  Widget overlay = const SizedBox(),
  Widget help = ds.HelpScope.None,
  bool loading = false,
  Widget cause = const SizedBox(),
  bool expanded = true,
}) {
  final render =
      expanded ? ds.Table.expanded<String>((item) => Text(item)) : ds.Table.inline<String>((item) => Text(item));
  return ds.Table<String>(
    render,
    children: children,
    leading: leading,
    trailing: trailing,
    empty: empty,
    overlay: overlay,
    help: help,
    loading: loading,
    cause: cause,
  );
}

void main() {
  // ---------------------------------------------------------------------------
  // Layout: constrained parents
  // ---------------------------------------------------------------------------
  group('Table constrained layout', () {
    testWidgets('renders in SizedBox with fixed dimensions', (tester) async {
      await tester.pumpApp(
        SizedBox(width: kWide, height: 400, child: buildTable()),
      );
      await tester.pumpAndSettle();

      expect(find.text('alpha'), findsOneWidget);
      expect(find.text('beta'), findsOneWidget);
      expect(find.text('gamma'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in finite Container', (tester) async {
      await tester.pumpApp(
        Container(width: kWide, height: 400, child: buildTable()),
      );
      await tester.pumpAndSettle();

      expect(find.text('alpha'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in Card with finite constraints', (tester) async {
      await tester.pumpApp(
        Card(child: buildTable()),
        physicalSize: const Size(kWide, 400),
        fit: FlexFit.tight,
      );
      await tester.pumpAndSettle();

      expect(find.byType(Card), findsOneWidget);
      expect(find.text('alpha'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders as Expanded child in Column', (tester) async {
      await tester.pumpApp(
        Column(
          children: [
            Text('header'),
            Expanded(child: buildTable()),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('header'), findsOneWidget);
      expect(find.text('alpha'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in Row with Expanded', (tester) async {
      await tester.pumpApp(
        Row(
          children: [
            Expanded(child: buildTable()),
            SizedBox(width: 50, child: Text('side')),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('side'), findsOneWidget);
      expect(find.text('alpha'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in ConstrainedBox', (tester) async {
      await tester.pumpApp(
        ConstrainedBox(
          constraints: BoxConstraints(maxWidth: kWide, maxHeight: 400),
          child: buildTable(),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('alpha'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  // ---------------------------------------------------------------------------
  // Layout: unconstrained / scrollable parents
  // ---------------------------------------------------------------------------
  group('Table unconstrained layout', () {
    testWidgets('renders inside SingleChildScrollView with Column', (tester) async {
      await tester.pumpApp(
        SingleChildScrollView(
          child: Column(
            children: [buildTable(), Text('after')],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('alpha'), findsOneWidget);
      expect(find.text('after'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders inside ListView', (tester) async {
      await tester.pumpApp(
        ListView(children: [buildTable(), Text('item')]),
      );
      await tester.pumpAndSettle();

      expect(find.text('alpha'), findsOneWidget);
      expect(find.text('item'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders inside CustomScrollView with SliverToBoxAdapter', (tester) async {
      await tester.pumpApp(
        CustomScrollView(
          slivers: [
            SliverToBoxAdapter(child: buildTable()),
            SliverToBoxAdapter(child: Text('sliver')),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('alpha'), findsOneWidget);
      expect(find.text('sliver'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('multiple Tables in Column coexist', (tester) async {
      await tester.pumpApp(
        SingleChildScrollView(
          child: Column(
            children: [
              buildTable(children: ['first-a', 'first-b']),
              buildTable(children: ['second-a', 'second-b']),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('first-a'), findsOneWidget);
      expect(find.text('second-a'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  // ---------------------------------------------------------------------------
  // Compact breakpoint (threshold: 400)
  // ---------------------------------------------------------------------------
  group('Table compact breakpoint', () {
    testWidgets('wide layout uses VerticalDirection.down', (tester) async {
      await tester.pumpApp(
        buildTable(leading: Text('search-bar')),
        physicalSize: const Size(kWide, 500),
        fit: FlexFit.tight,
      );
      await tester.pumpAndSettle();

      expect(find.text('search-bar'), findsOneWidget);

      final column = tester.widget<Column>(
        find
            .ancestor(
              of: find.text('search-bar'),
              matching: find.byType(Column),
            )
            .first,
      );
      expect(column.verticalDirection, VerticalDirection.down);
      expect(tester.takeException(), isNull);
    });

    testWidgets('narrow layout uses VerticalDirection.up', (tester) async {
      await tester.pumpApp(
        buildTable(leading: Text('search-bar')),
        physicalSize: const Size(kNarrow, 500),
        fit: FlexFit.tight,
      );
      await tester.pumpAndSettle();

      expect(find.text('search-bar'), findsOneWidget);

      final column = tester.widget<Column>(
        find
            .ancestor(
              of: find.text('search-bar'),
              matching: find.byType(Column),
            )
            .first,
      );
      expect(column.verticalDirection, VerticalDirection.up);
      expect(tester.takeException(), isNull);
    });

    testWidgets('two Tables side-by-side at wide and narrow widths', (tester) async {
      // 500 + 250 = 750, fits within the 800px test viewport width.
      await tester.pumpApp(
        Row(
          children: [
            SizedBox(
              width: 500,
              height: 300,
              child: buildTable(children: ['wide-item']),
            ),
            SizedBox(
              width: 250,
              height: 300,
              child: buildTable(children: ['narrow-item']),
            ),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('wide-item'), findsOneWidget);
      expect(find.text('narrow-item'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  // ---------------------------------------------------------------------------
  // Render modes: expanded vs inline
  // ---------------------------------------------------------------------------
  group('Table render modes', () {
    testWidgets('expanded wraps rows in SingleChildScrollView', (tester) async {
      await tester.pumpApp(
        buildTable(expanded: true),
        physicalSize: const Size(kWide, 400),
        fit: FlexFit.tight,
      );
      await tester.pumpAndSettle();

      expect(find.byType(SingleChildScrollView), findsOneWidget);
      expect(find.text('alpha'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('inline does not add SingleChildScrollView', (tester) async {
      await tester.pumpApp(
        SingleChildScrollView(
          child: Column(children: [buildTable(expanded: false)]),
        ),
      );
      await tester.pumpAndSettle();

      // Only the outer scroll view — Table.inline uses a plain Column.
      expect(find.byType(SingleChildScrollView), findsOneWidget);
      expect(find.text('alpha'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('expanded renders all items', (tester) async {
      await tester.pumpApp(
        buildTable(
          children: ['one', 'two', 'three'],
          expanded: true,
        ),
        physicalSize: const Size(kWide, 400),
        fit: FlexFit.tight,
      );
      await tester.pumpAndSettle();

      expect(find.text('one'), findsOneWidget);
      expect(find.text('two'), findsOneWidget);
      expect(find.text('three'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('inline renders all items', (tester) async {
      await tester.pumpApp(
        SingleChildScrollView(
          child: Column(
            children: [
              buildTable(
                children: ['one', 'two', 'three'],
                expanded: false,
              ),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('one'), findsOneWidget);
      expect(find.text('two'), findsOneWidget);
      expect(find.text('three'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  // ---------------------------------------------------------------------------
  // Empty state
  // ---------------------------------------------------------------------------
  group('Table empty state', () {
    testWidgets('shows empty widget when children is empty', (tester) async {
      await tester.pumpApp(
        buildTable(children: [], empty: Text('nothing here')),
        physicalSize: const Size(kWide, 400),
        fit: FlexFit.tight,
      );
      await tester.pumpAndSettle();

      expect(find.text('nothing here'), findsOneWidget);
      expect(find.text('alpha'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('hides empty widget when children is non-empty', (tester) async {
      await tester.pumpApp(
        buildTable(children: ['alpha'], empty: Text('nothing here')),
        physicalSize: const Size(kWide, 400),
        fit: FlexFit.tight,
      );
      await tester.pumpAndSettle();

      expect(find.text('alpha'), findsOneWidget);
      expect(find.text('nothing here'), findsNothing);
      expect(tester.takeException(), isNull);
    });
  });

  // ---------------------------------------------------------------------------
  // Leading and trailing widgets
  // ---------------------------------------------------------------------------
  group('Table leading and trailing', () {
    testWidgets('leading widget renders', (tester) async {
      await tester.pumpApp(
        buildTable(leading: Text('lead')),
        physicalSize: const Size(kWide, 400),
        fit: FlexFit.tight,
      );
      await tester.pumpAndSettle();

      expect(find.text('lead'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('trailing widget renders', (tester) async {
      await tester.pumpApp(
        buildTable(trailing: Text('trail')),
        physicalSize: const Size(kWide, 400),
        fit: FlexFit.tight,
      );
      await tester.pumpAndSettle();

      expect(find.text('trail'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('leading as SearchTray (common usage pattern)', (tester) async {
      await tester.pumpApp(
        buildTable(
          leading: ds.SearchTray(
            onSubmitted: (_) async {},
            next: (_) {},
            current: ds.SearchTray.Zero,
            empty: true,
          ),
        ),
        physicalSize: const Size(kWide, 400),
        fit: FlexFit.tight,
      );
      await tester.pumpAndSettle();

      expect(find.byType(ds.SearchTray), findsOneWidget);
      expect(find.text('alpha'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('leading as Column with SearchTray and header (common usage pattern)', (tester) async {
      await tester.pumpApp(
        buildTable(
          leading: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              ds.SearchTray(
                onSubmitted: (_) async {},
                next: (_) {},
                current: ds.SearchTray.Zero,
                empty: true,
              ),
              ds.TableHeader([
                Expanded(child: Text('Name')),
                Expanded(child: Text('Type')),
              ]),
            ],
          ),
        ),
        physicalSize: const Size(kWide, 500),
        fit: FlexFit.tight,
      );
      await tester.pumpAndSettle();

      expect(find.byType(ds.SearchTray), findsOneWidget);
      expect(find.text('Name'), findsOneWidget);
      expect(find.text('alpha'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  // ---------------------------------------------------------------------------
  // Overlay
  // ---------------------------------------------------------------------------
  group('Table overlay', () {
    testWidgets('overlay widget is shown over content', (tester) async {
      await tester.pumpApp(
        buildTable(overlay: Text('overlay-content')),
        physicalSize: const Size(kWide, 400),
        fit: FlexFit.tight,
      );
      await tester.pumpAndSettle();

      expect(find.text('overlay-content'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('default SizedBox overlay does not obstruct rows', (tester) async {
      await tester.pumpApp(
        buildTable(children: ['item-a', 'item-b']),
        physicalSize: const Size(kWide, 400),
        fit: FlexFit.tight,
      );
      await tester.pumpAndSettle();

      expect(find.text('item-a'), findsOneWidget);
      expect(find.text('item-b'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  // ---------------------------------------------------------------------------
  // Help integration
  // ---------------------------------------------------------------------------
  group('Table help', () {
    testWidgets('registers description with HelpScope', (tester) async {
      await tester.pumpApp(
        ds.HelpScope(
          ds.Table<String>(
            ds.Table.expanded((item) => Text(item)),
            children: ['a'],
            help: ds.Hint(const Text('list of items')),
          ),
        ),
      );
      await tester.pumpAndSettle();

      final scope = tester.state<ds.HelpScopeState>(find.byType(ds.HelpScope));
      expect(scope.descriptions, hasLength(1));
      expect(tester.takeException(), isNull);
    });

    testWidgets('description appears in help overlay', (tester) async {
      await tester.pumpApp(
        modals.Node(
          ds.HelpScope(
            ds.Table<String>(
              ds.Table.expanded((item) => Text(item)),
              children: ['a'],
              help: ds.Hint(const Text('list of items')),
              key: Key('table'),
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
      await tester.tap(find.descendant(of: find.byKey(Key('table')), matching: find.byType(InkWell)));
      await tester.pumpAndSettle();

      expect(find.text('list of items'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('defaults to HelpScope.None and does not register', (tester) async {
      await tester.pumpApp(
        ds.HelpScope(
          ds.Table<String>(
            ds.Table.expanded((item) => Text(item)),
            children: ['a'],
          ),
        ),
      );
      await tester.pumpAndSettle();

      final scope = tester.state<ds.HelpScopeState>(find.byType(ds.HelpScope));
      expect(scope.descriptions, hasLength(0));
      expect(tester.takeException(), isNull);
    });
  });
}
