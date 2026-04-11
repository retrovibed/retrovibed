import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/community/publish.container.dart';
import 'package:retrovibed/media/media.pb.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

PublishContainer _build({
  VoidCallback? onPublished,
  VoidCallback? onCancel,
}) => PublishContainer(
  onPublished: onPublished ?? () {},
  onCancel: onCancel ?? () {},
);

void main() {
  group('PublishContainer', () {
    testWidgets('renders step 0 with Library indicator', (tester) async {
      await tester.pumpApp(_build(
        onCancel: () {},
        onPublished: () {},
      ));
      await tester.pumpAndSettle();

      expect(find.text('Library'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('back button absent on step 0', (tester) async {
      await tester.pumpApp(_build());
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.arrow_back), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('close button calls onCancel', (tester) async {
      bool cancelled = false;

      await tester.pumpApp(_build(onCancel: () => cancelled = true));
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.close));
      await tester.pump();

      expect(cancelled, isTrue);
      expect(tester.takeException(), isNull);
    });

    testWidgets('selecting content advances to step 1', (tester) async {
      await tester.pumpApp(
        PublishContainer(
          onPublished: () {},
          onCancel: () {},
          search: (req, {options = const []}) => Future.value(
            MediaSearchResponse(
              next: MediaSearchRequest(),
              items: [Media(description: 'Test Video', mimetype: 'video/mp4')],
            ),
          ),
          upload: (_) => Future.value(MediaUploadResponse()),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('Test Video'));
      await tester.pumpAndSettle();

      expect(find.text('Media'), findsOneWidget);
      expect(find.byIcon(Icons.arrow_back), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without overflow', (tester) async {
      final entry = _resolutions.currentValue!;

      await tester.pumpApp(
        physicalSize: entry.value,
        _build(),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    }, variant: _resolutions);
  });
}
