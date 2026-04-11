import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('FileDropWell', () {
    group('constrained contexts', () {
      testWidgets('renders within SizedBox with fixed dimensions', (WidgetTester tester) async {
        await tester.pumpApp(
          SizedBox(
            width: 200,
            height: 200,
            child: ds.FileDropWell(
              (_, {ValueNotifier<int>? progress}) async => null,
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(ds.FileDropWell), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders within Container with fixed dimensions', (WidgetTester tester) async {
        await tester.pumpApp(
          Container(
            width: 300,
            height: 150,
            child: ds.FileDropWell(
              (_, {ValueNotifier<int>? progress}) async => null,
              child: const Text('Drop here'),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(ds.FileDropWell), findsOneWidget);
        expect(find.text('Drop here'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders within Column with Expanded wrapper', (WidgetTester tester) async {
        await tester.pumpApp(
          Column(
            children: [
              const Text('Header'),
              Expanded(
                child: ds.FileDropWell(
                  (_, {ValueNotifier<int>? progress}) async => null,
                  child: const Text('Drop zone'),
                ),
              ),
              const Text('Footer'),
            ],
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(ds.FileDropWell), findsOneWidget);
        expect(find.text('Drop zone'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders with custom child in bounded context', (WidgetTester tester) async {
        await tester.pumpApp(
          Center(
            child: SizedBox(
              width: 100,
              height: 100,
              child: ds.FileDropWell(
                (_, {ValueNotifier<int>? progress}) async => null,
                child: const Icon(Icons.upload),
              ),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byIcon(Icons.upload), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('icon factory renders in constrained context', (WidgetTester tester) async {
        await tester.pumpApp(
          SizedBox(
            width: 48,
            height: 48,
            child: ds.FileDropWell.icon(
              (_, {ValueNotifier<int>? progress}) async => null,
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byIcon(Icons.file_upload_outlined), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('unconstrained/scrollable contexts', () {
      testWidgets('renders within ListView when given explicit height', (WidgetTester tester) async {
        await tester.pumpApp(
          ListView(
            children: [
              const Text('Item 1'),
              SizedBox(
                height: 200,
                child: ds.FileDropWell(
                  (_, {ValueNotifier<int>? progress}) async => null,
                  child: const Text('Drop zone'),
                ),
              ),
              const Text('Item 2'),
            ],
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(ds.FileDropWell), findsOneWidget);
        expect(find.text('Drop zone'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders within SingleChildScrollView when given explicit height', (WidgetTester tester) async {
        await tester.pumpApp(
          SingleChildScrollView(
            child: Column(
              children: [
                const Text('Header'),
                SizedBox(
                  height: 150,
                  child: ds.FileDropWell(
                    (_, {ValueNotifier<int>? progress}) async => null,
                    child: const Text('Scrollable drop zone'),
                  ),
                ),
                const Text('Footer'),
              ],
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(ds.FileDropWell), findsOneWidget);
        expect(find.text('Scrollable drop zone'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders within horizontal Row when given explicit width', (WidgetTester tester) async {
        await tester.pumpApp(
          Row(
            children: [
              const Text('Left'),
              SizedBox(
                width: 100,
                height: 100,
                child: ds.FileDropWell(
                  (_, {ValueNotifier<int>? progress}) async => null,
                  child: const Icon(Icons.add),
                ),
              ),
              const Text('Right'),
            ],
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(ds.FileDropWell), findsOneWidget);
        expect(find.byIcon(Icons.add), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders within horizontal ListView when given explicit width', (WidgetTester tester) async {
        await tester.pumpApp(
          SizedBox(
            height: 100,
            child: ListView(
              scrollDirection: Axis.horizontal,
              children: [
                const SizedBox(width: 50, child: Text('A')),
                SizedBox(
                  width: 100,
                  child: ds.FileDropWell(
                    (_, {ValueNotifier<int>? progress}) async => null,
                    child: const Text('Drop'),
                  ),
                ),
                const SizedBox(width: 50, child: Text('B')),
              ],
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(ds.FileDropWell), findsOneWidget);
        expect(find.text('Drop'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('icon factory renders in scrollable context with constraints', (WidgetTester tester) async {
        await tester.pumpApp(
          ListView(
            children: [
              SizedBox(
                height: 48,
                width: 48,
                child: ds.FileDropWell.icon(
                  (_, {ValueNotifier<int>? progress}) async => null,
                ),
              ),
            ],
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byIcon(Icons.file_upload_outlined), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('default child behavior', () {
      testWidgets('default child renders in full-screen constrained context', (WidgetTester tester) async {
        await tester.pumpApp(
          ds.FileDropWell(
            (_, {ValueNotifier<int>? progress}) async => null,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(ds.FileDropWell), findsOneWidget);
        expect(find.byIcon(Icons.filter_rounded), findsOneWidget);
        expect(find.text('Drop Files to add them to your library.'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('default child with MainAxisSize.max requires bounded height in Column', (WidgetTester tester) async {
        await tester.pumpApp(
          Column(
            children: [
              Expanded(
                child: ds.FileDropWell(
                  (_, {ValueNotifier<int>? progress}) async => null,
                ),
              ),
            ],
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(ds.FileDropWell), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('custom child with MainAxisSize.min works in flexible contexts', (WidgetTester tester) async {
        await tester.pumpApp(
          SingleChildScrollView(
            child: Column(
              children: [
                ds.FileDropWell(
                  (_, {ValueNotifier<int>? progress}) async => null,
                  child: const Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Icon(Icons.upload),
                      Text('Upload files'),
                    ],
                  ),
                ),
              ],
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(ds.FileDropWell), findsOneWidget);
        expect(find.byIcon(Icons.upload), findsOneWidget);
        expect(find.text('Upload files'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('layout dimensions', () {
      testWidgets('respects margin parameter', (WidgetTester tester) async {
        const testMargin = EdgeInsets.all(16.0);
        await tester.pumpApp(
          SizedBox(
            width: 200,
            height: 200,
            child: ds.FileDropWell(
              (_, {ValueNotifier<int>? progress}) async => null,
              margin: testMargin,
              child: const Text('Margined'),
            ),
          ),
        );
        await tester.pumpAndSettle();

        final container = tester.widget<Container>(
          find.descendant(
            of: find.byType(ds.FileDropWell),
            matching: find.byType(Container),
          ),
        );
        expect(container.margin, testMargin);
      });

      testWidgets('fills available space in Expanded widget', (WidgetTester tester) async {
        await tester.pumpApp(
          Column(
            children: [
              SizedBox(height: 50, child: const Text('Fixed')),
              Expanded(
                child: ds.FileDropWell(
                  (_, {ValueNotifier<int>? progress}) async => null,
                  child: const Text('Fills space'),
                ),
              ),
            ],
          ),
        );
        await tester.pumpAndSettle();

        final fileDropWellBox = tester.getSize(find.byType(ds.FileDropWell));
        final scaffoldSize = tester.getSize(find.byType(Scaffold));

        expect(fileDropWellBox.height, scaffoldSize.height - 50);
        expect(tester.takeException(), isNull);
      });

      testWidgets('adapts to tight constraints', (WidgetTester tester) async {
        const tightSize = Size(50, 50);
        await tester.pumpApp(
          Center(
            child: SizedBox(
              width: tightSize.width,
              height: tightSize.height,
              child: ds.FileDropWell.icon(
                (_, {ValueNotifier<int>? progress}) async => null,
              ),
            ),
          ),
        );
        await tester.pumpAndSettle();

        final widgetSize = tester.getSize(find.byType(ds.FileDropWell));
        expect(widgetSize.width, tightSize.width);
        expect(widgetSize.height, tightSize.height);
        expect(tester.takeException(), isNull);
      });
    });

    group('nested in complex layouts', () {
      testWidgets('renders in GridView cell with fixed extent', (WidgetTester tester) async {
        await tester.pumpApp(
          GridView.extent(
            maxCrossAxisExtent: 150,
            childAspectRatio: 1.0,
            children: [
              ds.FileDropWell(
                (_, {ValueNotifier<int>? progress}) async => null,
                child: const Icon(Icons.upload),
              ),
              Container(color: Colors.blue),
              Container(color: Colors.green),
            ],
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(ds.FileDropWell), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders in Stack', (WidgetTester tester) async {
        await tester.pumpApp(
          Stack(
            children: [
              Container(color: Colors.grey),
              Positioned(
                top: 50,
                left: 50,
                width: 200,
                height: 200,
                child: ds.FileDropWell(
                  (_, {ValueNotifier<int>? progress}) async => null,
                  child: const Text('Stacked'),
                ),
              ),
            ],
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(ds.FileDropWell), findsOneWidget);
        expect(find.text('Stacked'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders in Flex with bounded cross axis', (WidgetTester tester) async {
        await tester.pumpApp(
          Flex(
            direction: Axis.vertical,
            children: [
              Flexible(
                child: ds.FileDropWell(
                  (_, {ValueNotifier<int>? progress}) async => null,
                  child: const Text('Flex child'),
                ),
              ),
            ],
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(ds.FileDropWell), findsOneWidget);
        expect(find.text('Flex child'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });
  });
}
