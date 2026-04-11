import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/forms/field.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('Field layout', () {
    group('without label', () {
      testWidgets('renders input without overflow in unconstrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Field(input: TextField()),
        );
        await tester.pumpAndSettle();

        expect(find.byType(Field), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders input without overflow in narrow SizedBox (200px)', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 200,
            child: Field(input: TextField()),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(Field), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in very narrow SizedBox (100px)', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 100,
            child: Field(input: TextField()),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(Field), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('with short label', () {
      testWidgets('renders label and input without overflow at default size', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Field(
            label: Text('port'),
            input: TextField(),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('port'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in narrow SizedBox (200px)', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 200,
            child: Field(
              label: Text('port'),
              input: TextField(),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('port'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in very narrow SizedBox (100px)', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 100,
            child: Field(
              label: Text('port'),
              input: TextField(),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(Field), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('with long label', () {
      testWidgets('renders without overflow at default size', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Field(
            label: Text('download rate'),
            input: TextField(),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('download rate'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in narrow SizedBox (200px)', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 200,
            child: Field(
              label: Text('download rate'),
              input: TextField(),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(Field), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in narrow SizedBox (216px)', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 216,
            child: Field(
              label: Text('download rate'),
              input: TextField(),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(Field), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in SizedBox (300px)', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 300,
            child: Field(
              label: Text('download rate'),
              input: TextField(),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('download rate'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('with very long label', () {
      testWidgets('renders without overflow at default size', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Field(
            label: Text('an extremely long label that could potentially overflow'),
            input: TextField(),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(Field), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in narrow SizedBox (200px)', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 200,
            child: Field(
              label: Text(
                'an extremely long label that could potentially overflow',
              ),
              input: TextField(),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(Field), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('with rich label', () {
      testWidgets('renders Row label without overflow at narrow width', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 216,
            child: Field(
              label: Row(
                children: [
                  Icon(Icons.info, size: 16),
                  SizedBox(width: 4),
                  Text('download rate'),
                ],
              ),
              input: TextField(),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(Field), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('in Wrap', () {
      testWidgets('multiple fields in Wrap do not overflow at 400px', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 400,
            child: Wrap(
              children: [
                Field(label: Text('download rate'), input: TextField()),
                Field(label: Text('upload rate'), input: TextField()),
                Field(label: Text('peers'), input: TextField()),
              ],
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('download rate'), findsOneWidget);
        expect(find.text('upload rate'), findsOneWidget);
        expect(find.text('peers'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('multiple fields in Wrap do not overflow at 216px', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 216,
            child: Wrap(
              children: [
                Field(label: Text('download rate'), input: TextField()),
                Field(label: Text('upload rate'), input: TextField()),
              ],
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(Field), findsNWidgets(2));
        expect(tester.takeException(), isNull);
      });
    });

    group('in Column', () {
      testWidgets('multiple fields in Column do not overflow at 200px', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 200,
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Field(label: Text('download rate'), input: TextField()),
                Field(label: Text('upload rate'), input: TextField()),
                Field(label: Text('peers'), input: TextField()),
              ],
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(Field), findsNWidgets(3));
        expect(tester.takeException(), isNull);
      });
    });
  });
}
