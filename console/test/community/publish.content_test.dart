import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/community/publish.content.dart';
import 'package:retrovibed/media/media.pb.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

void main() {
  group('PublishContent', () {
    testWidgets('renders without overflow', (tester) async {
      final entry = _resolutions.currentValue!;

      await tester.pumpApp(
        physicalSize: entry.value,
        SingleChildScrollView(
          child: PublishContent(
            onSelect: (_) {},
            search: (req, {host, options = const []}) => Future.value(
              MediaSearchResponse(
                next: MediaSearchRequest(),
                items: [],
              ),
            ),
            upload: (_) => Future.value(MediaUploadResponse()),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('selecting an item calls onSelect with the download',
        (tester) async {
      Download? received;
      final item = Media(description: 'My Video', mimetype: 'video/mp4');

      await tester.pumpApp(
        PublishContent(
          onSelect: (d) => received = d,
          search: (req, {host, options = const []}) => Future.value(
            MediaSearchResponse(
              next: MediaSearchRequest(),
              items: [item],
            ),
          ),
          upload: (_) => Future.value(MediaUploadResponse()),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('My Video'));
      await tester.pump();

      expect(received, isNotNull);
      expect(received!.media.description, equals('My Video'));
    });
  });
}
