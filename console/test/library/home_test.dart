import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/library/dropdown.upload.dart';
import 'package:retrovibed/library/home.dart';
import 'package:retrovibed/library/known.media.display.dart';
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

void main() {
  group('Home', () {
    testWidgets('hides library items once discover mode is activated', (
      WidgetTester tester,
    ) async {
      final entry = _resolutions.currentValue!;
      final search = ValueNotifier<media.MediaSearchState>(
        media.MediaSearchState(next: media.media.request(limit: 32)),
      );
      addTearDown(search.dispose);
      await tester.pumpApp(
        Home(
          apisearch: (req, {options = const []}) async {
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
              ],
              next: req,
            );
          },
          search: search,
          highlighted: '',
        ),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(find.byType(KnownMediaDisplay), findsOneWidget);

      await tester.tap(find.byType(DropdownUpload));
      await tester.pumpAndSettle();
      await tester.ensureVisible(find.text('Search'));
      await tester.tap(find.text('Search'));
      await tester.pumpAndSettle();

      expect(find.byType(KnownMediaDisplay), findsNothing);
    }, variant: _resolutions);
  });
}
