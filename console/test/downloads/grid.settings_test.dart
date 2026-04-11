import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/downloads/grid.settings.dart';
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/torrents/api.dart' as torrents_api;
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/wireguard/api.dart' as wireguard_api;
import 'package:nock/nock.dart';

final _resolutions = Resolutions.variant();

void main() {
  group('GridSettings', () {
    setUp(() {
      nock('https://api.retrovibe.space')
          .get(RegExp(r'/s/torrents/'))
          .reply(
            200,
            torrents_api.TorrentSettings.create().toProto3Json(),
            headers: {'Content-Type': 'application/json'},
          );
    });

    testWidgets('renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        GridSettings(
          media.discoveredsearch.request(limit: 32),
          onChange: (_) {},
          wgcurrent: () async {
            return wireguard_api.WireguardCurrentResponse.create()
              ..wireguard = wireguard_api.Wireguard();
          },
          wgupdate: (_, {options = const []}) async {
            return wireguard_api.WireguardUpdateResponse.create();
          },
        ),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);
  });
}
