import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/filesystem/autoimport.api.dart' as api;
import 'package:retrovibed/filesystem/autoimport.edit.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('AutoimportEdit', () {
    testWidgets('renders without overflow', (WidgetTester tester) async {
      await tester.pumpApp(
        Scaffold(
          body: AutoimportEdit(
            current: api.AutoimportDirectory(path: '/tmp/inbox', description: 'inbox', debounce: ds.Int64(3600)),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
      expect(find.text('/tmp/inbox'), findsOneWidget);
      expect(find.text('inbox'), findsOneWidget);
      expect(find.text('60'), findsOneWidget);
      expect(find.text('copy'), findsOneWidget);
    });

    testWidgets('renders without overflow in constrained environment', (WidgetTester tester) async {
      await tester.pumpApp(
        ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 300, maxHeight: 400),
          child: SingleChildScrollView(
            child: AutoimportEdit(
              current: api.AutoimportDirectory(path: '/tmp/inbox', debounce: ds.Int64(3600)),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('path is read only when editing an existing directory', (WidgetTester tester) async {
      await tester.pumpApp(
        Scaffold(
          body: AutoimportEdit(
            current: api.AutoimportDirectory(path: '/tmp/inbox', debounce: ds.Int64(3600)),
          ),
        ),
      );
      await tester.pumpAndSettle();

      final path = tester.widget<TextField>(find.widgetWithText(TextField, '/tmp/inbox'));
      expect(path.enabled, isFalse);
    });

    testWidgets('path is editable when creating', (WidgetTester tester) async {
      api.AutoimportDirectory? changed;
      await tester.pumpApp(
        Scaffold(
          body: AutoimportEdit(pathEditable: true, onChange: (v) => changed = v),
        ),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextFormField).at(0), '/tmp/new');
      expect(changed?.path, '/tmp/new');
    });

    testWidgets('edits description, debounce, and mode', (WidgetTester tester) async {
      api.AutoimportDirectory? changed;
      await tester.pumpApp(
        Scaffold(
          body: AutoimportEdit(
            current: api.AutoimportDirectory(path: '/tmp/inbox', debounce: ds.Int64(3600)),
            onChange: (v) => changed = v,
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextFormField).at(1), 'movies');
      expect(changed?.description, 'movies');

      await tester.enterText(find.byType(TextFormField).at(2), '5');
      expect(changed?.debounce, ds.Int64(300));

      await tester.tap(find.text('copy'));
      await tester.pumpAndSettle();
      await tester.tap(find.text('move').last);
      await tester.pumpAndSettle();
      expect(changed?.mode, api.autoimportModeMove);
    });

    testWidgets('ignores a non numeric debounce', (WidgetTester tester) async {
      api.AutoimportDirectory? changed;
      await tester.pumpApp(
        Scaffold(
          body: AutoimportEdit(
            current: api.AutoimportDirectory(path: '/tmp/inbox', debounce: ds.Int64(3600)),
            onChange: (v) => changed = v,
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextFormField).at(2), 'abc');
      expect(changed, isNull);
    });
  });
}
