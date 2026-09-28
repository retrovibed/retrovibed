import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/filesystem.dart' as filesystem;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _held = media.Media(
  id: uuidx.withSuffix(1),
  description: 'held.bin',
  mimetype: 'application/octet-stream',
  createdAt: '2026-01-01T00:00:00Z',
  archiveId: uuidx.min(),
  torrentId: uuidx.min(),
  knownMediaId: uuidx.min(),
  directoryId: uuidx.min(),
);

final _photos = media.Media(
  id: uuidx.withSuffix(2),
  description: 'photos',
  mimetype: mimex.directory,
  createdAt: '2026-01-01T00:00:00Z',
  archiveId: uuidx.min(),
  torrentId: uuidx.min(),
  knownMediaId: uuidx.min(),
  directoryId: uuidx.min(),
);

// the button stays in its loading state (an indefinite CircularProgressIndicator) while the
// picker is open, so taps are followed by pump() rather than pumpAndSettle(). the picker
// renders through the modal node, so the button needs one above it.
void main() {
  group('ButtonFileMove', () {
    testWidgets('moves into the picked directory', (WidgetTester tester) async {
      final searched = <List<String>>[];
      final moved = <(String, String)>[];
      final changed = <media.Media?>[];

      await tester.pumpApp(
        ds.Node(
          filesystem.ButtonFileMove(
            current: _held,
            onChange: changed.add,
            apisearch: (media.MediaSearchRequest req, {String? host, List<httpx.Option> options = const []}) async {
              searched.add(List.of(req.mimetypes));
              return media.MediaSearchResponse(next: req, items: [_photos]);
            },
            apimove: (String id, filesystem.FilesystemMoveRequest req, {List<httpx.Option> options = const []}) async {
              moved.add((id, req.directoryId));
              return filesystem.FilesystemMoveResponse(media: _held.deepCopy()..directoryId = req.directoryId);
            },
          ),
        ),
      );
      await tester.pump();

      await tester.tap(find.byIcon(Icons.drive_file_move_outline));
      await tester.pump();
      await tester.pump();

      expect(searched, [
        [mimex.directory],
      ]);

      await tester.tap(find.text('photos'));
      await tester.pump();
      await tester.pump();

      expect(moved, [(_held.id, _photos.id)]);
      expect(changed, [null]);
      expect(find.text('move held.bin'), findsNothing);
    });

    testWidgets('the root row moves to the root', (WidgetTester tester) async {
      final moved = <String>[];

      await tester.pumpApp(
        ds.Node(
          filesystem.ButtonFileMove(
            current: _held,
            apisearch: (media.MediaSearchRequest req, {String? host, List<httpx.Option> options = const []}) async {
              return media.MediaSearchResponse(next: req, items: [_photos]);
            },
            apimove: (String id, filesystem.FilesystemMoveRequest req, {List<httpx.Option> options = const []}) async {
              moved.add(req.directoryId);
              return filesystem.FilesystemMoveResponse(media: _held);
            },
          ),
        ),
      );
      await tester.pump();

      await tester.tap(find.byIcon(Icons.drive_file_move_outline));
      await tester.pump();
      await tester.pump();

      await tester.tap(find.text('/'));
      await tester.pump();
      await tester.pump();

      expect(moved, [uuidx.min()]);
    });

    testWidgets('a directory is not offered as its own destination', (WidgetTester tester) async {
      await tester.pumpApp(
        ds.Node(
          filesystem.ButtonFileMove(
            current: _photos,
            apisearch: (media.MediaSearchRequest req, {String? host, List<httpx.Option> options = const []}) async {
              return media.MediaSearchResponse(next: req, items: [_photos]);
            },
            apimove: (String id, filesystem.FilesystemMoveRequest req, {List<httpx.Option> options = const []}) async {
              return filesystem.FilesystemMoveResponse(media: _photos);
            },
          ),
        ),
      );
      await tester.pump();

      await tester.tap(find.byIcon(Icons.drive_file_move_outline));
      await tester.pump();
      await tester.pump();

      expect(find.text('move photos'), findsOneWidget);
      expect(find.text('photos'), findsNothing);
    });

    testWidgets('closing moves nothing', (WidgetTester tester) async {
      final moved = <String>[];
      final changed = <media.Media?>[];

      await tester.pumpApp(
        ds.Node(
          filesystem.ButtonFileMove(
            current: _held,
            onChange: changed.add,
            apisearch: (media.MediaSearchRequest req, {String? host, List<httpx.Option> options = const []}) async {
              return media.MediaSearchResponse(next: req, items: [_photos]);
            },
            apimove: (String id, filesystem.FilesystemMoveRequest req, {List<httpx.Option> options = const []}) async {
              moved.add(id);
              return filesystem.FilesystemMoveResponse(media: _held);
            },
          ),
        ),
      );
      await tester.pump();

      await tester.tap(find.byIcon(Icons.drive_file_move_outline));
      await tester.pump();
      await tester.pump();

      await tester.tap(find.byIcon(Icons.close));
      await tester.pump();
      await tester.pump();

      expect(find.text('move held.bin'), findsNothing);
      expect(moved, isEmpty);
      expect(changed, isEmpty);
    });
  });
}
