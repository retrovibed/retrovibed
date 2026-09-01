import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/filesystem.dart' as filesystem;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _photos = uuidx.withSuffix(1);
final _nested = uuidx.withSuffix(2);

media.Media _directory(String id, String description, {String parent = ""}) => media.Media(
  id: id,
  description: description,
  mimetype: mimex.directory,
  createdAt: '2026-01-01T00:00:00Z',
  archiveId: uuidx.min(),
  torrentId: uuidx.min(),
  knownMediaId: uuidx.min(),
  directoryId: parent.isEmpty ? uuidx.min() : parent,
);

media.Media _file(String id, String description, {String parent = ""}) => media.Media(
  id: id,
  description: description,
  mimetype: 'audio/mp3',
  createdAt: '2026-01-01T00:00:00Z',
  archiveId: uuidx.min(),
  torrentId: uuidx.min(),
  knownMediaId: uuidx.min(),
  directoryId: parent.isEmpty ? uuidx.min() : parent,
);

// records what the browser asked for so navigation can be asserted on the request rather
// than on rendering alone.
class _Searches {
  final List<String> directories = [];

  Future<filesystem.FilesystemSearchResponse> call(
    filesystem.FilesystemSearchRequest req, {
    String? host,
    List<httpx.Option> options = const [],
  }) async {
    directories.add(req.directoryId);

    if (req.directoryId == _photos) {
      return filesystem.FilesystemSearchResponse(
        next: req,
        items: [
          _directory(_nested, "2026", parent: _photos),
          _file(uuidx.withSuffix(3), "take.five.mp3", parent: _photos),
        ],
        breadcrumb: [_directory(_photos, "photos")],
      );
    }

    return filesystem.FilesystemSearchResponse(
      next: req,
      items: [_directory(_photos, "photos"), _file(uuidx.withSuffix(4), "loose.mp3")],
      breadcrumb: [],
    );
  }
}

Widget _harness(_Searches searches) => filesystem.FilesystemBrowser(
  search: searches.call,
  mode: ValueNotifier(media.SearchMode.filesystem),
  onModeChanged: (_) {},
);

void main() {
  testWidgets('lists the root and shows no parent entry there', (WidgetTester tester) async {
    final searches = _Searches();
    await tester.pumpApp(_harness(searches));
    await tester.pumpAndSettle();

    expect(searches.directories, [uuidx.min()]);
    expect(find.text('photos'), findsOneWidget);
    expect(find.text('loose.mp3'), findsOneWidget);

    // the root has nowhere to go up to.
    expect(find.text('..'), findsNothing);
  });

  testWidgets('opening a directory lists it and offers a parent entry', (WidgetTester tester) async {
    final searches = _Searches();
    await tester.pumpApp(_harness(searches));
    await tester.pumpAndSettle();

    await tester.tap(find.text('photos'));
    await tester.pumpAndSettle();

    expect(searches.directories, [uuidx.min(), _photos]);
    expect(find.text('2026'), findsOneWidget);
    expect(find.text('take.five.mp3'), findsOneWidget);
    expect(find.text('..'), findsOneWidget);
    expect(find.text('loose.mp3'), findsNothing);
  });

  testWidgets('the parent entry returns to the containing directory', (WidgetTester tester) async {
    final searches = _Searches();
    await tester.pumpApp(_harness(searches));
    await tester.pumpAndSettle();

    await tester.tap(find.text('photos'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('..'));
    await tester.pumpAndSettle();

    // a single ancestor means the parent is the root.
    expect(searches.directories, [uuidx.min(), _photos, uuidx.min()]);
    expect(find.text('loose.mp3'), findsOneWidget);
  });

  testWidgets('the path appears in the search hint', (WidgetTester tester) async {
    final searches = _Searches();
    await tester.pumpApp(_harness(searches));
    await tester.pumpAndSettle();

    expect(find.text('search the library'), findsOneWidget);

    await tester.tap(find.text('photos'));
    await tester.pumpAndSettle();

    expect(find.text('search in photos'), findsOneWidget);
  });
}
