import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('Confirmation widget rendering', () {
    testWidgets('renders with custom content and buttons', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ds.Confirmation(
          content: Text('Are you sure?'),
          confirmation: Text('Confirm'),
          cancellation: Text('Deny'),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Are you sure?'), findsOneWidget);
      expect(find.text('Confirm'), findsOneWidget);
      expect(find.text('Deny'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders with complex content widget', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ds.Confirmation(
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text('Delete Item'),
              Text('This action cannot be undone.'),
            ],
          ),
          confirmation: Text('Delete'),
          cancellation: Text('Keep'),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Delete Item'), findsOneWidget);
      expect(find.text('This action cannot be undone.'), findsOneWidget);
      expect(find.text('Delete'), findsOneWidget);
      expect(find.text('Keep'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders with icon buttons', (WidgetTester tester) async {
      await tester.pumpApp(
        ds.Confirmation(
          content: Text('Confirm action?'),
          confirmation: Icon(Icons.check),
          cancellation: Icon(Icons.close),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Confirm action?'), findsOneWidget);
      expect(find.byIcon(Icons.check), findsOneWidget);
      expect(find.byIcon(Icons.close), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('Confirmation factory methods', () {
    testWidgets('yesNo creates dialog with Yes and No buttons', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ds.Confirmation.yesNo(content: Text('Continue?')),
      );
      await tester.pumpAndSettle();

      expect(find.text('Continue?'), findsOneWidget);
      expect(find.text('Yes'), findsOneWidget);
      expect(find.text('No'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('createCancel creates dialog with Create and Cancel buttons', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ds.Confirmation.createCancel(content: Text('Create new item?')),
      );
      await tester.pumpAndSettle();

      expect(find.text('Create new item?'), findsOneWidget);
      expect(find.text('Create'), findsOneWidget);
      expect(find.text('Cancel'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('Confirmation callbacks', () {
    testWidgets('onConfirm is called when confirmation is tapped', (
      WidgetTester tester,
    ) async {
      bool confirmed = false;

      await tester.pumpApp(
        ds.Confirmation.yesNo(
          content: Text('Confirm?'),
          onConfirm: () => confirmed = true,
        ),
      );
      await tester.pumpAndSettle();

      expect(confirmed, isFalse);

      await tester.tap(find.text('Yes'));
      await tester.pumpAndSettle();

      expect(confirmed, isTrue);
      expect(tester.takeException(), isNull);
    });

    testWidgets('onCancel is called when cancellation is tapped', (
      WidgetTester tester,
    ) async {
      bool cancelled = false;

      await tester.pumpApp(
        ds.Confirmation.yesNo(
          content: Text('Confirm?'),
          onCancel: () => cancelled = true,
        ),
      );
      await tester.pumpAndSettle();

      expect(cancelled, isFalse);

      await tester.tap(find.text('No'));
      await tester.pumpAndSettle();

      expect(cancelled, isTrue);
      expect(tester.takeException(), isNull);
    });

    testWidgets('both callbacks work independently', (
      WidgetTester tester,
    ) async {
      int confirmCount = 0;
      int cancelCount = 0;

      await tester.pumpApp(
        ds.Confirmation.createCancel(
          content: Text('Action?'),
          onConfirm: () => confirmCount++,
          onCancel: () => cancelCount++,
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('Create'));
      await tester.pumpAndSettle();
      expect(confirmCount, 1);
      expect(cancelCount, 0);

      await tester.tap(find.text('Cancel'));
      await tester.pumpAndSettle();
      expect(confirmCount, 1);
      expect(cancelCount, 1);

      expect(tester.takeException(), isNull);
    });

    testWidgets('null callbacks do not throw', (WidgetTester tester) async {
      await tester.pumpApp(
        ds.Confirmation.yesNo(content: Text('No callbacks')),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('Yes'));
      await tester.pumpAndSettle();

      await tester.tap(find.text('No'));
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });
  });

  group('Confirmation finite constraints', () {
    testWidgets('renders in finite container', (WidgetTester tester) async {
      await tester.pumpApp(
        SizedBox(
          width: 400,
          height: 200,
          child: ds.Confirmation.yesNo(content: Text('Constrained')),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Constrained'), findsOneWidget);
      expect(find.byType(ds.Confirmation), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in Column with constrained parent', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 400,
          height: 300,
          child: Column(
            children: [
              ds.Confirmation.yesNo(content: Text('In column')),
              Text('Other content'),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('In column'), findsOneWidget);
      expect(find.text('Other content'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('Confirmation infinite constraints', () {
    testWidgets('renders in ListView', (WidgetTester tester) async {
      await tester.pumpApp(
        ListView(
          children: [
            ds.Confirmation.yesNo(content: Text('In list')),
            Text('Item below'),
          ],
        ),
        fit: FlexFit.tight,
      );
      await tester.pumpAndSettle();

      expect(find.text('In list'), findsOneWidget);
      expect(find.text('Item below'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in SingleChildScrollView', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SingleChildScrollView(
          child: Column(
            children: [
              ds.Confirmation.createCancel(content: Text('Scrollable')),
              Text('More content'),
            ],
          ),
        ),
        fit: FlexFit.tight,
      );
      await tester.pumpAndSettle();

      expect(find.text('Scrollable'), findsOneWidget);
      expect(find.text('More content'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders multiple Confirmations in Column', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SingleChildScrollView(
          child: Column(
            children: [
              ds.Confirmation.yesNo(content: Text('First')),
              ds.Confirmation.createCancel(content: Text('Second')),
            ],
          ),
        ),
        fit: FlexFit.tight,
      );
      await tester.pumpAndSettle();

      expect(find.text('First'), findsOneWidget);
      expect(find.text('Second'), findsOneWidget);
      expect(find.byType(ds.Confirmation), findsNWidgets(2));
      expect(tester.takeException(), isNull);
    });
  });

  group('Confirmation InkWell behavior', () {
    testWidgets('confirmation button has InkWell', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ds.Confirmation.yesNo(content: Text('Test')),
      );
      await tester.pumpAndSettle();

      expect(find.byType(InkWell), findsNWidgets(2));
      expect(tester.takeException(), isNull);
    });

    testWidgets('InkWell wraps button content', (WidgetTester tester) async {
      await tester.pumpApp(
        ds.Confirmation(
          content: Text('Content'),
          confirmation: Text('OK'),
          cancellation: Text('Back'),
        ),
      );
      await tester.pumpAndSettle();

      final okText = find.text('OK');
      final backText = find.text('Back');

      expect(
        find.ancestor(of: okText, matching: find.byType(InkWell)),
        findsOneWidget,
      );
      expect(
        find.ancestor(of: backText, matching: find.byType(InkWell)),
        findsOneWidget,
      );
      expect(tester.takeException(), isNull);
    });
  });

  group('Confirmation button cursors', () {
    testWidgets('confirmation button shows click cursor on hover', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ds.Confirmation.yesNo(content: Text('Delete?')),
      );
      await tester.pumpAndSettle();

      expect(
        tester.resolvedCursorAt(find.text('Yes')),
        SystemMouseCursors.click,
      );
      expect(tester.takeException(), isNull);
    });

    testWidgets('cancellation button shows click cursor on hover', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ds.Confirmation.yesNo(content: Text('Delete?')),
      );
      await tester.pumpAndSettle();

      expect(
        tester.resolvedCursorAt(find.text('No')),
        SystemMouseCursors.click,
      );
      expect(tester.takeException(), isNull);
    });
  });
}
