import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:media_kit/media_kit.dart';
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/routes.dart' as routes;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

Widget _harness() => media.Playlist(const routes.Routes());

const _movieIndex = 0;
const _downloadIndex = 1;
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
// download page) overrides the tap-driven animateTo. Pumping in small
// increments lets that frame actually land.
Future<void> _pumpTabSwitch(WidgetTester tester) async {
  for (var i = 0; i < 10; i++) {
    await tester.pump(const Duration(milliseconds: 50));
  }
}

// Below Defaults._defaultCompact (400.0 logical pixels), forcing
// ds.Defaults.isCompact true, which switches routes.dart to
// bottomNavigationBar and reverses ds.Table's leading/trailing order.
const _compactPhysicalSize = Size(360, 640);

void main() {
  setUpAll(() {
    MediaKit.ensureInitialized();
  });

  testWidgets('movie -> download selects download, not community', (
    WidgetTester tester,
  ) async {
    await tester.pumpApp(_harness(), fit: FlexFit.tight);
    await _pumpUntilTabBarVisible(tester);
    expect(_selectedTabIndex(tester), _movieIndex);

    await tester.tap(find.byIcon(Icons.download));
    await _pumpTabSwitch(tester);

    expect(_selectedTabIndex(tester), _downloadIndex);
  });

  testWidgets('movie -> settings -> download selects download, not community', (
    WidgetTester tester,
  ) async {
    await tester.pumpApp(_harness(), fit: FlexFit.tight);
    await _pumpUntilTabBarVisible(tester);

    await tester.tap(find.byIcon(Icons.settings));
    await _pumpTabSwitch(tester);
    expect(_selectedTabIndex(tester), _settingsIndex);

    await tester.tap(find.byIcon(Icons.download));
    await _pumpTabSwitch(tester);

    expect(_selectedTabIndex(tester), _downloadIndex);
  });

  testWidgets('compact: movie -> download selects download, not community', (
    WidgetTester tester,
  ) async {
    await tester.pumpApp(_harness(), fit: FlexFit.tight, physicalSize: _compactPhysicalSize);
    await _pumpUntilTabBarVisible(tester);
    expect(_selectedTabIndex(tester), _movieIndex);

    await tester.tap(find.byIcon(Icons.download));
    await _pumpTabSwitch(tester);

    expect(_selectedTabIndex(tester), _downloadIndex);
  });

  testWidgets('compact: movie -> settings -> download selects download, not community', (
    WidgetTester tester,
  ) async {
    await tester.pumpApp(_harness(), fit: FlexFit.tight, physicalSize: _compactPhysicalSize);
    await _pumpUntilTabBarVisible(tester);

    await tester.tap(find.byIcon(Icons.settings));
    await _pumpTabSwitch(tester);
    expect(_selectedTabIndex(tester), _settingsIndex);

    await tester.tap(find.byIcon(Icons.download));
    await _pumpTabSwitch(tester);

    expect(_selectedTabIndex(tester), _downloadIndex);
  });
}
