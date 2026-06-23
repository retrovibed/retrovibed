import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/modals.dart' as modals;
import 'package:retrovibed/discovery/api.dart' as api;
import 'package:retrovibed/discovery/settings.locate.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('LocateSettings', () {
    testWidgets('renders the p2p locate checkbox', (tester) async {
      await tester.pumpApp(
        modals.Node(LocateSettings(api.DiscoverySettings(locateP2p: false))),
      );
      await tester.pumpAndSettle();

      expect(find.text('p2p locate'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('tapping while not yet acknowledged shows the confirmation instead of toggling', (tester) async {
      bool changed = false;
      await tester.pumpApp(
        modals.Node(
          LocateSettings(
            api.DiscoverySettings(locateP2p: false),
            onChange: (v) async {
              changed = true;
              return v;
            },
            disclaimer: (_) => false,
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('p2p locate'));
      await tester.pumpAndSettle();

      expect(find.text('Yes'), findsOneWidget);
      expect(find.text('No'), findsOneWidget);
      expect(changed, isFalse);
    });

    testWidgets('declining the confirmation does not acknowledge or toggle', (tester) async {
      bool changed = false;
      bool acknowledged = false;
      await tester.pumpApp(
        modals.Node(
          LocateSettings(
            api.DiscoverySettings(locateP2p: false),
            onChange: (v) async {
              changed = true;
              return v;
            },
            disclaimer: (_) => false,
            acknowledge: (_) => acknowledged = true,
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('p2p locate'));
      await tester.pumpAndSettle();
      await tester.tap(find.text('No'));
      await tester.pumpAndSettle();

      expect(changed, isFalse);
      expect(acknowledged, isFalse);
      expect(tester.takeException(), isNull);
    });

    testWidgets('accepting the confirmation acknowledges and toggles locateP2p', (tester) async {
      api.DiscoverySettings? changedTo;
      String? acknowledgedId;
      await tester.pumpApp(
        modals.Node(
          LocateSettings(
            api.DiscoverySettings(locateP2p: false),
            onChange: (v) async {
              changedTo = v;
              return v;
            },
            disclaimer: (_) => false,
            acknowledge: (id) => acknowledgedId = id,
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('p2p locate'));
      await tester.pumpAndSettle();
      await tester.tap(find.text('Yes'));
      await tester.pumpAndSettle();

      expect(acknowledgedId, equals('discovery.p2p'));
      expect(changedTo?.locateP2p, isTrue);
      expect(tester.takeException(), isNull);
    });

    testWidgets('already acknowledged toggles immediately without showing the confirmation', (tester) async {
      api.DiscoverySettings? changedTo;
      await tester.pumpApp(
        modals.Node(
          LocateSettings(
            api.DiscoverySettings(locateP2p: false),
            onChange: (v) async {
              changedTo = v;
              return v;
            },
            disclaimer: (_) => true,
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('p2p locate'));
      await tester.pumpAndSettle();

      expect(find.text('Yes'), findsNothing);
      expect(changedTo?.locateP2p, isTrue);
      expect(tester.takeException(), isNull);
    });
  });
}
