import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/filesystem.dart' as filesystem;
import 'package:retrovibed/filesystem/details.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _photos = media.Media(
  id: uuidx.withSuffix(1),
  description: 'photos',
  mimetype: mimex.directory,
  createdAt: '2026-01-01T00:00:00Z',
  archiveId: uuidx.min(),
  torrentId: uuidx.min(),
  knownMediaId: uuidx.min(),
  directoryId: uuidx.min(),
);

// asyncfn keeps the LoadingIconButton in loading state (indefinite CircularProgressIndicator)
// until the confirmation is answered, so the delete tap is followed by pump() rather than
// pumpAndSettle(). the confirmation renders through the modal node, so the details need one
// above them or ds.modals.of returns null and the push is silently dropped.
void main() {
  group('FilesystemDetails', () {
    testWidgets('deleting warns that the contents go with it', (WidgetTester tester) async {
      final deleted = <String>[];
      await tester.pumpApp(
        ds.Node(
          FilesystemDetails(
            _photos,
            onChange: (_) {},
            apiremove: (String id, {List<httpx.Option> options = const []}) async {
              deleted.add(id);
              return filesystem.FilesystemDeleteResponse();
            },
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.delete));
      await tester.pump();

      expect(
        find.text('Delete photos? Everything inside it is removed from your library too.'),
        findsOneWidget,
      );

      // the warning is not advisory: nothing is deleted until it is answered.
      expect(deleted, isEmpty);
    });

    testWidgets('cancelling the warning deletes nothing', (WidgetTester tester) async {
      final deleted = <String>[];
      final changes = <media.Media?>[];
      await tester.pumpApp(
        ds.Node(
          FilesystemDetails(
            _photos,
            onChange: changes.add,
            apiremove: (String id, {List<httpx.Option> options = const []}) async {
              deleted.add(id);
              return filesystem.FilesystemDeleteResponse();
            },
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.delete));
      await tester.pump();
      await tester.tap(find.text('No'));
      await tester.pumpAndSettle();

      expect(deleted, isEmpty);
      expect(changes, isEmpty);
    });

    testWidgets('confirming the warning deletes and reports the removal', (WidgetTester tester) async {
      final deleted = <String>[];
      final changes = <media.Media?>[];
      await tester.pumpApp(
        ds.Node(
          FilesystemDetails(
            _photos,
            onChange: changes.add,
            apiremove: (String id, {List<httpx.Option> options = const []}) async {
              deleted.add(id);
              return filesystem.FilesystemDeleteResponse();
            },
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.delete));
      await tester.pump();
      await tester.tap(find.text('Yes'));
      await tester.pumpAndSettle();

      expect(deleted, [_photos.id]);
      // null tells the listing the entry is gone.
      expect(changes, [null]);
    });
  });
}
