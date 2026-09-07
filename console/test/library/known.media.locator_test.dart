import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:retrovibed/ddisc.dart' as ddisc;
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

    testWidgets('future factory renders without overflow', (tester) async {
      await tester.pumpApp(
        KnownMediaLocator.future(
          Future.value(api.Known(description: 'Test', summary: 'summary')),
        ),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    });

    testWidgets('tapping calls locate with the knownMediaId', (tester) async {
      api.Locate? requested;
      final item = api.Known(id: 'row-1', uid: 'known-1', description: 'Test', summary: 'summary');
      await tester.pumpApp(
        KnownMediaLocator(
          item,
          ensureP2P: (context, {options = const []}) async => true,
          locate: (req, {options = const []}) async {
            requested = req;
            return api.LocateCreateResponse(locate: req);
          },
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(KnownMediaCard));
      await tester.pumpAndSettle();

      expect(requested?.knownMediaId, equals(item.uid));
    });

    testWidgets('tapping forwards the known item adult flag', (tester) async {
      api.Locate? requested;
      final item = api.Known(id: 'row-1', uid: 'known-1', description: 'Test', summary: 'summary', adult: true);
      await tester.pumpApp(
        KnownMediaLocator(
          item,
          ensureP2P: (context, {options = const []}) async => true,
          locate: (req, {options = const []}) async {
            requested = req;
            return api.LocateCreateResponse(locate: req);
          },
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(KnownMediaCard));
      await tester.pumpAndSettle();

      expect(requested?.adult, isTrue);
    });

    testWidgets('a failed locate is handled internally', (tester) async {
      final item = api.Known(id: 'known-1', description: 'Test', summary: 'summary');
      await tester.pumpApp(
        KnownMediaLocator(
          item,
          ensureP2P: (context, {options = const []}) async => true,
          locate: (req, {options = const []}) => Future.error('boom'),
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(KnownMediaCard));
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('a 404 deleting the recommendation after a successful locate is treated as success', (tester) async {
      bool onChangeCalled = false;
      final item = api.Known(id: 'known-1', description: 'Test', summary: 'summary');
      await tester.pumpApp(
        KnownMediaLocator(
          item,
          onChange: (v) => onChangeCalled = true,
          ensureP2P: (context, {options = const []}) async => true,
          locate: (req, {options = const []}) async => api.LocateCreateResponse(locate: req),
          delete: (id, {options = const []}) => Future.error(http.Response('', 404)),
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(KnownMediaCard));
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
      expect(onChangeCalled, isTrue);
    });

    testWidgets('tapping a discovered known item calls download with the uid and skips locate', (tester) async {
      String? downloadedId;
      bool locateCalled = false;
      final item =
          api.Known(id: 'known-1', uid: 'known-1', description: 'Test', summary: 'summary', source: ddisc.sources.discovered);
      await tester.pumpApp(
        KnownMediaLocator(
          item,
          ensureP2P: (context, {options = const []}) async => true,
          download: (id, {options = const []}) async {
            downloadedId = id;
            return ddisc.DiscoveryDownloadResponse.create();
          },
          locate: (req, {options = const []}) async {
            locateCalled = true;
            return api.LocateCreateResponse(locate: req);
          },
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(KnownMediaCard));
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
      expect(downloadedId, equals(item.uid));
      expect(locateCalled, isFalse);
    });

    testWidgets('tapping a searchplugin known item calls download with the uid and skips locate', (tester) async {
      String? downloadedId;
      bool locateCalled = false;
      final item = api.Known(
        id: 'known-1',
        uid: 'known-1',
        description: 'Test',
        summary: 'summary',
        source: ddisc.sources.searchplugin,
      );
      await tester.pumpApp(
        KnownMediaLocator(
          item,
          ensureP2P: (context, {options = const []}) async => true,
          download: (id, {options = const []}) async {
            downloadedId = id;
            return ddisc.DiscoveryDownloadResponse.create();
          },
          locate: (req, {options = const []}) async {
            locateCalled = true;
            return api.LocateCreateResponse(locate: req);
          },
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(KnownMediaCard));
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
      expect(downloadedId, equals(item.uid));
      expect(locateCalled, isFalse);
    });

    testWidgets('a 404 deleting the recommendation after a successful download is treated as success', (tester) async {
      bool onChangeCalled = false;
      final item = api.Known(id: 'known-1', description: 'Test', summary: 'summary', source: ddisc.sources.discovered);
      await tester.pumpApp(
        KnownMediaLocator(
          item,
          onChange: (v) => onChangeCalled = true,
          ensureP2P: (context, {options = const []}) async => true,
          download: (id, {options = const []}) async => ddisc.DiscoveryDownloadResponse.create(),
          delete: (id, {options = const []}) => Future.error(http.Response('', 404)),
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(KnownMediaCard));
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
      expect(onChangeCalled, isTrue);
    });

    testWidgets('a 404 deleting the recommendation on long press is treated as success', (tester) async {
      bool onChangeCalled = false;
      final item = api.Known(id: 'known-1', description: 'Test', summary: 'summary');
      await tester.pumpApp(
        KnownMediaLocator(
          item,
          onChange: (v) => onChangeCalled = true,
          delete: (id, {options = const []}) => Future.error(http.Response('', 404)),
        ),
      );
      await tester.pumpAndSettle();
      await tester.longPress(find.byType(KnownMediaCard));
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
      expect(onChangeCalled, isTrue);
    });
  });
}
