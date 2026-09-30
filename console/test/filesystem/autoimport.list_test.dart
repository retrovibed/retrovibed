import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/filesystem/autoimport.api.dart' as api;
import 'package:retrovibed/filesystem/autoimport.edit.dart';
import 'package:retrovibed/filesystem/autoimport.list.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('AutoimportItem', () {
    testWidgets('row shows path, description, mode, and debounce', (WidgetTester tester) async {
      await tester.pumpApp(
        SingleChildScrollView(
          child: AutoimportItem(
            current: api.AutoimportDirectory(
              id: 'dir-1',
              path: '/tmp/inbox',
              description: 'movies',
              debounce: ds.Int64(5400),
              mode: api.autoimportModeMove,
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('/tmp/inbox'), findsOneWidget);
      expect(find.text('movies'), findsOneWidget);
      expect(find.text('move'), findsOneWidget);
      expect(find.text('1h30m'), findsOneWidget);
      expect(find.byType(AutoimportEdit), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('tapping the row expands the edit form', (WidgetTester tester) async {
      await tester.pumpApp(
        SingleChildScrollView(
          child: AutoimportItem(
            current: api.AutoimportDirectory(id: 'dir-1', path: '/tmp/inbox', debounce: ds.Int64(3600)),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byType(ds.TableRow).first);
      await tester.pumpAndSettle();

      expect(find.byType(AutoimportEdit), findsOneWidget);
      expect(find.text('save'), findsOneWidget);
    });

    testWidgets('save sends the edits and propagates the result', (WidgetTester tester) async {
      String? updatedID;
      api.AutoimportDirectoryUpdateRequest? sent;
      api.AutoimportDirectory? changed;

      await tester.pumpApp(
        SingleChildScrollView(
          child: AutoimportItem(
            current: api.AutoimportDirectory(id: 'dir-1', path: '/tmp/inbox', debounce: ds.Int64(3600)),
            onChange: (v) => changed = v,
            update: (id, req, {List<httpx.Option> options = const []}) async {
              updatedID = id;
              sent = req;
              return api.AutoimportDirectoryUpdateResponse(directory: req.directory);
            },
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byType(ds.TableRow).first);
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextFormField).at(1), 'movies');
      await tester.tap(find.text('save'));
      await tester.pumpAndSettle();

      expect(updatedID, 'dir-1');
      expect(sent?.directory.description, 'movies');
      expect(changed?.description, 'movies');
    });

    testWidgets('delete confirms, removes, and propagates null', (WidgetTester tester) async {
      String? deletedID;
      var removed = false;

      // the confirmation renders through the modal node, without one the push is dropped.
      await tester.pumpApp(
        ds.Node(
          SingleChildScrollView(
            child: AutoimportItem(
              current: api.AutoimportDirectory(id: 'dir-1', path: '/tmp/inbox', debounce: ds.Int64(3600)),
              onChange: (v) => removed = v == null,
              delete: (id, {List<httpx.Option> options = const []}) async {
                deletedID = id;
                return api.AutoimportDirectoryDeleteResponse();
              },
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.delete_outline));
      await tester.pumpAndSettle();

      expect(find.textContaining('stop monitoring /tmp/inbox'), findsOneWidget);
      await tester.tap(find.text('Yes'));
      await tester.pumpAndSettle();

      expect(deletedID, 'dir-1');
      expect(removed, isTrue);
    });
  });
}
