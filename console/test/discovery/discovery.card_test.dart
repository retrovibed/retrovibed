import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/ddisc.dart' as ddisc;
import 'package:retrovibed/discovery/discovery.card.dart';
import 'package:retrovibed/discovery/locate.p2p.prompt.dart' as p2p;
import 'package:retrovibed/library/known.media.card.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('DiscoveredCard', () {
    testWidgets('renders without overflow', (tester) async {
      await tester.pumpApp(
        DiscoveredCard(ddisc.Discovery(id: 'd-1', title: 'Ubuntu', description: 'iso')),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    });

    testWidgets('tapping downloads once p2p consent resolves', (tester) async {
      String? downloadedId;
      final current = ddisc.Discovery(id: 'd-1', title: 'Ubuntu', description: 'iso');
      await tester.pumpApp(
        DiscoveredCard(
          current,
          ensureP2P: (context, {options = const []}) async {},
          download: (id, {discovery, autodownload = false, options = const []}) async {
            downloadedId = id;
            return ddisc.DiscoveryDownloadResponse.create();
          },
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(KnownMediaCard));
      await tester.pumpAndSettle();

      expect(downloadedId, equals(current.id));
      expect(tester.takeException(), isNull);
    });

    testWidgets('a declined p2p consent blocks download and shows no error', (tester) async {
      bool downloadCalled = false;
      final current = ddisc.Discovery(id: 'd-1', title: 'Ubuntu', description: 'iso');
      await tester.pumpApp(
        DiscoveredCard(
          current,
          ensureP2P: (context, {options = const []}) => Future.error(const p2p.P2PConsentDeclined()),
          download: (id, {discovery, autodownload = false, options = const []}) async {
            downloadCalled = true;
            return ddisc.DiscoveryDownloadResponse.create();
          },
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(KnownMediaCard));
      await tester.pumpAndSettle();

      expect(downloadCalled, isFalse);
      expect(tester.takeException(), isNull);
    });
  });
}
