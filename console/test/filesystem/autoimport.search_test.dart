import 'package:desktop_drop/desktop_drop.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/mimex.dart' as mimex;
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

      // the listing has items, so the tray's drop well is the only one.
      final well = tester.widget<ds.FileDropWell>(find.byType(ds.FileDropWell));
      expect(well.mimetypes, mimex.folders);
      await tester.runAsync(
        () => well.onDropped(
          ds.FilesEvent(files: [DropItemDirectory('/tmp/inbox', const [], mimeType: mimex.directory)]),
        ),
      );
      await tester.pumpAndSettle();

      // the selected directory opens the form with its path filled in.
      expect(find.byType(AutoImportEdit), findsOneWidget);
      // the path and the description, which defaults to the path.
      expect(find.text('/tmp/inbox'), findsNWidgets(2));
      expect(created, isNull);

      await tester.tap(find.text('create'));
      await tester.pumpAndSettle();

      expect(created?.directory.path, '/tmp/inbox');
      expect(created?.directory.debounce, ds.Int64(3600));
      expect(searches, 2);
      expect(find.byType(AutoImportEdit), findsNothing);
    });

    testWidgets('no monitored directories shows a directory drop well', (WidgetTester tester) async {
      await tester.pumpApp(
        AutoimportSearch(
          onClose: () {},
          authz: authn.AuthzCache.fake,
          search: (req, {List<httpx.Option> options = const []}) async =>
              api.AutoimportDirectorySearchResponse(next: req, items: []),
        ),
      );
      await tester.pumpAndSettle();

      // the empty state's drop well follows the tray's.
      final well = tester.widget<ds.FileDropWell>(find.byType(ds.FileDropWell).last);
      expect(well.mimetypes, mimex.folders);
      expect(find.text('Drop a directory to monitor it for automatic import.'), findsOneWidget);
      expect(find.byType(AutoImportEdit), findsNothing);
    });

    testWidgets('dropping a directory monitors it and refreshes', (WidgetTester tester) async {
      var searches = 0;
      final created = <api.AutoimportDirectoryCreateRequest>[];

      await tester.pumpApp(
        AutoimportSearch(
          onClose: () {},
          authz: authn.AuthzCache.fake,
          search: (req, {List<httpx.Option> options = const []}) async {
            searches++;
            return api.AutoimportDirectorySearchResponse(next: req, items: []);
          },
          create: (req, {List<httpx.Option> options = const []}) async {
            created.add(req);
            return api.AutoimportDirectoryCreateResponse(directory: req.directory);
          },
        ),
      );
      await tester.pumpAndSettle();
      expect(searches, 1);

      // the empty state's drop well follows the tray's.
      final well = tester.widget<ds.FileDropWell>(find.byType(ds.FileDropWell).last);
      await tester.runAsync(
        () => well.onDropped(
          ds.FilesEvent(files: [DropItemDirectory('/tmp/inbox', const [], mimeType: mimex.directory)]),
        ),
      );
      await tester.pumpAndSettle();

      expect(created, hasLength(1));
      expect(created.first.directory.path, '/tmp/inbox');
      expect(created.first.directory.debounce, ds.Int64(3600));
      expect(searches, 2);
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
