import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:media_kit/media_kit.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/routes.dart' as routes;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

// DaemonDropdown mounts on both the remote and settings tabs and scans for
// peers via api.daemons.discover/search on initState. Those hit a real
// socket/HTTP client, which flutter_test's HttpOverrides blocks with its own
// non-String error object - so every tab-switch test needs a no-op fake here
// rather than the live network defaults.
Future<meta.DaemonSearchResponse> _fakeDaemonSearch(meta.DaemonSearchRequest req) {
  return Future.value(meta.DaemonSearchResponse());
}

Future<Stream<meta.Daemon>> _noopDaemonDiscover({List<httpx.Option> options = const []}) async {
  return const Stream<meta.Daemon>.empty();
}

Widget _harness() => media.Playlist(
  routes.Routes(daemonSearch: _fakeDaemonSearch, daemonDiscover: _noopDaemonDiscover),
);

const _movieIndex = 0;
const _remoteIndex = 1;
const _settingsIndex = 3;

Future<void> _pumpUntilTabBarVisible(WidgetTester tester) async {
  for (var i = 0; i < 30 && find.byType(TabBar).evaluate().isEmpty; i++) {
    await tester.pump(const Duration(milliseconds: 100));
  }
}

int _selectedTabIndex(WidgetTester tester) {
  final context = tester.element(find.byType(TabBar));
  return DefaultTabController.of(context).index;
}

// A single large pump(duration) advances the clock in one jump and can skip
// over the frame where the community page's Scrollable.ensureVisible call
// (fired from its search tray mounting as a cache-adjacent neighbor of the
// remote page) overrides the tap-driven animateTo. Pumping in small
// increments lets that frame actually land.
Future<void> _pumpTabSwitch(WidgetTester tester) async {
  for (var i = 0; i < 10; i++) {
    await tester.pump(const Duration(milliseconds: 50));
  }
}

void main() {
  setUpAll(() {
    MediaKit.ensureInitialized();
  });

  // Runs against Resolutions.all (mobile through desktop) so both the
  // compact (bottomNavigationBar, reversed ds.Table leading/trailing) and
  // non-compact (appBar) layouts in routes.dart are exercised.
  final resolutions = Resolutions.variant();

  testWidgets(
    'movie -> remote selects remote, not community',
    (WidgetTester tester) async {
      await tester.pumpApp(_harness(), fit: FlexFit.tight, physicalSize: resolutions.currentValue!.value);
      await _pumpUntilTabBarVisible(tester);
      expect(_selectedTabIndex(tester), _movieIndex);

      await tester.tap(find.byIcon(Icons.settings_remote));
      await _pumpTabSwitch(tester);

      expect(_selectedTabIndex(tester), _remoteIndex);
    },
    variant: resolutions,
  );

  testWidgets(
    'movie -> settings -> remote selects remote, not community',
    (WidgetTester tester) async {
      await tester.pumpApp(_harness(), fit: FlexFit.tight, physicalSize: resolutions.currentValue!.value);
      await _pumpUntilTabBarVisible(tester);

      await tester.tap(find.byIcon(Icons.settings));
      await _pumpTabSwitch(tester);
      expect(_selectedTabIndex(tester), _settingsIndex);

      await tester.tap(find.byIcon(Icons.settings_remote));
      await _pumpTabSwitch(tester);

      expect(_selectedTabIndex(tester), _remoteIndex);
    },
    variant: resolutions,
  );
}
