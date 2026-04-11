import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:media_kit/media_kit.dart';
import 'package:retrovibed/media/player.settings.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

class _FakePlatformPlayer extends PlatformPlayer {
  _FakePlatformPlayer()
      : super(configuration: const PlayerConfiguration());
}

Player _fakePlayer() => Player(platformPlayer: _FakePlatformPlayer());

void main() {
  group('PlayerSettings', () {
    final resolutions = Resolutions.variant();

    testWidgets(
      'renders without overflow',
      (tester) async {
        final player = _fakePlayer();
        addTearDown(player.dispose);

        await tester.pumpApp(
          PlayerSettings(current: player),
          physicalSize: resolutions.currentValue!.value,
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      },
      variant: resolutions,
    );

    testWidgets(
      'renders with padding without overflow',
      (tester) async {
        final player = _fakePlayer();
        addTearDown(player.dispose);

        await tester.pumpApp(
          PlayerSettings(
            current: player,
            padding: const EdgeInsets.all(16),
          ),
          physicalSize: resolutions.currentValue!.value,
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      },
      variant: resolutions,
    );

    testWidgets(
      'renders with margin without overflow',
      (tester) async {
        final player = _fakePlayer();
        addTearDown(player.dispose);

        await tester.pumpApp(
          PlayerSettings(
            current: player,
            margin: const EdgeInsets.all(16),
          ),
          physicalSize: resolutions.currentValue!.value,
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      },
      variant: resolutions,
    );
  });
}
