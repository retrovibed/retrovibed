import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/filesystem.dart' as filesystem;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _photos = uuidx.withSuffix(1);

Future<filesystem.FilesystemSearchResponse> _search(
  filesystem.FilesystemSearchRequest req, {
  String? host,
  List<httpx.Option> options = const [],
}) async {
  return filesystem.FilesystemSearchResponse(
    next: req,
    breadcrumb: [],
    items: [
      media.Media(
        id: _photos,
        description: 'photos',
        mimetype: mimex.directory,
        createdAt: '2026-01-01T00:00:00Z',
        archiveId: uuidx.min(),
        torrentId: uuidx.min(),
        knownMediaId: uuidx.min(),
        directoryId: uuidx.min(),
      ),
    ],
  );
}

class _Deletes {
  final List<String> ids = [];

  Future<filesystem.FilesystemDeleteResponse> call(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    ids.add(id);
    return filesystem.FilesystemDeleteResponse();
  }
}

// the confirmation renders through the modal node, so the browser needs one above it or
// ds.modals.of returns null and the push is silently dropped.
Widget _harness(_Deletes deletes) => ds.Node(
  filesystem.FilesystemBrowser(
    search: _search,
    remove: deletes.call,
    mode: ValueNotifier(media.SearchMode.filesystem),
    onModeChanged: (_) {},
  ),
);

void main() {
  testWidgets('deleting a directory warns that its contents go with it', (WidgetTester tester) async {
    final deletes = _Deletes();
    await tester.pumpApp(_harness(deletes));
    await tester.pumpAndSettle();

    await tester.tap(find.byIcon(Icons.delete_outline));
    await tester.pumpAndSettle();

    expect(
      find.text('Delete photos? Everything inside it is removed from your library too.'),
      findsOneWidget,
    );

    // the warning is not advisory: nothing is deleted until it is answered.
    expect(deletes.ids, isEmpty);
  });

  testWidgets('cancelling the warning deletes nothing', (WidgetTester tester) async {
    final deletes = _Deletes();
    await tester.pumpApp(_harness(deletes));
    await tester.pumpAndSettle();

    await tester.tap(find.byIcon(Icons.delete_outline));
    await tester.pumpAndSettle();
    await tester.tap(find.text('No'));
    await tester.pumpAndSettle();

    expect(deletes.ids, isEmpty);
  });

  testWidgets('confirming the warning deletes the directory', (WidgetTester tester) async {
    final deletes = _Deletes();
    await tester.pumpApp(_harness(deletes));
    await tester.pumpAndSettle();

    await tester.tap(find.byIcon(Icons.delete_outline));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Yes'));
    await tester.pumpAndSettle();

    expect(deletes.ids, [_photos]);
  });
}
