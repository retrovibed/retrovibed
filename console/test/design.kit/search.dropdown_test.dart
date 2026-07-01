import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('SearchDropdown onSearch returns ds.Empty', () {
    testWidgets('no results panel shown and size unchanged', (WidgetTester tester) async {
      await tester.pumpApp(
        ds.SearchDropdown(
          onSearch: (query, onClick) async => ds.Empty,
        ),
      );
      await tester.pumpAndSettle();

      final sizeBefore = tester.getSize(find.byType(ds.SearchDropdown));

      await tester.enterText(find.byType(TextField), 'hello');
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
      expect(tester.getSize(find.byType(ds.SearchDropdown)), equals(sizeBefore));
      expect(find.byType(Visibility), findsWidgets);

      final visibilities = tester.widgetList<Visibility>(find.byType(Visibility)).toList();
      expect(visibilities.any((v) => v.visible), isFalse);
    });

    testWidgets('no results panel shown after clearing text', (WidgetTester tester) async {
      int callCount = 0;
      await tester.pumpApp(
        ds.SearchDropdown(
          onSearch: (query, onClick) async {
            callCount++;
            return ds.Empty;
          },
        ),
      );
      await tester.pumpAndSettle();

      final sizeBefore = tester.getSize(find.byType(ds.SearchDropdown));

      await tester.enterText(find.byType(TextField), 'search');
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), '');
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
      expect(tester.getSize(find.byType(ds.SearchDropdown)), equals(sizeBefore));
      expect(callCount, greaterThan(0));
    });

    testWidgets('results panel shown when onSearch returns non-Empty', (WidgetTester tester) async {
      await tester.pumpApp(
        ds.SearchDropdown(
          onSearch: (query, onClick) async => const Text('result'),
        ),
      );
      await tester.pumpAndSettle();

      final sizeBefore = tester.getSize(find.byType(ds.SearchDropdown));

      await tester.enterText(find.byType(TextField), 'query');
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
      expect(find.text('result'), findsOneWidget);
      expect(tester.getSize(find.byType(ds.SearchDropdown)).height, greaterThan(sizeBefore.height));
    });

    testWidgets('switches from non-Empty to ds.Empty collapses panel', (WidgetTester tester) async {
      var returnEmpty = false;
      await tester.pumpApp(
        StatefulBuilder(
          builder: (context, setState) {
            return ds.SearchDropdown(
              onSearch: (query, onClick) async {
                if (returnEmpty) return ds.Empty;
                return const Text('result');
              },
            );
          },
        ),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), 'query');
      await tester.pumpAndSettle();

      expect(find.text('result'), findsOneWidget);
      final sizeWithResults = tester.getSize(find.byType(ds.SearchDropdown));

      returnEmpty = true;
      await tester.enterText(find.byType(TextField), 'query2');
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
      expect(find.text('result'), findsNothing);
      expect(tester.getSize(find.byType(ds.SearchDropdown)).height, lessThan(sizeWithResults.height));
    });
  });
}
