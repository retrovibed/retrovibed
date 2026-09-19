import 'package:fixnum/fixnum.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/torrentx/display.dart';
import 'package:retrovibed/torrentx/api.dart' as api;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

api.TorrentInfoResponse _torrent() => api.TorrentInfoResponse(
  meta: api.TorrentMeta(
    comment: 'a test torrent',
    encoding: 'UTF-8',
    createdBy: 'mktorrent 1.1',
    creationDate: Int64(1700000000),
    announceList: ['https://example.com/announce'],
    urlList: ['https://example.com/webseed'],
  ),
  details: api.TorrentDetails(
    name: 'example.1',
    length: Int64(1024 * 1024 * 500),
    source: 'example',
    private: true,
  ),
  files: [
    api.TorrentFile(name: 'example.1', length: Int64(1024 * 1024 * 500), path: 'example.1'),
  ],
);

api.TorrentInfoResponse _torrentLongFields() => api.TorrentInfoResponse(
  meta: api.TorrentMeta(
    comment: 'An Extremely Long Comment That Could Potentially Overflow The Torrent Display Widget And Cause Layout Issues',
    encoding: 'UTF-8',
    createdBy: 'mktorrent 1.1 with a surprisingly long created-by string for testing overflow behavior',
    creationDate: Int64(1700000000),
    announceList: List.generate(10, (i) => 'https://tracker-$i.example.com/announce/with/a/long/path'),
    urlList: List.generate(5, (i) => 'https://webseed-$i.example.com/with/a/long/path'),
  ),
  details: api.TorrentDetails(
    name: 'a-very-long-torrent-name-that-could-potentially-overflow-the-display-widget-in-narrow-layouts',
    length: Int64(1024 * 1024 * 1024 * 50),
    source: 'a-very-long-source-string-for-testing-overflow-in-narrow-layouts',
    private: false,
  ),
  files: List.generate(
    20,
    (i) => api.TorrentFile(
      name: 'very/deeply/nested/path/that/could/overflow/the/display/widget/file-$i.bin',
      length: Int64(1024 * 1024 * (i + 1)),
      path: 'very/deeply/nested/path/that/could/overflow/the/display/widget/file-$i.bin',
    ),
  ),
);

final _resolutions = Resolutions.variant();

void main() {
  group('TorrentDisplay', () {
    testWidgets('renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(SingleChildScrollView(child: TorrentDisplay(_torrent())), physicalSize: entry.value);
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders with long fields and many files without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        SingleChildScrollView(child: TorrentDisplay(_torrentLongFields())),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders transfer ratio', (WidgetTester tester) async {
      await tester.pumpApp(
        SingleChildScrollView(
          child: TorrentDisplay(
            api.TorrentInfoResponse(
              details: api.TorrentDetails(downloaded: Int64(1024 * 1024), uploaded: Int64(1024 * 1024 * 3)),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();
      expect(find.text('ratio'), findsOneWidget);
      expect(find.text('3.00'), findsOneWidget);
    });

    testWidgets('renders placeholder ratio when nothing downloaded', (WidgetTester tester) async {
      await tester.pumpApp(
        SingleChildScrollView(
          child: TorrentDisplay(
            api.TorrentInfoResponse(details: api.TorrentDetails(uploaded: Int64(1024 * 1024))),
          ),
        ),
      );
      await tester.pumpAndSettle();
      expect(find.text('-'), findsOneWidget);
    });

    testWidgets('renders empty response without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        SingleChildScrollView(child: TorrentDisplay(api.TorrentInfoResponse())),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);
  });
}
