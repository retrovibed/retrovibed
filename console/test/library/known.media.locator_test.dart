import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/library/known.media.locator.dart';
import 'package:retrovibed/library/known.media.card.dart';
import 'package:retrovibed/library/api.dart' as api;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('KnownMediaLocator', () {
    testWidgets('renders without overflow', (tester) async {
      await tester.pumpApp(
        KnownMediaLocator(
          api.Known(description: 'Test', summary: 'summary'),
          locate: (req, {options = const []}) async => api.LocateCreateResponse(locate: req),
        ),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    });

    testWidgets('tapping calls locate with the knownMediaId', (tester) async {
      api.Locate? requested;
      final item = api.Known(id: 'known-1', description: 'Test', summary: 'summary');
      await tester.pumpApp(
        KnownMediaLocator(
          item,
          locate: (req, {options = const []}) async {
            requested = req;
            return api.LocateCreateResponse(locate: req);
          },
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(KnownMediaCard));
      await tester.pumpAndSettle();

      expect(requested?.knownMediaId, equals(item.id));
    });

    testWidgets('a failed locate is handled internally', (tester) async {
      final item = api.Known(id: 'known-1', description: 'Test', summary: 'summary');
      await tester.pumpApp(
        KnownMediaLocator(
          item,
          locate: (req, {options = const []}) => Future.error('boom'),
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(KnownMediaCard));
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });
  });
}
