import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/errors.dart' as errors;
import 'package:retrovibed/design.kit/screens/error.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('ErrorScreen constrained parent', () {
    testWidgets('renders within fixed SizedBox constraints', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 150,
          child: ErrorScreen(
            Container(color: Colors.blue, child: const Text('Child')),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Child'), findsOneWidget);
      expect(tester.takeException(), isNull);

      final size = tester.getSize(find.byType(ErrorScreen).last);
      expect(size.width, equals(200));
      expect(size.height, equals(150));
    });

    testWidgets('renders in Column with fixed height', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Column(
          children: [
            SizedBox(
              height: 100,
              child: ErrorScreen(
                const Text('In column'),
              ),
            ),
            const Expanded(child: SizedBox()),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('In column'), findsOneWidget);
      expect(tester.takeException(), isNull);

      final size = tester.getSize(find.byType(ErrorScreen).last);
      expect(size.height, equals(100));
    });

    testWidgets('renders in Row with fixed width', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Row(
          children: [
            SizedBox(
              width: 150,
              child: ErrorScreen(
                const Text('In row'),
              ),
            ),
            const Expanded(child: SizedBox()),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('In row'), findsOneWidget);
      expect(tester.takeException(), isNull);

      final size = tester.getSize(find.byType(ErrorScreen).last);
      expect(size.width, equals(150));
    });

    testWidgets('renders with small dimensions', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Center(
          child: SizedBox(
            width: 50,
            height: 50,
            child: ErrorScreen(
              Container(color: Colors.blue),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);

      final size = tester.getSize(find.byType(ErrorScreen).last);
      expect(size.width, equals(50));
      expect(size.height, equals(50));
    });

    testWidgets('renders with zero width constraint', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Center(
          child: SizedBox(
            width: 0,
            height: 100,
            child: ErrorScreen(
              Container(color: Colors.blue),
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
        Center(
          child: SizedBox(
            width: 100,
            height: 0,
            child: ErrorScreen(
              Container(color: Colors.blue),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });
  });

  group('ErrorScreen unconstrained parent', () {
    testWidgets('renders in ListView with fixed height child', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ListView(
          children: [
            SizedBox(
              height: 200,
              child: ErrorScreen(
                const Text('In ListView'),
              ),
            ),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('In ListView'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in SingleChildScrollView with fixed height', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SingleChildScrollView(
          child: SizedBox(
            height: 300,
            child: ErrorScreen(
              const Text('In scroll'),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('In scroll'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in horizontal ListView with fixed width child', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          height: 100,
          child: ListView(
            scrollDirection: Axis.horizontal,
            children: [
              SizedBox(
                width: 200,
                child: ErrorScreen(
                  const Text('Horizontal'),
                ),
              ),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Horizontal'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('ErrorScreen size stability with errors', () {
    testWidgets('size unchanged in fixed SizedBox when error occurs', (
      WidgetTester tester,
    ) async {
      Widget cause = errors.Error.zero;

      await tester.pumpApp(
        StatefulBuilder(
          builder: (context, setState) {
            return SizedBox(
              width: 200,
              height: 150,
              child: ErrorScreen(
                cause: cause,
                TextButton(
                  onPressed: () {
                    setState(() {
                      cause = errors.Error.text('something went wrong');
                    });
                  },
                  child: const Text('Child'),
                ),
              ),
            );
          },
        ),
      );
      await tester.pumpAndSettle();

      final sizeBefore = tester.getSize(find.byType(ErrorScreen).last);
      expect(sizeBefore.width, equals(200));
      expect(sizeBefore.height, equals(150));

      await tester.tap(find.text('Child'));
      await tester.pumpAndSettle();

      final sizeAfter = tester.getSize(find.byType(ErrorScreen).last);
      expect(sizeAfter.width, equals(sizeBefore.width));
      expect(sizeAfter.height, equals(sizeBefore.height));
      expect(tester.takeException(), isNull);
    });

    testWidgets('size unchanged in Column when error occurs', (
      WidgetTester tester,
    ) async {
      Widget cause = errors.Error.zero;

      await tester.pumpApp(
        StatefulBuilder(
          builder: (context, setState) {
            return Column(
              children: [
                ErrorScreen(
                  SizedBox(
                    height: 100,
                    child: TextButton(
                      onPressed: () {
                        setState(() {
                          cause = errors.Error.text('something went wrong');
                        });
                      },
                      child: const Text('Child'),
                    ),
                  ),
                  cause: cause,
                ),
                const Expanded(child: SizedBox()),
              ],
            );
          },
        ),
      );
      await tester.pumpAndSettle();

      final sizeBefore = tester.getSize(find.byType(ErrorScreen).last);

      await tester.tap(find.text('Child'));
      await tester.pumpAndSettle();

      final sizeAfter = tester.getSize(find.byType(ErrorScreen).last);
      expect(sizeAfter.width, equals(sizeBefore.width));
      expect(sizeAfter.height, equals(sizeBefore.height));
      expect(tester.takeException(), isNull);
    });

    testWidgets('size unchanged in Row when error occurs', (
      WidgetTester tester,
    ) async {
      Widget cause = errors.Error.zero;

      await tester.pumpApp(
        StatefulBuilder(
          builder: (context, setState) {
            return Row(
              children: [
                ErrorScreen(
                  SizedBox(
                    width: 150,
                    child: TextButton(
                      onPressed: () {
                        setState(() {
                          cause = errors.Error.text('something went wrong');
                        });
                      },
                      child: const Text('Child'),
                    ),
                  ),
                  cause: cause,
                ),
                const Expanded(child: SizedBox()),
              ],
            );
          },
        ),
      );
      await tester.pumpAndSettle();

      final sizeBefore = tester.getSize(find.byType(ErrorScreen).last);

      await tester.tap(find.text('Child'));
      await tester.pumpAndSettle();

      final sizeAfter = tester.getSize(find.byType(ErrorScreen).last);
      expect(sizeAfter.width, equals(sizeBefore.width));
      expect(sizeAfter.height, equals(sizeBefore.height));
      expect(tester.takeException(), isNull);
    });

    testWidgets('size unchanged in ListView when error occurs', (
      WidgetTester tester,
    ) async {
      Widget cause = errors.Error.zero;

      await tester.pumpApp(
        StatefulBuilder(
          builder: (context, setState) {
            return ListView(
              children: [
                ErrorScreen(
                  SizedBox(
                    height: 200,
                    child: TextButton(
                      onPressed: () {
                        setState(() {
                          cause = errors.Error.text('something went wrong');
                        });
                      },
                      child: const Text('Child'),
                    ),
                  ),
                  cause: cause,
                ),
              ],
            );
          },
        ),
      );
      await tester.pumpAndSettle();

      final sizeBefore = tester.getSize(find.byType(ErrorScreen).last);

      await tester.tap(find.text('Child'));
      await tester.pumpAndSettle();

      final sizeAfter = tester.getSize(find.byType(ErrorScreen).last);
      expect(sizeAfter.width, equals(sizeBefore.width));
      expect(sizeAfter.height, equals(sizeBefore.height));
      expect(tester.takeException(), isNull);
    });

    testWidgets('size unchanged with tint when error occurs', (
      WidgetTester tester,
    ) async {
      Widget cause = errors.Error.zero;
      final tint = [
        BoxShadow(
          color: const Color(0x666A0101),
          spreadRadius: 1,
          blurRadius: 10,
        ),
      ];

      await tester.pumpApp(
        StatefulBuilder(
          builder: (context, setState) {
            return SizedBox(
              width: 200,
              height: 150,
              child: ErrorScreen(
                TextButton(
                  onPressed: () {
                    setState(() {
                      cause = errors.Error.text('something went wrong');
                    });
                  },
                  child: const Text('Child'),
                ),
                cause: cause,
                tint: tint,
              ),
            );
          },
        ),
      );
      await tester.pumpAndSettle();

      final sizeBefore = tester.getSize(find.byType(ErrorScreen).last);
      expect(sizeBefore.width, equals(200));
      expect(sizeBefore.height, equals(150));

      await tester.tap(find.text('Child'));
      await tester.pumpAndSettle();

      final sizeAfter = tester.getSize(find.byType(ErrorScreen).last);
      expect(sizeAfter.width, equals(sizeBefore.width));
      expect(sizeAfter.height, equals(sizeBefore.height));
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders with ListView shrinkWrap - natural height when no error', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ListView(
          shrinkWrap: true,
          padding: EdgeInsets.zero,
          children: [
            ErrorScreen(
              const Text("derp"),
            ),
          ],
        ),
      );

      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);

      final size = tester.getSize(find.byType(Text));
      expect(size.height, equals(20.0));
    });
  });
}
