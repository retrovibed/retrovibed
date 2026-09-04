import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/library/search.dart';
import 'package:retrovibed/library/dropdown.upload.dart';
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

Future<media.MediaSearchResponse> _mockSearchEmpty(
  media.MediaSearchRequest req, {
  String? host,
  List<httpx.Option> options = const [],
}) async {
  return media.MediaSearchResponse(items: [], next: req);
}

Future<media.MediaSearchResponse> _mockSearchWithItems(
  media.MediaSearchRequest req, {
  String? host,
  List<httpx.Option> options = const [],
}) async {
  return media.MediaSearchResponse(
    items: [
      media.Media(
        id: uuidx.withSuffix(1),
        description: 'Test Media One',
        mimetype: 'video/mp4',
        createdAt: '2025-01-01T00:00:00Z',
        archiveId: uuidx.min(),
        torrentId: uuidx.min(),
        knownMediaId: uuidx.min(),
      ),
      media.Media(
        id: uuidx.withSuffix(2),
        description: 'Test Media Two',
        mimetype: 'audio/mp3',
        createdAt: '2025-01-02T00:00:00Z',
        archiveId: uuidx.min(),
        torrentId: uuidx.min(),
        knownMediaId: uuidx.min(),
      ),
    ],
    next: req,
  );
}

Future<media.MediaSearchResponse> _mockSearchWithLongNames(
  media.MediaSearchRequest req, {
  String? host,
  List<httpx.Option> options = const [],
}) async {
  return media.MediaSearchResponse(
    items: [
      media.Media(
        id: uuidx.withSuffix(1),
        description:
            'An Extremely Long Title That Could Potentially Overflow The Display Widget And Cause Layout Issues',
        mimetype: 'video/mp4',
        createdAt: '2025-01-01T00:00:00Z',
        archiveId: uuidx.min(),
        torrentId: uuidx.min(),
        knownMediaId: uuidx.min(),
      ),
    ],
    next: req,
  );
}

final _resolutions = Resolutions.variant();

ValueNotifier<media.MediaSearchState> _search({String query = ''}) =>
    ValueNotifier(media.MediaSearchState(next: media.media.request(limit: 32, query: query)));

Widget _buildSearch({
  required media.FnMediaSearch apisearch,
  required ValueNotifier<media.MediaSearchState> search,
  ValueNotifier<media.SearchMode>? mode,
  void Function(media.SearchMode)? onModeChanged,
  String highlighted = '',
}) {
  return Search(
    apisearch: apisearch,
    highlighted: highlighted,
    search: search,
    mode: mode ?? ValueNotifier(media.SearchMode.library),
    onModeChanged: onModeChanged ?? (_) {},
    downloading: const SizedBox.shrink(),
    onDownloadingChanged: (_) {},
  );
}

void main() {
  group('Search', () {
    testWidgets('renders empty without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      final search = _search();
      addTearDown(search.dispose);
      await tester.pumpApp(
        _buildSearch(apisearch: _mockSearchEmpty, search: search),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders with items without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      final search = _search();
      addTearDown(search.dispose);
      await tester.pumpApp(
        _buildSearch(apisearch: _mockSearchWithItems, search: search),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders with long names without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      final search = _search();
      addTearDown(search.dispose);
      await tester.pumpApp(
        _buildSearch(apisearch: _mockSearchWithLongNames, search: search),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders empty without overflow at narrow 300x600', (WidgetTester tester) async {
      final search = _search();
      addTearDown(search.dispose);
      await tester.pumpApp(
        _buildSearch(apisearch: _mockSearchEmpty, search: search),
        physicalSize: const Size(300, 600),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders with items without overflow at narrow 300x600', (WidgetTester tester) async {
      final search = _search();
      addTearDown(search.dispose);
      await tester.pumpApp(
        _buildSearch(apisearch: _mockSearchWithItems, search: search),
        physicalSize: const Size(300, 600),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    });

    testWidgets('submitting the search tray updates the query and resets the offset', (
      WidgetTester tester,
    ) async {
      final search = _search();
      addTearDown(search.dispose);
      search.value.next.offset = ds.Int64(5);
      await tester.pumpApp(
        _buildSearch(apisearch: _mockSearchEmpty, search: search),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), 'something');
      await tester.testTextInput.receiveAction(TextInputAction.done);
      await tester.pumpAndSettle();

      expect(search.value.next.query, equals('something'));
      expect(search.value.next.offset, equals(ds.Int64.ZERO));
      expect(tester.takeException(), isNull);
    });

    testWidgets('switching to discovery mode invokes onModeChanged', (WidgetTester tester) async {
      final search = _search();
      addTearDown(search.dispose);
      media.SearchMode? changedTo;
      await tester.pumpApp(
        _buildSearch(
          apisearch: _mockSearchEmpty,
          search: search,
          onModeChanged: (m) => changedTo = m,
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byType(DropdownUpload));
      await tester.pumpAndSettle();
      await tester.ensureVisible(find.text('Discover'));
      await tester.tap(find.text('Discover'));
      await tester.pumpAndSettle();

      expect(changedTo, equals(media.SearchMode.discovery));
      expect(tester.takeException(), isNull);
    });
  });
}
