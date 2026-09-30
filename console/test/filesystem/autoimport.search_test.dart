import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/filesystem/autoimport.api.dart' as api;
import 'package:retrovibed/filesystem/autoimport.edit.dart';
import 'package:retrovibed/filesystem/autoimport.search.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('AutoimportSearch', () {
    testWidgets('lists the monitored directories', (WidgetTester tester) async {
      await tester.pumpApp(
        AutoimportSearch(
          onClose: () {},
          authz: authn.AuthzCache.fake,
          search: (req, {List<httpx.Option> options = const []}) async => api.AutoimportDirectorySearchResponse(
            next: req,
            items: [
              api.AutoimportDirectory(id: 'dir-1', path: '/tmp/movies', debounce: ds.Int64(3600)),
              api.AutoimportDirectory(id: 'dir-2', path: '/tmp/music', debounce: ds.Int64(3600)),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('/tmp/movies'), findsOneWidget);
      expect(find.text('/tmp/music'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('submitting a query searches from the first page', (WidgetTester tester) async {
      final requests = <api.AutoimportDirectorySearchRequest>[];

      await tester.pumpApp(
        AutoimportSearch(
          onClose: () {},
          authz: authn.AuthzCache.fake,
          search: (req, {List<httpx.Option> options = const []}) async {
            requests.add(req.deepCopy());
            return api.AutoimportDirectorySearchResponse(next: req, items: []);
          },
        ),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField).first, 'movies');
      await tester.testTextInput.receiveAction(TextInputAction.done);
      await tester.pumpAndSettle();

      expect(requests.last.query, 'movies');
      expect(requests.last.offset, ds.Int64(0));
    });

    testWidgets('creating a directory submits it and refreshes', (WidgetTester tester) async {
      var searches = 0;
      api.AutoimportDirectoryCreateRequest? created;

      await tester.pumpApp(
        AutoimportSearch(
          onClose: () {},
          authz: authn.AuthzCache.fake,
          search: (req, {List<httpx.Option> options = const []}) async {
            searches++;
            return api.AutoimportDirectorySearchResponse(
              next: req,
              items: [api.AutoimportDirectory(id: 'dir-1', path: '/tmp/movies', debounce: ds.Int64(3600))],
            );
          },
          create: (req, {List<httpx.Option> options = const []}) async {
            created = req;
            return api.AutoimportDirectoryCreateResponse(directory: req.directory);
          },
        ),
      );
      await tester.pumpAndSettle();
      expect(searches, 1);

      await tester.tap(find.byIcon(Icons.add));
      await tester.pumpAndSettle();
      expect(find.byType(AutoimportEdit), findsOneWidget);

      await tester.enterText(find.byType(TextFormField).at(0), '/tmp/inbox');
      await tester.tap(find.text('create'));
      await tester.pumpAndSettle();

      expect(created?.directory.path, '/tmp/inbox');
      expect(created?.directory.debounce, ds.Int64(3600));
      expect(searches, 2);
      expect(find.byType(AutoimportEdit), findsNothing);
    });

    testWidgets('close button invokes onClose', (WidgetTester tester) async {
      var closed = false;

      await tester.pumpApp(
        AutoimportSearch(
          onClose: () => closed = true,
          authz: authn.AuthzCache.fake,
          search: (req, {List<httpx.Option> options = const []}) async =>
              api.AutoimportDirectorySearchResponse(next: req, items: []),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.close));
      await tester.pumpAndSettle();

      expect(closed, isTrue);
    });
  });
}
