import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('TableRow autoexpand', () {
    testWidgets('autoexpand: false starts collapsed', (tester) async {
      await tester.pumpApp(
        ds.TableRow.single(
          Text('row'),
          expanded: Text('panel'),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('panel'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('autoexpand: false shows panel after one tap', (tester) async {
      await tester.pumpApp(
        ds.TableRow.single(
          Text('row'),
          expanded: Text('panel'),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byType(InkWell));
      await tester.pumpAndSettle();

      expect(find.text('panel'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('autoexpand: true starts expanded with no tap', (tester) async {
      await tester.pumpApp(
        ds.TableRow.single(
          Text('row'),
          expanded: Text('panel'),
          autoexpand: true,
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('panel'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('autoexpand: true still toggles closed on tap', (tester) async {
      await tester.pumpApp(
        ds.TableRow.single(
          Text('row'),
          expanded: Text('panel'),
          autoexpand: true,
        ),
      );
      await tester.pumpAndSettle();
      expect(find.text('panel'), findsOneWidget);

      await tester.tap(find.byType(InkWell));
      await tester.pumpAndSettle();

      expect(find.text('panel'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('autoexpand: true with no expandable panel uses onTap', (tester) async {
      int taps = 0;

      await tester.pumpApp(
        ds.TableRow.single(
          Text('row'),
          autoexpand: true,
          onTap: () => taps++,
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byType(InkWell));
      await tester.pumpAndSettle();

      expect(taps, 1);
      expect(tester.takeException(), isNull);
    });
  });
}
