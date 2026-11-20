import 'package:material_ui/material_ui.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/library/network.settings.dart';
import 'package:retrovibed/torrents.dart' as torrents;
import 'package:retrovibed/dht.dart' as dht;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

void main() {
  group('NetworkSettings', () {
    testWidgets('renders torrent and dht sections without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        physicalSize: entry.value,
        SingleChildScrollView(
          child: NetworkSettings(
            torrent: ({options = const []}) => Future.value(torrents.TorrentSettings()),
            dhts: ({options = const []}) => Future.value(dht.DHTSettings()),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('sharing'), findsOneWidget);
      expect(find.text('download rate'), findsOneWidget);
      expect(find.text('firewall'), findsOneWidget);
      expect(find.text('dht'), findsOneWidget);
      expect(find.text('passive'), findsOneWidget);
      expect(find.text('bootstrap'), findsOneWidget);
      expect(find.text('info'), findsNWidgets(2));
      expect(find.text('debug'), findsNWidgets(2));
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);
  });
}
