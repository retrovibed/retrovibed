import 'package:fixnum/fixnum.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/media/download.display.dart';
import 'package:retrovibed/media/api.dart' as api;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

api.Download _download() => api.Download(
  media: api.Media(
    id: 'test-id-1',
    description: 'Test Media Description',
    mimetype: 'video/mp4',
    createdAt: '2025-01-01T00:00:00Z',
    archiveId: uuidx.min(),
    torrentId: uuidx.min(),
    knownMediaId: uuidx.min(),
  ),
  path: '/media/test-id-1.mp4',
  bytes: Int64(1024 * 1024 * 500),
);

api.Download _downloadLongFields() => api.Download(
  media: api.Media(
    id: 'very-long-identifier-that-could-potentially-cause-overflow-in-the-display-widget',
    description:
        'An Extremely Long Title That Could Potentially Overflow The Download Display Widget And Cause Layout Issues In The UI',
    mimetype: 'video/mp4',
    createdAt: '2025-01-01T00:00:00Z',
    archiveId: uuidx.min(),
    torrentId: uuidx.min(),
    knownMediaId: uuidx.min(),
  ),
  path: '/media/very/deeply/nested/path/that/could/overflow/the/display/widget/test-id-1.mp4',
  bytes: Int64(1024 * 1024 * 500),
);

final _resolutions = Resolutions.variant();

void main() {
  group('DownloadDisplay', () {
    testWidgets('renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        DownloadDisplay(_download()),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders with long fields without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        DownloadDisplay(_downloadLongFields()),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders with onTap and onVerify without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        DownloadDisplay(
          _download(),
          onTap: () async {},
          onVerify: (_) async {},
        ),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);
  });
}
