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
          disclaimer: (_) => true,
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
          disclaimer: (_) => true,
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
          disclaimer: (_) => true,
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(KnownMediaCard));
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('shows the disclaimer prompt instead of the card when not yet acknowledged', (tester) async {
      final item = api.Known(id: 'known-1', description: 'Test', summary: 'summary');
      await tester.pumpApp(
        KnownMediaLocator(item, disclaimer: (_) => false),
      );
      await tester.pumpAndSettle();

      expect(find.byType(KnownMediaCard), findsNothing);
      expect(find.text('Nevermind'), findsOneWidget);
      expect(find.text('P2P'), findsOneWidget);
      expect(find.text('Listed Only'), findsOneWidget);
    });

    testWidgets('choosing Nevermind does not acknowledge the disclaimer', (tester) async {
      bool acknowledged = false;
      final item = api.Known(id: 'known-1', description: 'Test', summary: 'summary');
      await tester.pumpApp(
        KnownMediaLocator(
          item,
          disclaimer: (_) => false,
          acknowledge: (_) => acknowledged = true,
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.text('Nevermind'));
      await tester.pumpAndSettle();

      expect(acknowledged, isFalse);
      expect(find.byType(KnownMediaCard), findsNothing);
    });

    testWidgets('choosing P2P acknowledges the disclaimer', (tester) async {
      bool acknowledged = false;
      final item = api.Known(id: 'known-1', description: 'Test', summary: 'summary');
      await tester.pumpApp(
        KnownMediaLocator(
          item,
          disclaimer: (_) => false,
          acknowledge: (_) => acknowledged = true,
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.text('P2P'));
      await tester.pumpAndSettle();

      expect(acknowledged, isTrue);
    });

    testWidgets('choosing Listed Only acknowledges the disclaimer', (tester) async {
      bool acknowledged = false;
      final item = api.Known(id: 'known-1', description: 'Test', summary: 'summary');
      await tester.pumpApp(
        KnownMediaLocator(
          item,
          disclaimer: (_) => false,
          acknowledge: (_) => acknowledged = true,
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.text('Listed Only'));
      await tester.pumpAndSettle();

      expect(acknowledged, isTrue);
    });
  });
}
